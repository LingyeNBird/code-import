package migrate

import (
	"ccrctl/pkg/api/cnb"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/system"
	"ccrctl/pkg/vcs"
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

const (
	GitDirName     = "source_git_dir"
	CnbUserName    = "cnb"
	MaxConcurrency = 20
)

var (
	CnbURL                                                                  = config.Cfg.GetString("cnb.url")
	CnbApiURL                                                               = config.ConvertToApiURL(CnbURL)
	CnbToken                                                                = config.Cfg.GetString("cnb.token")
	MigrateType                                                             = config.Cfg.GetString("migrate.type")
	Concurrency                                                             = config.Cfg.GetInt("migrate.concurrency")
	totalRepoNumber, skipRepoNumber, successfulRepoNumber, failedRepoNumber int
	IsForcePush                                                             = config.Cfg.GetBool("migrate.force_push")
	IgnoreLFSNotFoundError                                                  = config.Cfg.GetBool("migrate.ignore_lfs_notfound_error")
	useLfsMigrate                                                           = config.Cfg.GetBool("migrate.use_lfs_migrate")
	organizationMappingLevel                                                = config.Cfg.GetInt("migrate.organization_mapping_level")
	allowIncompletePush                                                     = config.Cfg.GetBool("migrate.allow_incomplete_push")
	SourcePlatformName                                                      = config.Cfg.GetString("source.platform")
	SkipExistsRepo                                                          = config.Cfg.GetBool("migrate.skip_exists_repo")
)

func Run() {
	if SourcePlatformName == "github" {
		organizationMappingLevel = 2
	}
	startTime := time.Now()     // 记录迁移开始时间
	err := config.CheckConfig() // 检查配置文件
	if err != nil {
		panic(err)
	}
	err = system.SetFileDescriptorLimit(system.Limit) // 设置文件描述符限制
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(GitDirName, 0755) // 创建Git工作目录
	if err != nil {
		panic(err)
	}
	logger.Logger.Infof("创建仓库工作目录%s成功", GitDirName)
	pwdDir, err := os.Getwd() // 获取当前工作目录
	if err != nil {
		panic(err)
	}
	gitDirABSPath := filepath.Join(pwdDir, GitDirName) // 构建Git目录的绝对路径
	defer func(path string) {                          // 延迟删除Git目录
		err := os.RemoveAll(path)
		if err != nil {
			panic(err)
		}
	}(gitDirABSPath)
	system.HandleInterrupt(gitDirABSPath) // 处理中断信号
	err = os.Chdir(GitDirName)            // 切换到Git工作目录
	if err != nil {
		panic(err)
	}
	depotList, err := vcs.NewVcs(SourcePlatformName) // 获取仓库列表
	if err != nil {
		panic(err)
	}
	totalRepoNumber, failedRepoNumber = len(depotList), len(depotList) // 初始化统计变量
	successfulRepoNumber, skipRepoNumber = 0, 0
	//err = cnb.CreateRootOrganizationIfNotExists(CnbApiURL, CnbToken) // 创建根组织
	//if err != nil {
	//	panic(err)
	//}
	exist, err := cnb.RootOrganizationExists(CnbApiURL, CnbToken)
	if err != nil {
		panic(err)
	}
	if !exist {
		logger.Logger.Errorf("根组织%s不存在，请先创建根组织", config.Cfg.GetString("cnb.root_organization"))
		return
	}
	if organizationMappingLevel == 1 { // 如果组织映射级别为1，则创建子组织
		err = cnb.CreateSubOrganizationIfNotExists(CnbApiURL, CnbToken, depotList)
		if err != nil {
			panic(err)
		}
	}
	if Concurrency > MaxConcurrency { // 限制并发数不超过最大值
		Concurrency = MaxConcurrency
	}
	logger.Logger.Infof("开始迁移仓库，当前并发数:%d", Concurrency)
	sem := semaphore.NewWeighted(int64(Concurrency)) // 创建信号量控制并发
	var wg sync.WaitGroup                            // 创建WaitGroup等待所有goroutine完成

	for _, depot := range depotList { // 遍历仓库列表
		wg.Add(1)                // 增加WaitGroup计数
		go func(depot vcs.VCS) { // 启动goroutine进行迁移
			defer wg.Done() // 完成后减少WaitGroup计数

			if err := sem.Acquire(context.Background(), 1); err != nil { // 获取信号量
				panic(err)
				return
			}
			defer sem.Release(1)   // 释放信号量
			err = migrateDo(depot) // 执行迁移操作
			if err != nil {
				logger.Logger.Errorf("%s 仓库迁移失败: %s", depot.GetRepoPath(), err) // 记录迁移失败信息
			}
		}(depot)
	}
	wg.Wait()                                                                                                                                                  // 等待所有goroutine完成
	duration := time.Since(startTime)                                                                                                                          // 计算迁移耗时
	logger.Logger.Infof("代码仓库迁移完成，耗时%s。\n【仓库总数】%d【成功迁移】%d【忽略迁移】%d【迁移失败】%d", duration, totalRepoNumber, successfulRepoNumber, skipRepoNumber, failedRepoNumber) // 记录迁移完成信息
}

func migrateDo(depot vcs.VCS) (err error) {
	repoName, subGroupName, repoPath, repoPrivate := depot.GetRepoName(), depot.GetSubGroupName(), depot.GetRepoPath(), depot.GetRepoPrivate()
	err, migrated := isMigrated(repoPath, logger.SuccessfulLogFilePath)
	if err != nil {
		return fmt.Errorf("%s 判断是否迁移失败%s", repoPath, err)
	}
	if migrated {
		skipRepoNumber++
		failedRepoNumber--
		logger.Logger.Infof("%s 已迁移，跳过同步", repoPath)
		return nil
	}
	logger.Logger.Infof("%s 开始迁移", repoPath)
	startTime := time.Now()
	isSvn := git.IsSvnRepo(depot.GetRepoType())
	if isSvn {
		skipRepoNumber++
		failedRepoNumber--
		logger.Logger.Infof("%s svn仓库，跳过同步", repoPath)
		return nil
	}
	cnbRepoPath := cnb.GetCnbRepoPath(subGroupName, repoName, organizationMappingLevel)
	has, err := cnb.HasRepoV2(CnbApiURL, CnbToken, subGroupName, repoName, organizationMappingLevel)
	if err != nil {
		return err
	}

	if !has {
		err = cnb.CreateRepo(CnbApiURL, CnbToken, subGroupName, repoName, organizationMappingLevel, repoPrivate)
		if err != nil {

			return fmt.Errorf("%s 仓库创建失败: %s", repoPath, err)
		}
		logger.Logger.Infof("%s 仓库创建成功", repoPath)
	} else if has && SkipExistsRepo {
		skipRepoNumber++
		failedRepoNumber--
		logger.Logger.Warnf("%s CNB仓库%s已存在，跳过同步", repoPath, cnbRepoPath)
		return nil
	}
	//err = system.CreateDirIfNotExists(depot.ProjectName)
	//if err != nil {
	//	return fmt.Errorf("创建%s目录失败", depot.Name)
	//}
	//cloneURL := depot.GetCloneUrl()
	err = depot.Clone()
	if err != nil {
		logger.Logger.Errorf(err.Error())
		return fmt.Errorf(err.Error())
	}

	// 设置要进入的目录路径
	pwdDir, err := os.Getwd()
	if err != nil {
		return err
	}
	// 完整的仓库目录路径
	fullRepoDir := filepath.Join(pwdDir, repoPath)
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			logger.Logger.Errorf("%s 删除失败: %s", path, err)
		}
	}(fullRepoDir)

	pushURL := cnb.GetPushUrl(organizationMappingLevel, CnbURL, CnbUserName, CnbToken, subGroupName, repoName)

	output, err := git.Push(repoPath, pushURL, IsForcePush)
	if err != nil && useLfsMigrate && git.IsExceededLimitError(output) {
		logger.Logger.Warnf("%s 历史提交文件大小超过%sM", repoPath, git.FileLimitSize)
		fixError := git.FixExceededLimitError(repoPath)
		if fixError != nil {
			return fmt.Errorf("%s 修复大文件超过限制: %s", repoPath, fixError)
		}
		output, err = git.Push(repoPath, pushURL, IsForcePush)
		if err != nil {
			return fmt.Errorf("%s push失败: %s\n %s", repoPath, err, output)
		}
	}
	if err != nil {
		return fmt.Errorf("%s push失败: %s\n %s", repoPath, err, output)
	}

	err, IsLFSRepo := git.IsLFSRepo(repoPath)
	if err != nil {
		return fmt.Errorf("%s 判断是否是LFS仓库失败: %s", repoPath, err)
	}
	if IsLFSRepo {
		out, err := git.FetchLFS(repoPath, allowIncompletePush)
		if err != nil {
			return fmt.Errorf("%s 下载LFS文件失败: %s\n%s", repoPath, err, out)
		}
		out, err = git.PushLFS(repoPath, pushURL)
		if err != nil {
			return fmt.Errorf("%s 上传LFS文件失败: %s\n%s", repoPath, err, out)
		}
	}

	successfulRepoNumber++
	failedRepoNumber--
	duration := time.Since(startTime)
	logger.Logger.Infof("%s 迁移至CNB %s 成功,耗时%s", repoPath, cnbRepoPath, duration)
	logger.RecordSuccessfulRepo(repoPath)
	return nil
}

func isMigrated(repoPath, filePath string) (error, bool) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err, false
	}
	//pattern := fmt.Sprintf(`(?<=^|[\s])%s(?=$|[\s])`, regexp.QuoteMeta(repoPath))
	pattern := fmt.Sprintf(`(?:^|[\s])%s(?:$|[\s])`, regexp.QuoteMeta(repoPath))
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("正则表达式编译失败: %s", err), false
	}
	return nil, re.MatchString(string(content))
}
