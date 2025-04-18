package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	configName           = "config.yaml"
	defaultConfigName    = "config.yaml.default"
	defaultConfigContent = `coding:
  #CODING URL
  url: https://coding.example.com
  token: xxx
  #项目名，多个以英文逗号分割
  project: [ProjectName]
  仓库路径，团队名/项目名/仓库名，多个以英文逗号分割
  repo: [TeamName/ProjectName/RepoName]
cnb:
  #CNB URL
  url: https://cnb.example.com
  token: xxx
  #CNB根组织，需提前手动创建
  root_organization: "coding"
migrate:
  #迁移类型，project/repo
  type: project
  #仓库迁移并发数，最大10
  concurrency: 10
  #强制推送
  force_push: true
  #忽略源仓库LFS文件丢失错误，直接push代码
  Ignore_lfs_notfound_error: false
  #是否使用lfs migrate处理历史提交的大文件，业务仓库不建议开启，会造成迁移后commit ID不一致
  use_lfs_migrate: false
  #CODING与CNB组织映射关系，1表示CODING团队映射为CNB根组织，2表示CODING项目映射为CNB根组织
  organization_mapping_level: 1`
)

var Cfg *viper.Viper

type Config struct {
	Source  Source  `yaml:"source"`
	CNB     CNB     `yaml:"cnb"`
	Migrate Migrate `yaml:"migrate"`
}

type Source struct {
	URL           string   `yaml:"url"`
	Token         string   `yaml:"token"`
	Project       []string `yaml:"project"`
	Repo          []string `yaml:"repo"`
	Platform      string   `yaml:"platform"`
	UserName      string   `yaml:"username"`
	Password      string   `yaml:"password"`
	SshPrivateKey string   `yaml:"ssh_private_key"`
}

type CNB struct {
	URL              string `yaml:"url"`
	Token            string `yaml:"token"`
	RootOrganization string `yaml:"root_organization"`
}

type Migrate struct {
	Type                     string `yaml:"type"`
	Concurrency              int    `yaml:"concurrency"`
	ForcePush                bool   `yaml:"force_push"`
	IgnoreLFSNotFoundError   bool   `yaml:"ignore_lfs_not_found_error"`
	useLfsMigrate            bool   `yaml:"use_lfs_migrate"`
	organizationMappingLevel int    `yaml:"organization_mapping_level"`
	//针对LFS源文件丢失的仓库，LFS推送时，允许本地缓存中丢失对象，而无需停止 Git 推送。https://github.com/git-lfs/git-lfs/blob/main/docs/man/git-lfs-config.adoc
	allowIncompletePush bool   `yaml:"allow_incomplete_push"`
	LogLevel            string `yaml:"log_level"`
	fileLimitSize       int64  `yaml:"file_limit_size"`
	SkipExistsRepo      bool   `yaml:"skip_exists_repo"`
	Ssh                 bool   `yaml:"ssh"`
}

func CheckConfig() error {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	// 解析 YAML 文件
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	platform := config.Source.Platform
	if platform != "common" && platform != "coding" && platform != "gitlab" && platform != "github" && platform != "gitee" && platform != "aliyun" {
		return fmt.Errorf("source.platform error only support common、coding、gitlab、github、gitee、aliyun")
	}
	if platform != "aliyun" {
		// 检查 source 参数
		if config.Source.URL == "" {
			return fmt.Errorf("source.url is required")
		}

		// 检查 source.url 前缀
		if !strings.HasPrefix(config.Source.URL, "http://") && !strings.HasPrefix(config.Source.URL, "https://") {
			return fmt.Errorf("source.url must start with 'http://' or 'https://'")
		}
	}

	//非阿里云平台和通用第三方平台迁移，检查 source.token 参数
	if platform != "common" && platform != "aliyun" {
		if config.Source.Token == "" {
			return fmt.Errorf("source.token is required")
		}
	}

	//common http迁移
	if platform == "common" && !config.Migrate.Ssh {
		if config.Source.UserName == "" || config.Source.Password == "" {
			return fmt.Errorf("when platform is common, source.username、password is required")
		}
		if platform == "common" {
			if len(config.Source.Repo) == 0 || config.Source.Repo[0] == "" {
				return fmt.Errorf("when platform is common, source.repo is required")
			}
		}
	}

	// 检查 migrate 参数
	if config.Migrate.Type == "" {
		return fmt.Errorf("migrate.type is required")
	}

	if config.Migrate.Type != "repo" && config.Migrate.Type != "project" && config.Migrate.Type != "team" {
		return fmt.Errorf("migrate.type error only support repo or project or team")
	}

	if config.Migrate.Type == "project" && (len(config.Source.Project) == 0 || config.Source.Project[0] == "") {
		return fmt.Errorf("coding.project is required")
	}

	if config.Migrate.Type == "repo" && (len(config.Source.Repo) == 0 || config.Source.Repo[0] == "") {
		return fmt.Errorf("coding.repo is required")
	}

	// 检查 cnb 参数
	if config.CNB.URL == "" {
		return fmt.Errorf("cnb.url is required")
	}

	// 检查 cnb.url 前缀
	if !strings.HasPrefix(config.CNB.URL, "http://") && !strings.HasPrefix(config.CNB.URL, "https://") {
		return fmt.Errorf("cnb.url must start with 'http://' or 'https://'")
	}

	if config.CNB.Token == "" {
		return fmt.Errorf("cnb.token is required")
	}

	if config.CNB.RootOrganization == "" {
		return fmt.Errorf("cnb.RootOrganization is required")
	}

	if config.Migrate.Concurrency < 1 {
		return fmt.Errorf("migrate.concurrency must be greater than 0")
	}

	if config.Migrate.organizationMappingLevel > 2 {
		return fmt.Errorf("organization_mapping_level error only support 1 or 2 ")
	}

	//logger.Logger.Infof("配置检查通过")
	return nil
}

func init() {
	Cfg = viper.New()
	Cfg.SetConfigName("config")
	Cfg.SetConfigType("yaml")
	Cfg.AddConfigPath(".")

	// 读取环境变量
	Cfg.AutomaticEnv()

	//envVariables := os.Environ()
	//println(envVariables)

	// 设置环境变量前缀
	Cfg.SetEnvPrefix("PLUGIN")
	Cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// 将环境变量绑定到配置项
	err := bindEnvVariables(Cfg)
	if err != nil {
		panic(err)
	}

	// 设置默认值
	setDefaultValues(Cfg)

	stringCovertToListAndSetConfigValue(Cfg, "source.project", "source.repo", "migrate.rebase_branch")

	err = parseStringEnvValueToBool(Cfg, "migrate.force_push", "migrate.ignore_lfs_notfound_error", "migrate.use_lfs_migrate", "migrate.allow_incomplete_push", "migrate.skip_exists_repo", "migrate.release", "migrate.code", "migrate.ssh", "migrate.rebase")
	if err != nil {
		panic(err)
	}
	err = parseStringEnvValueToInt(Cfg, "migrate.concurrency", "migrate.organization_mapping_level")
	if err != nil {
		panic(err)
	}

	// 保存配置到文件
	err = Cfg.WriteConfigAs(configName)
	if err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}

	// 读取配置文件
	err = Cfg.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

}

func ConvertToApiURL(baseUrl string) (apiUrl string) {
	parts := strings.Split(baseUrl, "://")
	apiUrl = parts[0] + "://" + "api." + parts[1]
	return apiUrl
}

func InitConfig() error {
	file, err := os.Create(defaultConfigName)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	data := []byte(defaultConfigContent)
	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("Failed to write into file: %v", err)
	}
	return nil
}

func parseStringEnvValueToBool(v *viper.Viper, keys ...string) error {
	for _, key := range keys {
		stringValue := v.GetString(key)
		boolValue, err := strconv.ParseBool(stringValue)
		if err != nil {
			return fmt.Errorf("failed to parse %s to bool: %v", key, err)
		}
		v.Set(key, boolValue)
	}
	return nil
}

func parseStringEnvValueToInt(v *viper.Viper, keys ...string) error {
	for _, key := range keys {
		stringValue := v.GetString(key)
		intValue, err := strconv.Atoi(stringValue)
		if err != nil {
			return fmt.Errorf("failed to parse %s to bool: %v", key, err)
		}
		v.Set(key, intValue)
	}
	return nil
}

func stringCovertToListAndSetConfigValue(v *viper.Viper, keys ...string) {
	for _, key := range keys {
		value := v.GetString(key)
		listValue := strings.Split(value, ",")
		if listValue[0] == "" {
			listValue = nil
		}
		v.Set(key, listValue)
	}
	return
}

func bindEnvVariables(config *viper.Viper) error {
	envKeys := []string{
		"source.url",
		"source.token",
		"source.platform",
		"source.project",
		"source.repo",
		"source.username",
		"source.password",
		"cnb.url",
		"cnb.token",
		"cnb.root_organization",
		"migrate.type",
		"migrate.concurrency",
		"migrate.force_push",
		"migrate.ignore_lfs_notfound_error",
		"migrate.use_lfs_migrate",
		"migrate.organization_mapping_level",
		"migrate.allow_incomplete_push",
		"migrate.log_level",
		"migrate.file_limit_size",
		"migrate.skip_exists_repo",
		"migrate.release",
		"migrate.code",
		"source.ak",
		"source.as",
		"source.endpoint",
		"migrate.ssh",
		"source.ssh_private_key",
		"migrate.rebase_branch",
	}
	for _, key := range envKeys {
		err := config.BindEnv(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func setDefaultValues(config *viper.Viper) {
	// 定义默认值
	defaults := map[string]string{
		"migrate.force_push":                 "false",
		"migrate.ignore_lfs_notfound_error":  "false",
		"migrate.use_lfs_migrate":            "false",
		"migrate.organization_mapping_level": "1",
		"migrate.concurrency":                "10",
		"migrate.type":                       "team",
		"migrate.allow_incomplete_push":      "false",
		"migrate.log_level":                  "info",
		"source.platform":                    "coding",
		"migrate.file_limit_size":            "500",
		"migrate.skip_exists_repo":           "true",
		"migrate.release":                    "false",
		"migrate.code":                       "true",
		"source.endpoint":                    "devops.cn-hangzhou.aliyuncs.com",
		"migrate.ssh":                        "false",
		"migrate.rebase":                     "false",
		"cnb.url":                            "https://cnb.cool",
		"source.url":                         "https://e.coding.net",
	}

	// 使用循环来设置默认值
	for key, value := range defaults {
		config.SetDefault(key, value)
	}
}
