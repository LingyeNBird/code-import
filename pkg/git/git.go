package git

import (
	"ccrctl/pkg/api/coding"
	"ccrctl/pkg/config"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/system"
	"fmt"
	"os"
	"path"
	"strings"
)

const (
	SourceOriginName                = "source"
	ListAllBranchesAndGrep          = "git branch -a | grep -E '(^|/)%s$'"
	CheckoutBranch                  = "git checkout %s --"
	RebaseBranch                    = "git rebase %s"
	ForcePushBranch                 = "git push -f"
	CNBYamlFileName                 = ".cnb.yml"
	SetCheckOutDefaultRemoteCommand = " git config --global checkout.defaultRemote origin"
	GetOriginBranchesCommand        = "git branch -r | grep '^  origin/'"
	GitPushToLocalBareRepo          = "git push -f source"
	ErrInvalidUpstreamEN            = "fatal: invalid upstream"
	ErrInvalidUpstreamZH            = "致命错误：无效的上游"
)

var FileLimitSize = config.Cfg.GetString("migrate.file_limit_size")

func Clone(cloneURL, repoPath string, allowIncompletePush bool) error {
	logger.Logger.Infof("%s 开始clone", repoPath)
	logger.Logger.Debugf("git clone --mirror %s %s", cloneURL, repoPath)
	out, err := system.RunCommand("git", "./", "clone", "--mirror", cloneURL, repoPath)
	if err != nil {
		return fmt.Errorf("%s clone失败: %s\n %s", repoPath, err, out)
	}
	out, err = FetchLFS(repoPath, allowIncompletePush)
	if err != nil {
		return fmt.Errorf("%s 下载LFS文件失败: %s\n %s", repoPath, err, out)
	}
	logger.Logger.Infof("%s clone成功", repoPath)
	return nil
}

func NormalClone(cloneURL, repoPath string) error {
	logger.Logger.Infof("%s 开始clone", repoPath)
	logger.Logger.Debugf("git clone  %s %s", cloneURL, repoPath)
	out, err := system.RunCommand("git", "./", "clone", cloneURL, repoPath)
	if err != nil {
		return fmt.Errorf("%s clone失败: %s\n %s", repoPath, err, out)
	}
	logger.Logger.Infof("%s clone成功", repoPath)
	return nil
}

func RebasePush(rebaseRepoPath string, rebaseSuccessBranches []string) error {
	for _, branchName := range rebaseSuccessBranches {
		checkBranchOut, CheckoutBranchErr := system.ExecCommand(fmt.Sprintf(CheckoutBranch, branchName), rebaseRepoPath)
		if CheckoutBranchErr != nil {
			logger.Logger.Errorf("%s 切换分支到%s失败: %s \n%s", rebaseRepoPath, branchName, CheckoutBranchErr, checkBranchOut)
			return CheckoutBranchErr
		}
		pushOut, err := system.ExecCommand(ForcePushBranch, rebaseRepoPath)
		if err != nil {
			return fmt.Errorf("%s %s rebase后 push失败: %s\n %s", rebaseRepoPath, branchName, err, pushOut)
		}
		logger.Logger.Infof("%s %s rebase后 push成功", rebaseRepoPath, branchName)
	}
	logger.Logger.Infof("%s rebase push成功", rebaseRepoPath)
	return nil
}

func Rebase(rebaseRepoPath, repoPath string) error {
	logger.Logger.Infof("%s 开始rebase", rebaseRepoPath)
	pwdDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("%s 获取当前目录失败: %s", rebaseRepoPath, err)
	}
	logger.Logger.Debugf("pwd: %s", pwdDir)
	bareRepoPath := path.Join(pwdDir, repoPath)
	logger.Logger.Debugf("bareRepoPath: %s", bareRepoPath)
	out, err := system.RunCommand("git", rebaseRepoPath, "remote", "add", SourceOriginName, bareRepoPath)
	if err != nil {
		return fmt.Errorf("%s 添加source远程仓库失败: %s\n %s", rebaseRepoPath, err, out)
	}
	logger.Logger.Infof("%s 添加source远程仓库成功", rebaseRepoPath)
	out, err = system.RunCommand("git", rebaseRepoPath, "fetch", SourceOriginName)
	if err != nil {
		return fmt.Errorf("%s 拉取souce远程仓库失败: %s\n %s", rebaseRepoPath, err, out)
	}
	logger.Logger.Infof("%s 拉取souce远程仓库成功", rebaseRepoPath)
	rebaseBranches := config.Cfg.GetStringSlice("migrate.rebase_branch")
	// 如果没有指定rebaseBranches, 则遍历所有分支进行处理
	if len(rebaseBranches) == 0 {
		// 获取所有远程分支
		out, err := system.ExecCommand(GetOriginBranchesCommand, rebaseRepoPath)
		if err != nil {
			return fmt.Errorf("%s 获取远程分支列表失败: %s\n %s", rebaseRepoPath, err, out)
		}

		// 解析分支列表，过滤掉 HEAD 和 origin/HEAD
		remoteBranches := strings.Split(strings.TrimSpace(out), "\n")
		var branches []string
		for _, branch := range remoteBranches {
			branch = strings.TrimSpace(branch)
			if strings.Contains(branch, "->") || strings.Contains(branch, "HEAD") {
				continue
			}
			// 去除 "origin/" 前缀
			branch = strings.TrimPrefix(branch, "origin/")
			branches = append(branches, branch)
		}
		logger.Logger.Infof("%s 分支列表: %s", rebaseRepoPath, branches)

		// 遍历所有分支进行rebase
		for _, branch := range branches {
			// 切换到指定分支
			checkBranchErr := checkoutBranch(rebaseRepoPath, branch)
			if checkBranchErr != nil {
				return fmt.Errorf("%s 分支 %s checkout失败: %s", rebaseRepoPath, branch, checkBranchErr)
			}
			// 检查 .cnb.yaml文件是否存在
			CNBYamlFileAbsPath := path.Join(rebaseRepoPath, CNBYamlFileName)
			exist := system.FileExists(CNBYamlFileAbsPath)
			if !exist {
				logger.Logger.Infof("%s分支%s .cnb.yml文件不存在,跳过 rebease", rebaseRepoPath, branch)
				continue
			}
			rebaseBranch := SourceOriginName + "/" + branch
			// rebase指定分支
			rebaseOut, rebaseErr := system.ExecCommand(fmt.Sprintf(RebaseBranch, rebaseBranch), rebaseRepoPath)
			if rebaseErr != nil {
				// 检查是否是分支不存在的情况
				if isInvalidUpstreamError(rebaseOut) {
					logger.Logger.Warnf("%s 分支 %s 在源仓库中不存在，跳过 rebase", rebaseRepoPath, branch)
					continue
				}
				return fmt.Errorf("分支 %s rebase失败: %s\n %s", branch, rebaseErr.Error(), rebaseOut)
			}
			logger.Logger.Infof("%s %s rebase成功", rebaseRepoPath, rebaseBranch)
			rebasePushOut, rebasePushErr := system.ExecCommand(GitPushToLocalBareRepo, rebaseRepoPath)
			if rebasePushErr != nil {
				return fmt.Errorf("分支 %s rebase后push失败: %s\n %s", branch, rebasePushErr.Error(), rebasePushOut)
			}
			logger.Logger.Infof("%s %s rebase后push成功", rebaseRepoPath, rebaseBranch)
		}
	} else {
		for _, branch := range rebaseBranches {
			// 检查分支是否存在
			_, grepBranchErr := system.ExecCommand(fmt.Sprintf(ListAllBranchesAndGrep, branch), rebaseRepoPath)
			if grepBranchErr != nil {
				logger.Logger.Debugf("%s 分支: %s 不存在,跳过 rebase", rebaseRepoPath, branch)
				continue
			}
			// 切换到指定分支
			checkBranchErr := checkoutBranch(rebaseRepoPath, branch)
			if checkBranchErr != nil {
				return fmt.Errorf("%s 分支 %s checkout失败: %s", rebaseRepoPath, branch, checkBranchErr)
			}
			//检查 .cnb.yaml文件是否存在
			CNBYamlFileAbsPath := path.Join(rebaseRepoPath, CNBYamlFileName)
			exist := system.FileExists(CNBYamlFileAbsPath)
			if !exist {
				logger.Logger.Infof("%s分支%s .cnb.yml文件不存在,跳过 rebease", rebaseRepoPath, branch)
				continue
			}
			rebaseBranch := SourceOriginName + "/" + branch
			// rebase指定分支
			rebaseOut, rebaseErr := system.ExecCommand(fmt.Sprintf(RebaseBranch, rebaseBranch), rebaseRepoPath)
			if rebaseErr != nil {
				// 检查是否是分支不存在的情况
				if isInvalidUpstreamError(rebaseOut) {
					logger.Logger.Warnf("%s 分支 %s 在源仓库中不存在，跳过 rebase", rebaseRepoPath, branch)
					continue
				}
				return fmt.Errorf("分支 %s rebase失败: %s\n %s", branch, rebaseErr.Error(), rebaseOut)
			}
			logger.Logger.Infof("%s rebase %s 成功", rebaseRepoPath, rebaseBranch)
			rebasePushOut, rebasePushErr := system.ExecCommand(GitPushToLocalBareRepo, rebaseRepoPath)
			if rebasePushErr != nil {
				return fmt.Errorf("分支 %s rebase后push失败: %s\n %s", branch, rebasePushErr.Error(), rebasePushOut)
			}
			logger.Logger.Infof("%s %s rebase后push成功", rebaseRepoPath, rebaseBranch)
		}
	}
	return nil
}

func Push(repoPath, pushURL string, forcePush bool) (output string, err error) {
	logger.Logger.Infof("%s 开始push", repoPath)
	if forcePush {
		out, err := ForcePush(repoPath, pushURL)
		if err != nil {
			return out, err
		}
	} else {
		out, err := NormalPush(repoPath, pushURL)
		if err != nil {
			return out, err
		}
	}
	logger.Logger.Infof("%s 裸仓push成功", repoPath)
	out, err := PushLFS(repoPath, pushURL)
	if err != nil {
		return out, err
	}
	return output, nil
}

// 强制推送
func ForcePush(workDir, pushURL string) (output string, err error) {
	logger.Logger.Debugf("git push -f %s refs/heads/*:refs/heads/* refs/tags/*:refs/tags/*", pushURL)
	output, err = system.RunCommand("git", workDir, "push", "-f", pushURL, "refs/heads/*:refs/heads/*", "refs/tags/*:refs/tags/*")
	if err != nil {
		return output, err
	}
	return output, nil
}

// 普通推送不带-f参数
func NormalPush(workDir, pushURL string) (output string, err error) {
	logger.Logger.Debugf("git push %s refs/heads/*:refs/heads/* refs/tags/*:refs/tags/*", pushURL)
	output, err = system.RunCommand("git", workDir, "push", pushURL, "refs/heads/*:refs/heads/*", "refs/tags/*:refs/tags/*")
	if err != nil {
		return output, err
	}
	return output, nil
}

func IsLFSRepo(repoPath string) (error, bool) {
	workDir := repoPath
	output, err := system.RunCommand("git", workDir, "lfs", "ls-files", "--all")
	logger.Logger.Debugf("%s 检查是否是LFS仓库\n%s", repoPath, output)
	if err != nil {
		return err, false
	}
	return nil, len(output) > 0
}

func FetchLFS(repoPath string, allowIncompletePush bool) (string, error) {
	workDir := repoPath
	logger.Logger.Infof("%s 下载LFS文件", repoPath)
	output, err := system.RunCommand("git", workDir, "lfs", "fetch", "--all", "origin")
	logger.Logger.Debugf("%s 下载LFS文件\n%s", repoPath, output)
	if err != nil && allowIncompletePush {
		logger.Logger.Warnf("%s 下载LFS文件失败,忽略报错继续执行lfs Push", repoPath)
		logger.Logger.Infof("%s 正在设置lfs.allowincompletepush为true", repoPath)
		output, err := system.RunCommand("git", workDir, "config", "lfs.allowincompletepush", "true")
		if err != nil {
			return output, err
		}
		logger.Logger.Infof("%s 设置lfs.allowincompletepush为true成功", repoPath)
		return output, nil
	}
	if err != nil {
		logger.Logger.Errorf("%s 下载LFS文件失败", repoPath)
		return output, err
	}
	return output, err
}

func PushLFS(repoPath, pushUrl string) (string, error) {
	logger.Logger.Infof("%s 开始推送LFS文件", repoPath)
	workDir := repoPath
	output, err := system.RunCommand("git", workDir, "lfs", "push", "--all", pushUrl)
	if err != nil {
		logger.Logger.Errorf("%s LFS文件推送失败", repoPath)
		return output, err
	}
	logger.Logger.Infof("%s LFS文件推送成功", repoPath)
	return output, err
}

func FixExceededLimitError(repoPath string) error {
	workDir := repoPath
	above := "--above=" + FileLimitSize + "Mb"
	logger.Logger.Infof("%s 使用git lfs migrate 处理历史提交中的大文件", repoPath)
	output, err := system.RunCommand("git", workDir, "lfs", "migrate", "import", "--everything", above)
	if err != nil {
		return fmt.Errorf("git lfs migrate import 失败: %s\n%s", err, output)
	}
	logger.Logger.Infof("%s 使用git lfs migrate 处理历史提交中的大文件成功", repoPath)
	logger.Logger.Debugf("%s 使用git lfs migrate 处理历史提交中的大文件成功\n%s", repoPath, output)
	return nil
}

func IsExceededLimitError(output string) bool {
	if strings.Contains(output, "exceeded limit") {
		return true
	}
	return false
}

func IsSvnRepo(vcsType string) bool {
	if vcsType == coding.SvnVcsType {
		return true
	}
	return false
}

func SetCheckOutDefaultRemote() error {
	output, err := system.ExecCommand(SetCheckOutDefaultRemoteCommand, ".")
	if err != nil {
		return fmt.Errorf("git config remote.origin.fetch 失败: %s\n%s", err, output)
	}
	return nil
}

func checkoutBranch(repoPath, branch string) error {
	_, err := system.ExecCommand(fmt.Sprintf(CheckoutBranch, branch), repoPath)
	if err != nil {
		logger.Logger.Warnf(fmt.Sprintf("%s 切换分支 %s 失败,尝试指定ref切换", repoPath, branch))
		refBranch := "refs/remotes/origin/" + branch
		// 兼容分支名匹配到 tree object 问题
		out, refErr := system.ExecCommand(fmt.Sprintf(CheckoutBranch, refBranch), repoPath)
		if refErr != nil {
			logger.Logger.Errorf(fmt.Sprintf("%s 指定ref切换分支 %s 失败: %s\n%s", repoPath, refBranch, refErr, out))
			return refErr
		}
	}
	return nil
}

// isInvalidUpstreamError 检查是否是分支不存在的错误（同时处理英文和中文环境）
func isInvalidUpstreamError(output string) bool {
	return strings.Contains(output, ErrInvalidUpstreamEN) ||
		strings.Contains(output, ErrInvalidUpstreamZH)
}

// IsBareRepoInitialized 判断本地裸仓库是否初始化（有分支或tag）
// repoDir: 本地裸仓库目录
func IsBareRepoInitialized(repoDir string) bool {
	// 执行 git for-each-ref，若有输出则说明已初始化
	output, err := system.ExecCommand("git for-each-ref", repoDir)
	if err != nil {
		logger.Logger.Warnf("检测裸仓库初始化状态失败: %s", err)
		return false
	}
	return len(output) > 0
}
