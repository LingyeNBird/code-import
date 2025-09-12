package migrate

import (
	"ccrctl/pkg/api/target"
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
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	GitDirName      = "source_git_dir"
	CnbUserName     = "cnb"
	MaxConcurrency  = 10
	RebaseDirPrefix = "rebase"
	RepoPathFile    = "repo-path.txt"
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
	DownloadOnly             = config.Cfg.GetBool("migrate.download_only")
	RootGroupName            = config.Cfg.GetString("cnb.root_organization")
	rebaseBackDirPath        string
	rebaseBranchesMap        sync.Map
	workDirCreated           bool
)

// checkAndGetRepoList 检查并获取仓库列表
// 如果启用了仓库选择功能且 repo-path.txt 不存在，则获取仓库列表并写入文件
// 返回是否需要继续执行迁移
func checkAndGetRepoList(source vcs.VCS) (bool, error) {
	if !config.Cfg.GetBool("migrate.allow_select_repos") {
		return true, nil
	}

	if _, err := os.Stat(RepoPathFile); err == nil {
		return true, nil
	}

	logger.Logger.Info("已启用仓库选择功能，正在获取仓库列表...")
	_, err := GetRepoList(source)
	if err != nil {
		return false, fmt.Errorf("获取仓库列表失败: %s", err)
	}
	logger.Logger.Infof("仓库列表已写入 %s，请编辑该文件保留需要迁移的仓库，然后重新运行迁移工具", RepoPathFile)
	return false, nil
}

// filterReposBySelection 根据 repo-path.txt 过滤仓库列表
func filterReposBySelection(depotList []vcs.VCS) ([]vcs.VCS, error) {
	if !config.Cfg.GetBool("migrate.allow_select_repos") {
		return depotList, nil
	}

	selectedRepos, err := ReadSelectedRepos()
	if err != nil {
		return nil, fmt.Errorf("读取仓库列表文件失败: %s", err)
	}

	filteredDepotList := make([]vcs.VCS, 0, len(depotList))
	for _, depot := range depotList {
		if selectedRepos[depot.GetRepoPath()] {
			filteredDepotList = append(filteredDepotList, depot)
		} else {
			logger.Logger.Infof("跳过仓库 %s（未在 %s 中选择）", depot.GetRepoPath(), RepoPathFile)
		}
	}
	return filteredDepotList, nil
}

// initMigrationStats 初始化迁移统计信息
func initMigrationStats(depotList []vcs.VCS) {
	atomic.StoreInt64(&totalRepoNumber, int64(len(depotList)))
	atomic.StoreInt64(&failedRepoNumber, int64(len(depotList)))
	atomic.StoreInt64(&successfulRepoNumber, 0)
	atomic.StoreInt64(&skipRepoNumber, 0)
}

func Run() {
	startTime := time.Now()     // 记录迁移开始时间
	err := config.CheckConfig() // 检查配置文件
	if err != nil {
		logger.Logger.Errorf("配置文件校验失败: %s", err)
		return
	}
	logger.Logger.Infof("SOURCE_URL: %s", config.Cfg.GetString("source.url"))
	logger.Logger.Infof("CNB_URL: %s", config.Cfg.GetString("cnb.url"))
	err = system.SetFileDescriptorLimit(system.Limit) // 设置文件描述符限制
	if err != nil {
		logger.Logger.Errorf("设置文件描述符限制失败: %s", err)

		return
	}

	// 获取源平台的 VCS 实例列表
	sourceVcsList, err := vcs.NewVcs(SourcePlatformName)
	if err != nil {
		logger.Logger.Errorf("获取源平台仓库列表失败，请检查配置参数: %s", err)
		return
	}
	if len(sourceVcsList) == 0 {
		logger.Logger.Warnf("源平台仓库列表为空，无需迁移")
		return
	}
	sourceVcs := sourceVcsList[0] // 使用第一个实例

	// 检查是否需要获取仓库列表
	shouldContinue, err := checkAndGetRepoList(sourceVcs)
	if err != nil {
		logger.Logger.Errorf("%s", err)
		return
	}
	if !shouldContinue {
		return
	}

	// 获取并过滤仓库列表
	depotList := sourceVcsList
	depotList, err = filterReposBySelection(depotList)
	if err != nil {
		logger.Logger.Errorf("%s", err)
		return
	}

	logger.Logger.Infof("待迁移仓库总数%d", len(depotList))
	initMigrationStats(depotList)

	// 如果不是只下载模式，则执行 CNB 相关操作
	if !DownloadOnly {
		// 检查根组织
		logger.Logger.Infof("检查根组织%s是否存在", RootGroupName)
		exist, err := target.RootOrganizationExists(CnbApiURL, CnbToken)
		if err != nil {
			logger.Logger.Errorf("判断根组织是否存在失败: %s", err)
			return
		}
		if !exist {
			logger.Logger.Errorf("根组织%s不存在，请先创建根组织", RootGroupName)
			return
		}

		// 创建子组织（如果需要）
		if organizationMappingLevel == 1 {
			err = target.CreateSubOrganizationIfNotExists(CnbApiURL, CnbToken, depotList)
			if err != nil {
				logger.Logger.Errorf("创建子组织失败: %s", err)
				return
			}
		}
	}

	// 设置并发数
	if Concurrency > MaxConcurrency {
		Concurrency = MaxConcurrency
	}

	// 处理 SSH 配置
	if config.Cfg.GetBool("migrate.ssh") {
		if err := setupSSH(); err != nil {
			panic(err)
		}
	}

	// 设置工作目录
	if err := setupWorkDir(); err != nil {
		logger.Logger.Errorf("%s", err)
		return
	}

	// defer 删除 source_git_dir 目录，确保所有迁移操作完成后再清理（仅删除本工具创建的目录，且非 local 平台）
	pwdDir, err := os.Getwd()
	if err == nil && !DownloadOnly && workDirCreated && SourcePlatformName != "local" {
		gitDirABSPath := filepath.Join(pwdDir, "..", GitDirName)
		defer func(path string) {
			_ = os.RemoveAll(path)
		}(gitDirABSPath)
	}

	// 执行迁移
	executeMigration(depotList, startTime)
}

// setupSSH 设置 SSH 配置
func setupSSH() error {
	sourceKeyPath := "ssh.key"
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")

	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("创建.ssh目录失败: %v", err)
	}

	if _, err := os.Stat(sourceKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("SSH私钥文件 %s 不存在于当前目录", sourceKeyPath)
	}

	privateKeyPath := filepath.Join(sshDir, "id_rsa")
	keyData, err := os.ReadFile(sourceKeyPath)
	if err != nil {
		return fmt.Errorf("读取SSH私钥文件失败: %v", err)
	}

	if err := os.WriteFile(privateKeyPath, keyData, 0600); err != nil {
		return fmt.Errorf("复制SSH私钥文件失败: %v", err)
	}

	logger.Logger.Infof("已成功复制SSH私钥文件从 %s 到 %s", sourceKeyPath, privateKeyPath)
	return nil
}

// setupWorkDir 设置工作目录
func setupWorkDir() error {
	var workDirName string
	if DownloadOnly {
		// 在只下载模式下，使用带时间戳的目录名
		timestamp := time.Now().Format("20060102150405")
		workDirName = fmt.Sprintf("%s_%s", GitDirName, timestamp)
	} else {
		workDirName = GitDirName
	}

	if st, err := os.Stat(workDirName); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(workDirName, 0755); err != nil {
				return fmt.Errorf("创建Git工作目录失败: %s", err)
			}
			logger.Logger.Infof("创建仓库工作目录%s成功", workDirName)
			workDirCreated = true
		} else {
			return fmt.Errorf("检查Git工作目录失败: %s", err)
		}
	} else if !st.IsDir() {
		return fmt.Errorf("%s 已存在但不是目录", workDirName)
	} else {
		logger.Logger.Infof("使用已有仓库工作目录%s", workDirName)
	}

	pwdDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %s", err)
	}

	gitDirABSPath := filepath.Join(pwdDir, workDirName)

	system.HandleInterrupt(gitDirABSPath)

	if MigrateRebase {
		if err := setupRebase(pwdDir); err != nil {
			return err
		}
	}

	if err := os.Chdir(workDirName); err != nil {
		return fmt.Errorf("切换到Git工作目录失败: %s", err)
	}

	return nil
}

// setupRebase 设置 rebase 相关配置
func setupRebase(pwdDir string) error {
	if err := system.SetGlobalGitUser(); err != nil {
		return fmt.Errorf("设置全局Git用户失败: %s", err)
	}

	if err := git.SetCheckOutDefaultRemote(); err != nil {
		return fmt.Errorf("设置默认远程仓库失败: %s", err)
	}

	rebaseBackDirPath = filepath.Join(pwdDir, time.Now().Format("200601021504")+"bak")
	if err := os.Mkdir(rebaseBackDirPath, 0755); err != nil {
		return fmt.Errorf("创建rebase备份目录失败: %s", err)
	}
	logger.Logger.Infof("创建rebase备份目录%s成功", rebaseBackDirPath)

	return nil
}

// executeMigration 执行迁移操作
func executeMigration(depotList []vcs.VCS, startTime time.Time) {
	if DownloadOnly {
		logger.Logger.Infof("开始下载仓库，当前并发数:%d", Concurrency)
	} else {
		logger.Logger.Infof("开始迁移仓库，当前并发数:%d", Concurrency)
	}
	sem := semaphore.NewWeighted(int64(Concurrency))
	var wg sync.WaitGroup

	for _, depot := range depotList {
		wg.Add(1)
		depotCopy := depot
		go func(depot vcs.VCS) {
			defer wg.Done()
			if err := sem.Acquire(context.Background(), 1); err != nil {
				panic(err)
			}
			defer sem.Release(1)
			if err := migrateDo(depot); err != nil {
				logger.Logger.Errorf("%s 仓库%s失败: %s", depot.GetRepoPath(), getOperationType(), err)
			}
		}(depotCopy)
	}

	wg.Wait()
	duration := int(time.Since(startTime).Seconds())
	if DownloadOnly {
		logger.Logger.Infof("代码仓库下载完成，耗时%d秒。\n【仓库总数】%d【成功下载】%d【忽略下载】%d【下载失败】%d",
			duration, totalRepoNumber, successfulRepoNumber, skipRepoNumber, failedRepoNumber)
	} else {
		logger.Logger.Infof("代码仓库迁移完成，耗时%d秒。\n【仓库总数】%d【成功迁移】%d【忽略迁移】%d【迁移失败】%d",
			duration, totalRepoNumber, successfulRepoNumber, skipRepoNumber, failedRepoNumber)
	}
	// 检查是否有忽略迁移或迁移失败的仓库
	if skipRepoNumber > 0 || failedRepoNumber > 0 {
		logger.Logger.Errorf("存在忽略迁移或迁移失败的仓库，请检查error级别日志查看详情")
	}
}

// getOperationType 根据当前模式返回操作类型
func getOperationType() string {
	if DownloadOnly {
		return "下载"
	}
	return "迁移"
}

func migrateDo(depot vcs.VCS) error {
	var err error
	repoName, subGroup, repoPath, repoPrivate := depot.GetRepoName(), depot.GetSubGroup(), depot.GetRepoPath(), depot.GetRepoPrivate()
	subGroupName := subGroup.Name

	// 如果不是只下载模式，则检查是否已迁移
	if !DownloadOnly {
		err, migrated := isMigrated(repoPath, logger.SuccessfulLogFilePath)
		if err != nil {
			logger.Logger.Errorf("判断是否迁移失败: %s", err)
			return fmt.Errorf("%s 判断是否迁移失败%s", repoPath, err)
		}
		if migrated {
			atomic.AddInt64(&skipRepoNumber, 1)
			atomic.AddInt64(&failedRepoNumber, -1)
			logger.Logger.Infof("%s 已迁移，忽略迁移", repoPath)
			return nil
		}
	}

	logger.Logger.Infof("%s 开始迁移", repoPath)
	startTime := time.Now()
	isSvn := git.IsSvnRepo(depot.GetRepoType())
	if isSvn {
		atomic.AddInt64(&skipRepoNumber, 1)
		atomic.AddInt64(&failedRepoNumber, -1)
		logger.Logger.Errorf("%s svn仓库，忽略迁移", repoPath)
		return nil
	}
	// 执行 clone 操作
	err = depot.Clone()
	if err != nil {
		logger.Logger.Errorf(err.Error())
		return fmt.Errorf(err.Error())
	}
	// 如果是只下载模式，则直接返回
	if DownloadOnly {
		atomic.AddInt64(&successfulRepoNumber, 1)
		atomic.AddInt64(&failedRepoNumber, -1)
		duration := time.Since(startTime)
		logger.Logger.Infof("%s 下载完成，耗时%s", repoPath, duration)
		return nil
	}
	// 以下是原有的迁移逻辑
	cnbRepoPath, cnbRepoGroup := target.GetCnbRepoPathAndGroup(subGroupName, repoName, organizationMappingLevel)
	if MigrateCode {
		has, err := target.HasRepoV2(CnbApiURL, CnbToken, cnbRepoPath)
		if err != nil {
			return err
		}
		if !has {
			err = target.CreateRepo(CnbApiURL, CnbToken, cnbRepoGroup, repoName, depot.GetRepoDescription(), repoPrivate)
			if err != nil {
				return fmt.Errorf("%s 仓库创建失败: %s", repoPath, err)
			}
			logger.Logger.Infof("%s 仓库创建成功", repoPath)
			time.Sleep(1000 * time.Millisecond) // 添加0.5秒延迟，避免push操作太快导致报错找不到仓库
		} else if has && SkipExistsRepo {
			atomic.AddInt64(&skipRepoNumber, 1)
			atomic.AddInt64(&failedRepoNumber, -1)
			logger.Logger.Warnf("%s CNB仓库%s已存在，忽略迁移", repoPath, cnbRepoPath)
			return nil
		}
		// 检查源仓库是否初始化
		if !git.IsBareRepoInitialized(repoPath) {
			atomic.AddInt64(&successfulRepoNumber, 1)
			atomic.AddInt64(&failedRepoNumber, -1)
			logger.Logger.Infof("%s 源仓库未初始化", repoPath)
			return nil
		}
		// 设置要进入的目录路径
		pwdDir, err := os.Getwd()
		if err != nil {
			return err
		}
		// 完整的仓库目录路径
		fullRepoDir := filepath.Join(pwdDir, repoPath)
		if SourcePlatformName != "local" {
			defer func(path string) {
				err := os.RemoveAll(path)
				if err != nil {
					logger.Logger.Errorf("%s 删除失败: %s", path, err)
				}
			}(fullRepoDir)
		}

		pushURL := target.GetPushUrl(organizationMappingLevel, CnbURL, CnbUserName, CnbToken, subGroupName, repoName)
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
				logger.Logger.Errorf("备份仓库失败: %v", err)
				return fmt.Errorf("备份仓库失败: %w", err)
			}
			logger.Logger.Infof("%s 已备份仓库到 %s", repoPath, destPath)
			rebaseErr := git.Rebase(rebaseRepoPath, depot.GetRepoPath())
			if rebaseErr != nil {
				return rebaseErr
			}
		}

		output, err := git.Push(repoPath, pushURL, isForcePush)
		if err != nil && useLfsMigrate && git.IsExceededLimitError(output) {
			logger.Logger.Warnf("%s 历史提交文件大小超过%sM", repoPath, git.FileLimitSize)
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
	duration := int(time.Since(startTime).Seconds())
	logger.Logger.Infof("%s 迁移至CNB %s 成功,耗时%d秒", repoPath, cnbRepoPath, duration)
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
	if SourcePlatformName == "common" {
		return nil
	}
	releases := depot.GetReleases()
	if releases == nil || 0 == len(releases) {
		logger.Logger.Infof("%s 无release需要迁移", depot.GetRepoPath())
		return nil
	}
	repoPath := depot.GetRepoPath()

	logger.Logger.Infof("%s 开始迁移 release", repoPath)

	// 遍历处理每个release
	for _, release := range releases {
		if err := migrateOneRelease(depot, release, repoPath); err != nil {
			return err
		}
	}

	logger.Logger.Infof("%s 迁移 release 成功", repoPath)
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
	logger.Logger.Infof("%s 开始迁移release: %s", repoPath, release.Name)

	// 在目标平台创建release
	releaseID, exist, err := target.CreateRelease(repoPath, depot.GetProjectID(), release, depot)
	if err != nil {
		logger.Logger.Errorf("%s 迁移 release %s 失败: %s", repoPath, release.Name, err)
		return err
	}

	// 如果release已存在则跳过
	if exist {
		logger.Logger.Warnf("%s 迁移release: %s 已存在，忽略迁移", repoPath, release.Name)
		return nil
	}

	// 处理release附带的资源文件
	if len(release.Assets) > 0 {
		if err := migrateReleaseAssets(repoPath, releaseID, release); err != nil {
			return err
		}
	}

	logger.Logger.Infof("%s 迁移 release %s 成功", repoPath, release.Name)
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
	// 遍历处理每个资源文件
	for _, asset := range release.Assets {

		if err := migrateReleaseAsset(repoPath, releaseID, asset.Name, asset.Url); err != nil {
			logger.Logger.Errorf("%s 迁移 release %s asset %s 失败: %s",
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
	err = target.UploadReleaseAsset(repoPath, releaseID, fileName, data)
	if err != nil {
		logger.Logger.Errorf("%s 上传release asset %s 失败: %s", downloadUrl, fileName, err)
		return err
	}
	return nil
}

// GetRepoList 获取源平台仓库列表
func GetRepoList(source vcs.VCS) ([]string, error) {
	repos, err := source.ListRepos()
	if err != nil {
		return nil, fmt.Errorf("获取仓库列表失败: %v", err)
	}

	var repoPaths []string
	for _, repo := range repos {
		repoPaths = append(repoPaths, repo.GetRepoPath())
	}

	// 将仓库列表写入文件
	if err := os.WriteFile(RepoPathFile, []byte(strings.Join(repoPaths, "\n")), 0644); err != nil {
		return nil, fmt.Errorf("写入仓库列表文件失败: %v", err)
	}

	return repoPaths, nil
}

// ReadSelectedRepos 读取用户选择的仓库列表
func ReadSelectedRepos() (map[string]bool, error) {
	content, err := os.ReadFile(RepoPathFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("仓库列表文件 %s 不存在，请先运行工具获取仓库列表", RepoPathFile)
		}
		return nil, fmt.Errorf("读取仓库列表文件失败: %v", err)
	}

	selectedRepos := make(map[string]bool)
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			selectedRepos[line] = true
		}
	}

	return selectedRepos, nil
}
