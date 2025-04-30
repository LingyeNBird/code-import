package migrate

import (
	"ccrctl/pkg/api/source"
	"ccrctl/pkg/config"
	"ccrctl/pkg/git"
	"ccrctl/pkg/http_client"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/system"
	"ccrctl/pkg/util"
	"ccrctl/pkg/vcs"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

const (
	GitDirName      = "source_git_dir"
	CnbUserName     = "cnb"
	MaxConcurrency  = 10
	RebaseDirPrefix = "rebase"
)

var (
	CnbURL                   = config.Cfg.GetString("cnb.url")
	CnbApiURL                = config.ConvertToApiURL(CnbURL)
	CnbToken                 = config.Cfg.GetString("cnb.token")
	Concurrency              = config.Cfg.GetInt("migrate.concurrency")
	totalRepoNumber          int64
	skipRepoNumber           int64
	successfulRepoNumber     int64
	failedRepoNumber         int64
	useLfsMigrate            = config.Cfg.GetBool("migrate.use_lfs_migrate")
	organizationMappingLevel = config.Cfg.GetInt("migrate.organization_mapping_level")
	SourcePlatformName       = config.Cfg.GetString("source.platform")
	SkipExistsRepo           = config.Cfg.GetBool("migrate.skip_exists_repo")
	MigrateRelease           = config.Cfg.GetBool("migrate.release")
	MigrateCode              = config.Cfg.GetBool("migrate.code")
	MigrateRebase            = config.Cfg.GetBool("migrate.rebase")
	rebaseBackDirPath        string
	rebaseBranchesMap        sync.Map
)

func Run() {
	startTime := time.Now()     // 记录迁移开始时间
	err := config.CheckConfig() // 检查配置文件
	if err != nil {
		panic(err)
	}
	err = system.SetFileDescriptorLimit(system.Limit) // 设置文件描述符限制
	if err != nil {
		panic(err)
	}

	if config.Cfg.GetBool("migrate.ssh") {
		// 处理SSH私钥文件
		sourceKeyPath := "ssh.key" // 当前目录下的私钥文件
		sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")

		// 创建.ssh目录
		if err := os.MkdirAll(sshDir, 0700); err != nil {
			panic(fmt.Sprintf("创建.ssh目录失败: %v", err))
		}

		// 检查源文件是否存在
		if _, err := os.Stat(sourceKeyPath); os.IsNotExist(err) {
			panic(fmt.Sprintf("SSH私钥文件 %s 不存在于当前目录", sourceKeyPath))
		}

		// 设置目标路径
		privateKeyPath := filepath.Join(sshDir, "id_rsa")

		// 复制文件内容
		keyData, err := os.ReadFile(sourceKeyPath)
		if err != nil {
			panic(fmt.Sprintf("读取SSH私钥文件失败: %v", err))
		}
		logger.Logger.Debugf("SSH私钥内容: \n %s", keyData)
		// 写入目标文件
		if err := os.WriteFile(privateKeyPath, keyData, 0600); err != nil {
			panic(fmt.Sprintf("复制SSH私钥文件失败: %v", err))
		}

		privateKeyData, err := os.ReadFile(privateKeyPath)
		if err != nil {
			panic(fmt.Sprintf("读取SSH私钥文件失败: %v", err))
		}
		logger.Logger.Debugf("SSH私钥内容: \n %s", privateKeyData)

		output, err := exec.Command("sh", "-c", "ls -l ~/.ssh/").CombinedOutput()
		if err != nil {
			panic(fmt.Sprintf("ls -l 失败: %v\n命令输出: %s", err, string(output)))
		}
		logger.Logger.Debug("ls -l: %s", string(output))

		logger.Logger.Infof("已成功复制SSH私钥文件从 %s 到 %s", sourceKeyPath, privateKeyPath)
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
	if MigrateRebase {
		err = system.SetGlobalGitUser()
		if err != nil {
			panic(err)
		}
		err = git.SetCheckOutDefaultRemote()
		if err != nil {
			panic(err)
		}
		// 创建rebase备份目录
		rebaseBackDirPath = filepath.Join(pwdDir, time.Now().Format("200601021504")+"bak")
		err := os.Mkdir(rebaseBackDirPath, 0755)
		if err != nil {
			logger.Logger.Errorf("创建rebase备份目录失败: %s", err)
			panic(err)
		}
		logger.Logger.Infof("创建rebase备份目录%s成功", rebaseBackDirPath)

	}
	err = os.Chdir(GitDirName) // 切换到Git工作目录
	if err != nil {
		panic(err)
	}
	depotList, err := vcs.NewVcs(SourcePlatformName) // 获取仓库列表
	if err != nil {
		panic(err)
	}
	logger.Logger.Infof("仓库总数%d", len(depotList))
	atomic.StoreInt64(&totalRepoNumber, int64(len(depotList)))
	atomic.StoreInt64(&failedRepoNumber, int64(len(depotList)))
	atomic.StoreInt64(&successfulRepoNumber, 0)
	atomic.StoreInt64(&skipRepoNumber, 0)
	//err = cnb.CreateRootOrganizationIfNotExists(CnbApiURL, CnbToken) // 创建根组织
	//if err != nil {
	//	panic(err)
	//}
	exist, err := source.RootOrganizationExists(CnbApiURL, CnbToken)
	if err != nil {
		panic(err)
	}
	if !exist {
		logger.Logger.Errorf("根组织%s不存在，请先创建根组织", config.Cfg.GetString("cnb.root_organization"))
		return
	}
	if organizationMappingLevel == 1 { // 如果组织映射级别为1，则创建子组织
		err = source.CreateSubOrganizationIfNotExists(CnbApiURL, CnbToken, depotList)
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
		wg.Add(1) // 增加WaitGroup计数
		// 创建depot的副本
		depotCopy := depot
		go func(depot vcs.VCS) { // 启动goroutine进行迁移
			defer wg.Done() // 完成后减少WaitGroup计数

			if err := sem.Acquire(context.Background(), 1); err != nil { // 获取信号量
				panic(err)
			}
			defer sem.Release(1)    // 释放信号量
			err := migrateDo(depot) // 执行迁移操作
			if err != nil {
				logger.Logger.Errorf("%s 仓库迁移失败: %s", depot.GetRepoPath(), err) // 记录迁移失败信息
			}
		}(depotCopy)
	}
	wg.Wait()                                                                                                                                                  // 等待所有goroutine完成
	duration := time.Since(startTime)                                                                                                                          // 计算迁移耗时
	logger.Logger.Infof("代码仓库迁移完成，耗时%s。\n【仓库总数】%d【成功迁移】%d【忽略迁移】%d【迁移失败】%d", duration, totalRepoNumber, successfulRepoNumber, skipRepoNumber, failedRepoNumber) // 记录迁移完成信息
}

func migrateDo(depot vcs.VCS) error {
	var err error
	repoName, subGroup, repoPath, repoPrivate := depot.GetRepoName(), depot.GetSubGroup(), depot.GetRepoPath(), depot.GetRepoPrivate()
	subGroupName := subGroup.Name
	// 使用zap的With方法添加repo字段
	log := logger.Logger.With(zap.String("repo", repoPath))

	err, migrated := isMigrated(repoPath, logger.SuccessfulLogFilePath)
	if err != nil {
		log.Errorf("判断是否迁移失败: %s", err)
		return fmt.Errorf("%s 判断是否迁移失败%s", repoPath, err)
	}
	if migrated {
		atomic.AddInt64(&skipRepoNumber, 1)
		atomic.AddInt64(&failedRepoNumber, -1)
		log.Infof("%s 已迁移，跳过同步", repoPath)
		return nil
	}
	log.Infof("%s 开始迁移", repoPath)
	startTime := time.Now()
	isSvn := git.IsSvnRepo(depot.GetRepoType())
	if isSvn {
		atomic.AddInt64(&skipRepoNumber, 1)
		atomic.AddInt64(&failedRepoNumber, -1)
		log.Infof("%s svn仓库，跳过同步", repoPath)
		return nil
	}
	cnbRepoPath, cnbRepoGroup := source.GetCnbRepoPathAndGroup(subGroupName, repoName, organizationMappingLevel)
	if MigrateCode {
		err = depot.Clone()
		if err != nil {
			log.Errorf(err.Error())
			return fmt.Errorf(err.Error())
		}
		has, err := source.HasRepoV2(CnbApiURL, CnbToken, cnbRepoPath)
		if err != nil {
			return err
		}
		if !has {
			err = source.CreateRepo(CnbApiURL, CnbToken, cnbRepoGroup, repoName, depot.GetRepoDescription(), repoPrivate)
			if err != nil {
				return fmt.Errorf("%s 仓库创建失败: %s", repoPath, err)
			}
			log.Infof("%s 仓库创建成功", repoPath)
		} else if has && SkipExistsRepo {
			atomic.AddInt64(&skipRepoNumber, 1)
			atomic.AddInt64(&failedRepoNumber, -1)
			log.Warnf("%s CNB仓库%s已存在，跳过同步", repoPath, cnbRepoPath)
			return nil
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
				log.Errorf("%s 删除失败: %s", path, err)
			}
		}(fullRepoDir)

		pushURL := source.GetPushUrl(organizationMappingLevel, CnbURL, CnbUserName, CnbToken, subGroupName, repoName)
		rebaseRepoPath := filepath.Join(RebaseDirPrefix, repoPath)
		isForcePush := config.Cfg.GetBool("migrate.force_push")
		if MigrateRebase {
			// 如果使用rebase同步，那么需要开启强制 push，避免出现冲突
			isForcePush = true
			rebaseCloneErr := git.NormalClone(pushURL, rebaseRepoPath)
			if rebaseCloneErr != nil {
				return fmt.Errorf("git rebase clone失败: %s", rebaseCloneErr)
			}
			destPath := filepath.Join(rebaseBackDirPath, repoPath)
			//备份CNB侧仓库
			err = util.CopyDir(rebaseRepoPath, destPath)
			if err != nil {
				log.Errorf("备份仓库失败: %v", err)
				return fmt.Errorf("备份仓库失败: %w", err)
			}
			log.Infof("%s 已备份仓库到 %s", repoPath, destPath)
			rebaseErr := git.Rebase(rebaseRepoPath, depot.GetRepoPath())
			if rebaseErr != nil {
				return rebaseErr
			}
		}

		output, err := git.Push(repoPath, pushURL, isForcePush)
		if err != nil && useLfsMigrate && git.IsExceededLimitError(output) {
			log.Warnf("%s 历史提交文件大小超过%sM", repoPath, git.FileLimitSize)
			fixError := git.FixExceededLimitError(repoPath)
			if fixError != nil {
				return fmt.Errorf("%s 修复大文件超过限制: %s", repoPath, fixError)
			}
			output, err = git.Push(repoPath, pushURL, isForcePush)
			if err != nil {
				return fmt.Errorf("%s push失败: %s\n %s", repoPath, err, output)
			}
		}
		if err != nil {
			return fmt.Errorf("%s push失败: %s\n %s", repoPath, err, output)
		}
	}
	if MigrateRelease {
		err = migrateRelease(depot)
		if err != nil {
			return err
		}
	}
	atomic.AddInt64(&successfulRepoNumber, 1)
	atomic.AddInt64(&failedRepoNumber, -1)
	duration := time.Since(startTime)
	log.Infof("%s 迁移至CNB %s 成功,耗时%s", repoPath, cnbRepoPath, duration)
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

// migrateRelease 处理仓库的所有release迁移
// 参数:
//   - depot: VCS接口实例，包含仓库信息
//
// 返回:
//   - error: 迁移过程中的错误信息
func migrateRelease(depot vcs.VCS) error {
	releases := depot.GetReleases()
	if releases == nil || 0 == len(releases) {
		logger.Logger.Infof("%s 无release需要迁移", depot.GetRepoPath())
		return nil
	}
	repoPath := depot.GetRepoPath()
	log := logger.Logger

	log.Infof("%s 开始迁移 release", repoPath)

	// 遍历处理每个release
	for _, release := range releases {
		if err := migrateOneRelease(depot, release, repoPath); err != nil {
			return err
		}
	}

	log.Infof("%s 迁移 release 成功", repoPath)
	return nil
}

// migrateOneRelease 处理单个release的迁移
// 参数:
//   - depot: VCS接口实例
//   - release: 需要迁移的release信息
//   - repoPath: 仓库路径
//
// 返回:
//   - error: 迁移过程中的错误信息
func migrateOneRelease(depot vcs.VCS, release vcs.Releases, repoPath string) error {
	log := logger.Logger
	log.Infof("%s 开始迁移release: %s", repoPath, release.Name)

	// 在目标平台创建release
	releaseID, exist, err := source.CreateRelease(repoPath, depot.GetProjectID(), release, depot)
	if err != nil {
		log.Errorf("%s 迁移 release %s 失败: %s", repoPath, release.Name, err)
		return err
	}

	// 如果release已存在则跳过
	if exist {
		log.Warnf("%s 迁移release: %s 已存在，跳过迁移", repoPath, release.Name)
		return nil
	}

	// 处理release附带的资源文件
	if len(release.Assets) > 0 {
		if err := migrateReleaseAssets(repoPath, releaseID, release); err != nil {
			return err
		}
	}

	log.Infof("%s 迁移 release %s 成功", repoPath, release.Name)
	return nil
}

// migrateReleaseAssets 处理release相关的资源文件迁移
// 参数:
//   - repoPath: 仓库路径
//   - releaseID: release的唯一标识
//   - release: release信息
//
// 返回:
//   - error: 迁移过程中的错误信息
func migrateReleaseAssets(repoPath, releaseID string, release vcs.Releases) error {
	log := logger.Logger
	// 遍历处理每个资源文件
	for _, asset := range release.Assets {
		fileName, err := util.GetFileNameFromURL(asset.Url)
		if err != nil {
			log.Errorf("%s 获取release asset 文件名失败: %s", asset.Url, err)
			return err
		}

		if err := migrateReleaseAsset(repoPath, releaseID, fileName, asset.Url); err != nil {
			log.Errorf("%s 迁移 release %s asset %s 失败: %s",
				repoPath, release.Name, asset.Name, err)
			return err
		}
	}
	return nil
}

func migrateReleaseAsset(repoPath, releaseID, fileName, downloadUrl string) (err error) {
	data, err := http_client.DownloadFromUrl(downloadUrl)
	if err != nil {
		logger.Logger.Errorf("%s 下载release asset %s 失败: %s", downloadUrl, fileName, err)
		return err
	}
	err = source.UploadReleaseAsset(repoPath, releaseID, fileName, data)
	if err != nil {
		logger.Logger.Errorf("%s 上传release asset %s 失败: %s", downloadUrl, fileName, err)
		return err
	}
	return nil
}
