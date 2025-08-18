package vcs

import (
	"ccrctl/pkg/logger"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	LocalPlatformName = "local"
)

type LocalVcs struct {
	RepoPath string
	RepoName string
}

func (l *LocalVcs) GetRepoPath() string { return l.RepoPath }
func (l *LocalVcs) GetRepoName() string { return l.RepoName }
func (l *LocalVcs) GetSubGroup() *SubGroup {
	parts := strings.Split(l.RepoPath, "/")
	if len(parts) > 0 {
		parts = parts[:len(parts)-1]
	}
	return &SubGroup{Name: strings.Join(parts, "/")}
}
func (l *LocalVcs) GetRepoType() string     { return Git }
func (l *LocalVcs) GetCloneUrl() string     { return "" }
func (l *LocalVcs) GetUserName() string     { return "" }
func (l *LocalVcs) GetToken() string        { return "" }
func (l *LocalVcs) Clone() error            { return nil }
func (l *LocalVcs) GetRepoPrivate() bool    { return true }
func (l *LocalVcs) GetReleases() []Releases { return nil }
func (l *LocalVcs) GetProjectID() string    { return "0" }
func (l *LocalVcs) GetReleaseAttachments(desc string, repoPath string, projectID string) ([]Attachment, error) {
	return nil, nil
}
func (l *LocalVcs) GetRepoDescription() string { return "" }
func (l *LocalVcs) ListRepos() ([]VCS, error)  { return nil, nil }

// newLocalRepo scans ./source_git_dir directory and builds VCS list from local bare repos.
func newLocalRepo() ([]VCS, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取当前目录失败: %v", err)
	}
	root := filepath.Join(pwd, "source_git_dir")
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("本地目录 %s 不存在，请先将待迁移仓库放在该目录下", root)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s 不是目录", root)
	}
	var repos []VCS
	walkFn := func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() {
			return nil
		}
		// 识别裸仓库：目录下存在 HEAD 和 objects 目录
		headPath := filepath.Join(path, "HEAD")
		objectsPath := filepath.Join(path, "objects")
		if _, err := os.Stat(headPath); err == nil {
			if info, err := os.Stat(objectsPath); err == nil && info.IsDir() {
				rel, relErr := filepath.Rel(root, path)
				if relErr != nil {
					return relErr
				}
				rel = filepath.ToSlash(rel)
				name := rel
				parts := strings.Split(rel, "/")
				if len(parts) > 0 {
					name = parts[len(parts)-1]
				}
				repos = append(repos, &LocalVcs{RepoPath: rel, RepoName: name})
				logger.Logger.Infof("发现本地仓库: %s", rel)
				return filepath.SkipDir
			}
		}
		return nil
	}
	if err := filepath.WalkDir(root, walkFn); err != nil {
		return nil, err
	}
	return repos, nil
}
