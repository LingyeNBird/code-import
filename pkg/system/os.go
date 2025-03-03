package system

import (
	"ccrctl/pkg/logger"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	Limit        = uint64(65535)
	GitUserName  = "cnb"
	GitUserEmail = "cnb@cnb.cool"
)

func SetFileDescriptorLimit(limit uint64) error {
	var rLimit syscall.Rlimit

	// 获取当前文件描述符限制
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("error getting rlimit: %v", err)
	}

	// 设置新的文件描述符限制
	rLimit.Cur = limit
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("error setting rlimit: %v", err)
	}

	// 验证新的文件描述符限制已生效
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("error getting rlimit: %v", err)
	}

	return nil
}

func RunCommand(command, workDir string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func ExecCommand(command, workDir string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func HandleInterrupt(path string) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ch

		logger.Logger.Infof("收到中断信号，程序即将终止")

		err := os.RemoveAll(path)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
}

func CreateDirIfNotExists(dirPath string) error {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// 拼接完整的目录路径
	fullPath := filepath.Join(currentDir, dirPath)

	// 检查目录是否存在
	_, err = os.Stat(fullPath)
	if os.IsNotExist(err) {
		// 目录不存在,创建目录
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		// 其他错误
		return err
	}

	// 目录已存在
	return nil
}

// SetGlobalGitUser 设置全局Git用户信息
func SetGlobalGitUser() error {
	// 设置用户名
	if output, err := RunCommand("git", "", "config", "--global", "user.name", GitUserName); err != nil {
		return fmt.Errorf("设置用户名失败: %v\n输出: %s", err, output)
	}

	// 设置邮箱
	if output, err := RunCommand("git", "", "config", "--global", "user.email", GitUserEmail); err != nil {
		return fmt.Errorf("设置邮箱失败: %v\n输出: %s", err, output)
	}

	return nil
}

func FileExists(path string) bool {
	pwd, err := os.Getwd()
	if err != nil {
		return false
	}
	logger.Logger.Debugf("当前工作目录: %s", pwd)
	_, err = os.Stat(path)
	return !os.IsNotExist(err)
}
