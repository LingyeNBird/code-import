package git

import (
	"ccrctl/pkg/api/coding"
	"ccrctl/pkg/config"
	"ccrctl/pkg/logger"
	"ccrctl/pkg/system"
	"fmt"
	"strings"
)

const (
	CnbOriginName    = "cnb"
	SourceOriginName = "source"
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

func Rebase(repoPath, cloneURL string) error {
	logger.Logger.Infof("%s 开始rebase", repoPath)
	logger.Logger.Debugf("git rebase %s %s", repoPath, cloneURL)
	out, err := system.RunCommand("git", repoPath, "remote", "add", SourceOriginName, cloneURL)
	if err != nil {
		return fmt.Errorf("%s 添加source远程仓库失败: %s\n %s", repoPath, err, out)
	}
	logger.Logger.Infof("%s 添加source远程仓库成功", repoPath)
	out, err = system.RunCommand("git", repoPath, "fetch", SourceOriginName)
	if err != nil {
		return fmt.Errorf("%s 拉取souce远程仓库失败: %s\n %s", repoPath, err, out)
	}
	logger.Logger.Infof("%s 拉取souce远程仓库成功", repoPath)

	defaultBranchName, err := system.RunCommand("git", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return fmt.Errorf("%s 获取当前分支失败: %s\n %s", repoPath, err, out)
	}
	defaultBranchName = strings.TrimSpace(defaultBranchName)
	rebaseBranch := SourceOriginName + "/" + defaultBranchName
	logger.Logger.Debugf("%s 当前默认分支: %s", repoPath, rebaseBranch)
	out, err = system.RunCommand("git", repoPath, "rebase", rebaseBranch)
	if err != nil {
		return fmt.Errorf("%s rebase失败: %s\n %s", repoPath, err, out)
	}
	logger.Logger.Infof("%s rebase成功", repoPath)
	pushOut, err := system.RunCommand("git", repoPath, "push", "-f")
	if err != nil {
		return fmt.Errorf("%s rebase后 push失败: %s\n %s", repoPath, err, pushOut)
	}
	logger.Logger.Infof("%s push成功", repoPath)
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
	logger.Logger.Infof("%s push成功", repoPath)
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
	logger.Logger.Infof("%s 上传LFS文件", repoPath)
	workDir := repoPath
	output, err := system.RunCommand("git", workDir, "lfs", "push", "--all", pushUrl)
	if err != nil {
		logger.Logger.Errorf("%s 上传LFS文件失败", repoPath)
		return output, err
	}
	logger.Logger.Infof("%s 上传LFS文件成功", repoPath)
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
