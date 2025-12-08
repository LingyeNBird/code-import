package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"ccrctl/pkg/util"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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
	codingTokenLength = 40
)

var Cfg *viper.Viper

type Config struct {
	Source  Source  `yaml:"source"`
	CNB     CNB     `yaml:"cnb"`
	Migrate Migrate `yaml:"migrate"`
}

type Source struct {
	URL            string   `yaml:"url"`
	Token          string   `yaml:"token"`
	Project        []string `yaml:"project"`
	Repo           []string `yaml:"repo"`
	Platform       string   `yaml:"platform"`
	UserName       string   `yaml:"username"`
	Password       string   `yaml:"password"`
	SshPrivateKey  string   `yaml:"ssh_private_key"`
	Group          string   `yaml:"group"`
	OrganizationId string   `yaml:"organizationid"`
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
	allowIncompletePush  bool   `yaml:"allow_incomplete_push"`
	LogLevel             string `yaml:"log_level"`
	fileLimitSize        int64  `yaml:"file_limit_size"`
	SkipExistsRepo       bool   `yaml:"skip_exists_repo"`
	Ssh                  bool   `yaml:"ssh"`
	AllowSelectRepos     bool   `yaml:"allow_select_repos"`
	DownloadOnly         bool   `yaml:"download_only"`
	MapCodingDisplayName bool   `yaml:"map_coding_display_name"`
	MapCodingDescription bool   `yaml:"map_coding_description"`
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
	// 检查是否为只下载模式
	downloadOnly := config.Migrate.DownloadOnly

	platform := config.Source.Platform
	if platform != "aliyun" && platform != "local" {
		err := checkURL(config.Source.URL)
		if err != nil {
			return err
		}
	}

	//非通用第三方平台迁移，检查 source.token 参数（local 不需要）
	if platform != "common" && platform != "local" {
		if err := checkTokenValid(config.Source.Token, platform); err != nil {
			return err
		}
	}

	// aliyun 平台需要检查 source.organizationid
	if platform == "aliyun" {
		if config.Source.OrganizationId == "" {
			return fmt.Errorf("when platform is aliyun, source.organizationid is required")
		}
	}

	//common http迁移（local 不需要）
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
	// migrate.type 不再是必填项，如果未配置则默认为 team
	if config.Migrate.Type == "" {
		config.Migrate.Type = "team"
	}

	if config.Migrate.Type != "repo" && config.Migrate.Type != "project" && config.Migrate.Type != "team" {
		return fmt.Errorf("migrate.type error only support repo or project or team")
	}

	// CODING 平台的特殊校验已移除，改为由 source.repo 和 source.project 自动判断
	// 其他平台保持原有逻辑

	// repo 模式下，只有 common 和 local 平台需要强制要求 source.repo
	// 其他平台（gitlab、github、gitee、gongfeng等）可以通过 source.repo 过滤仓库
	if config.Migrate.Type == "repo" && (platform == "common" || platform == "local") && (len(config.Source.Repo) == 0 || config.Source.Repo[0] == "") {
		return fmt.Errorf("when migrate.type is repo and platform is common or local, source.repo is required")
	}

	// 如果不是只下载模式，则检查 CNB 相关配置
	if !downloadOnly {
		err := checkURL(config.CNB.URL)
		if err != nil {
			return err
		}

		if config.CNB.Token == "" {
			return fmt.Errorf("cnb.token is required")
		}

		if config.CNB.RootOrganization == "" {
			return fmt.Errorf("cnb.RootOrganization is required")
		}

		if config.Migrate.organizationMappingLevel > 2 {
			return fmt.Errorf("organization_mapping_level error only support 1 or 2 ")
		}
	}

	if strings.HasPrefix(config.CNB.RootOrganization, "/") {
		return fmt.Errorf("cnb.RootOrganization 不能以 / 开头")
	}

	if config.Migrate.Concurrency < 1 {
		return fmt.Errorf("migrate.concurrency must be greater than 0")
	}

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

	stringCovertToListAndSetConfigValue(Cfg, "source.project", "source.repo")

	// 需要转换为布尔值的配置项
	boolKeys := []string{
		"migrate.force_push",
		"migrate.ignore_lfs_notfound_error",
		"migrate.use_lfs_migrate",
		"migrate.allow_incomplete_push",
		"migrate.skip_exists_repo",
		"migrate.release",
		"migrate.code",
		"migrate.ssh",
		"migrate.rebase",
		"migrate.allow_select_repos",
		"migrate.download_only",
		"migrate.map_coding_display_name",
		"migrate.map_coding_description",
	}

	err = parseStringEnvValueToBool(Cfg, boolKeys...)
	if err != nil {
		panic(err)
	}
	err = parseStringEnvValueToInt(Cfg, "migrate.concurrency", "migrate.organization_mapping_level")
	if err != nil {
		panic(err)
	}

	// 处理 cnb.root_organization，使用 TrimSlash 函数清理前后的斜杠和空格
	processCNBRootOrganization(Cfg)

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
	if len(parts) == 2 {
		apiUrl = parts[0] + "://" + "api." + parts[1]
	} else {
		apiUrl = baseUrl
	}
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

// processCNBRootOrganization 处理 cnb.root_organization 配置项
// 使用 util.TrimSlash 函数移除前后的斜杠和空格
func processCNBRootOrganization(v *viper.Viper) {
	key := "cnb.root_organization"
	value := v.GetString(key)
	if value != "" {
		// 使用 util.TrimSlash 函数处理值
		trimmedValue := util.TrimSlash(value)
		v.Set(key, trimmedValue)
	}
}

func stringCovertToListAndSetConfigValue(v *viper.Viper, keys ...string) {
	for _, key := range keys {
		value := v.GetString(key)
		// 去除前后空格和末尾的逗号
		value = strings.TrimSpace(value)
		value = strings.TrimSuffix(value, ",")
		
		listValue := strings.Split(value, ",")
		if listValue[0] == "" {
			listValue = nil
		}
		v.Set(key, listValue)
	}
}

func bindEnvVariables(config *viper.Viper) error {
	envKeys := []string{
		"source.group",
		"source.url",
		"source.token",
		"source.platform",
		"source.project",
		"source.repo",
		"source.username",
		"source.password",
		"source.region",
		"source.organizationId",
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
		"migrate.allow_select_repos",
		"migrate.download_only",
		"migrate.include_github_fork",
		"migrate.map_coding_display_name",
		"migrate.map_coding_description",
		"migrate.gitlab_projects_owned",
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
		"migrate.ignore_lfs_notfound_error":  "true",
		"migrate.use_lfs_migrate":            "true",
		"migrate.organization_mapping_level": "1",
		"migrate.concurrency":                "5",
		"migrate.type":                       "team",
		"migrate.allow_incomplete_push":      "false",
		"migrate.log_level":                  "info",
		"source.platform":                    "coding",
		"migrate.file_limit_size":            "256",
		"migrate.skip_exists_repo":           "false",
		"migrate.release":                    "false",
		"migrate.code":                       "true",
		"migrate.ssh":                        "false",
		"migrate.rebase":                     "false",
		"cnb.url":                            "https://cnb.cool",
		"source.url":                         "https://e.coding.net",
		"migrate.allow_select_repos":         "false",
		"migrate.download_only":              "false",
		"migrate.include_github_fork":        "true",
		"migrate.map_coding_display_name":    "true",
		"migrate.map_coding_description":     "true",
		"source.region":                      "cn-north-4",
		"migrate.gitlab_projects_owned":      "false",
	}

	// 使用循环来设置默认值
	for key, value := range defaults {
		config.SetDefault(key, value)
	}
}

// token 合规性检查函数
func checkTokenValid(token string, platform string) error {
	if err := checkCommonToken(token); err != nil {
		return err
	}
	// coding 平台 token 检查
	if platform == "coding" {
		if err := checkCodingToken(token); err != nil {
			return err
		}
	}
	return nil
}

func checkCodingToken(token string) error {
	codingPattern := regexp.MustCompile(`^[a-z0-9]{40}$`)
	if !codingPattern.MatchString(token) {
		return fmt.Errorf("source.token 不符合CODING平台token规范，只能包含小写字母、数字，长度 40 个字符，请重新配置。正则匹配规则:%s", codingPattern)
	}
	return nil
}

func checkCommonToken(token string) error {
	commonPattern := regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)
	if !commonPattern.MatchString(token) {
		return fmt.Errorf("source.token 包含非法字符，只能包含字母、数字、中划线、下划线、点号。正则匹配规则:%s", commonPattern)
	}
	return nil
}

func checkURL(url string) error {
	// 允许 http:// 或 https:// 开头，后跟域名或IP地址，支持端口号，支持路径
	// 支持域名格式：example.com, api.example.com
	// 支持IP地址格式：192.168.1.1, 10.0.0.1
	// 支持路径：/path/to/resource

	// IP地址模式：严格匹配IPv4地址，支持端口号和路径
	ipPattern := `^https?://((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(:[0-9]+)?(/.*)?$`

	// 域名模式：至少包含一个点的域名，且至少有一个部分包含字母，支持端口号和路径
	domainPattern := `^https?://(([A-Za-z0-9-]*[A-Za-z][A-Za-z0-9-]*\.)+[A-Za-z0-9-]*[A-Za-z][A-Za-z0-9-]*|([A-Za-z0-9-]*[A-Za-z][A-Za-z0-9-]*\.)+[A-Za-z0-9-]+|([A-Za-z0-9-]+\.)+[A-Za-z0-9-]*[A-Za-z][A-Za-z0-9-]*)(:[0-9]+)?(/.*)?$`

	ipRegex := regexp.MustCompile(ipPattern)
	domainRegex := regexp.MustCompile(domainPattern)

	// 先检查IP地址模式，如果匹配则直接返回成功
	if ipRegex.MatchString(url) {
		return nil
	}

	// 再检查域名模式
	if domainRegex.MatchString(url) {
		return nil
	}

	return fmt.Errorf("url %s 格式错误，必须以 'http://' 或 'https://' 开头，支持域名或IP地址，支持端口号和路径，如 https://e.coding.net、https://cnb.cool/api、https://example.com:8080/path、http://192.168.1.1:8080/api", url)
}
