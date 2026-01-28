# data-sanitization 规范增量

本文档描述 `trim-repo-description-spaces` 变更对数据清理能力的规范增量。

## ADDED Requirements

### Requirement: 仓库描述空格自动修剪

系统 SHALL 在调用 CNB API 创建仓库之前,自动修剪仓库描述字段的前导和尾随空格。

**理由**: CNB 平台的 API 验证规则要求描述不能以空格开头或结尾。源代码托管平台(CODING、GitHub、GitLab 等)允许描述包含前后空格,因此需要在目标平台创建时进行清理以确保兼容性。

#### Scenario: 描述前面有空格

- **假设** 源仓库的描述为 `"  这是一个测试仓库"`(前面有2个空格)
- **当** 系统调用 `CreateRepo` 函数创建 CNB 仓库
- **那么** 系统必须将描述修剪为 `"这是一个测试仓库"`
- **并且** API 请求体中的 `description` 字段必须不包含前导空格
- **并且** 仓库创建必须成功
- **并且** CNB 平台上显示的描述必须为 `"这是一个测试仓库"`

#### Scenario: 描述后面有空格

- **假设** 源仓库的描述为 `"这是一个测试仓库  "`(后面有2个空格)
- **当** 系统调用 `CreateRepo` 函数创建 CNB 仓库
- **那么** 系统必须将描述修剪为 `"这是一个测试仓库"`
- **并且** API 请求体中的 `description` 字段必须不包含尾随空格
- **并且** 仓库创建必须成功
- **并且** CNB 平台上显示的描述必须为 `"这是一个测试仓库"`

#### Scenario: 描述前后都有空格

- **假设** 源仓库的描述为 `"  这是一个测试仓库  "`(前后各有2个空格)
- **当** 系统调用 `CreateRepo` 函数创建 CNB 仓库
- **那么** 系统必须将描述修剪为 `"这是一个测试仓库"`
- **并且** API 请求体中的 `description` 字段必须不包含前导或尾随空格
- **并且** 仓库创建必须成功

#### Scenario: 描述中间的空格保留

- **假设** 源仓库的描述为 `"  一个  测试  仓库  "`
- **当** 系统调用 `CreateRepo` 函数创建 CNB 仓库
- **那么** 系统必须将描述修剪为 `"一个  测试  仓库"`
- **并且** 描述中间的所有空格(包括多个连续空格)必须被保留
- **并且** 仅删除前导和尾随空格

#### Scenario: 仅包含空格的描述

- **假设** 源仓库的描述为 `"     "`(仅包含5个空格)
- **当** 系统调用 `CreateRepo` 函数创建 CNB 仓库
- **那么** 系统必须将描述修剪为空字符串 `""`
- **并且** API 请求体中的 `description` 字段必须为空字符串
- **并且** 仓库创建必须成功

#### Scenario: 空描述不受影响

- **假设** 源仓库的描述为空字符串 `""`
- **当** 系统调用 `CreateRepo` 函数创建 CNB 仓库
- **那么** 描述必须保持为空字符串 `""`
- **并且** 不得发生任何错误

#### Scenario: 描述包含换行符和制表符

- **假设** 源仓库的描述为 `"\n\t这是测试\n\t"`(包含换行符和制表符)
- **当** 系统调用 `CreateRepo` 函数创建 CNB 仓库
- **那么** 系统必须使用 `strings.TrimSpace` 删除所有前导和尾随的空白字符(空格、制表符、换行符)
- **并且** 修剪后的描述必须为 `"这是测试"`
- **并且** 仓库创建必须成功

#### Scenario: 超长描述的空格修剪和截断顺序

- **假设** 源仓库的描述为 `"  ` + (360个字符的文本) + `  "`
- **假设** 修剪后的描述长度为 360 字符,超过 `RepoDescLimitSize`(350)
- **当** 系统调用 `CreateRepo` 函数创建 CNB 仓库
- **那么** 系统必须先修剪前后空格
- **然后** 系统必须检查长度并截断到 350 字符
- **并且** 最终的描述长度必须等于 350
- **并且** 最终的描述必须不包含前导空格
- **注意**: 截断后可能在末尾产生新的空格(如果截断点在单词边界),此场景不需要二次修剪

### Requirement: 修剪操作的一致性

系统 SHALL 对所有源平台(CODING、GitHub、GitLab、Gitee、Gitea、阿里云、华为云、工蜂等)的仓库描述应用相同的空格修剪逻辑。

**理由**: 确保迁移行为的一致性和可预测性,无论源平台的数据质量如何。

#### Scenario: 不同源平台的描述修剪一致性

- **假设** 从 CODING 平台迁移的仓库描述为 `"  CODING 仓库  "`
- **假设** 从 GitHub 平台迁移的仓库描述为 `"  GitHub 仓库  "`
- **假设** 从 GitLab 平台迁移的仓库描述为 `"  GitLab 仓库  "`
- **当** 系统为这三个仓库调用 `CreateRepo` 函数
- **那么** 所有三个描述必须都被修剪掉前后空格
- **并且** 修剪逻辑必须相同,不得因源平台不同而有差异
- **并且** 所有仓库创建必须成功

### Requirement: 修剪操作的透明性

系统 SHALL 在不需要用户配置的情况下自动执行描述空格修剪。

**理由**: 空格修剪是数据清理的基本操作,应作为默认行为,无需增加配置复杂度。

#### Scenario: 无配置自动修剪

- **假设** 用户未在配置文件中设置任何与描述修剪相关的选项
- **当** 系统创建仓库时
- **那么** 描述的空格修剪必须自动执行
- **并且** 不得因缺少配置而跳过修剪操作

#### Scenario: 修剪操作不记录额外日志

- **假设** 仓库描述为 `"  测试  "`
- **当** 系统修剪描述并创建仓库
- **那么** 系统不得为正常的修剪操作生成 INFO、WARN 或 ERROR 级别的日志
- **注意**: 未来如需要追踪修剪操作,可以考虑添加 DEBUG 级别日志,但这不是当前变更的要求

## MODIFIED Requirements

无现有需求被修改。

## REMOVED Requirements

无现有需求被移除。

## 实现指导

### 修改位置

**文件**: `pkg/api/target/api.go`  
**函数**: `CreateRepo`  
**行号**: 约 328-349

### 建议实现

```go
func CreateRepo(url, token, group, repoName, repoDesc string, private bool) (err error) {
	var visibility string
	endpoint := group + "/-/repos"
	if private {
		visibility = "private"
	} else {
		visibility = "public"
	}
	
	// 修剪描述的前后空格
	repoDesc = strings.TrimSpace(repoDesc)
	
	// 检查长度限制
	if len(repoDesc) > RepoDescLimitSize {
		repoDesc = repoDesc[:RepoDescLimitSize]
	}
	
	body := &CreateRepoBody{
		Name:        repoName,
		Visibility:  visibility,
		Description: repoDesc,
	}
	_, err = c.Request("POST", endpoint, token, body)
	if err != nil {
		return err
	}
	return nil
}
```

### 测试建议

创建或更新 `pkg/api/target/api_test.go`,添加以下测试用例:

```go
func TestTrimRepoDescription(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"前导空格", "  测试描述", "测试描述"},
		{"尾随空格", "测试描述  ", "测试描述"},
		{"前后空格", "  测试描述  ", "测试描述"},
		{"中间空格保留", "测试 描述", "测试 描述"},
		{"空描述", "", ""},
		{"仅空格", "   ", ""},
		{"换行符和制表符", "\n\t测试\n\t", "测试"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimSpace(tt.input)
			if result != tt.expected {
				t.Errorf("期望 %q, 得到 %q", tt.expected, result)
			}
		})
	}
}
```

## 相关规范

- **error-handling**: 本规范遵循优雅降级原则,在数据不完美时自动修复而非失败
- **migration-logging**: 修剪操作不生成额外日志,保持日志简洁
