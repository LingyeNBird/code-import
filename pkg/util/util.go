package util

import (
	"bytes"
	"ccrctl/pkg/logger"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	imgUrlRegexp = `!\[.*?\]\((.*?)\)`
)

func GetFunctionName() string {
	pc := make([]uintptr, 1) // at least 1 entry needed
	runtime.Callers(2, pc)   // 2 skips runtime.Callers and printFunctionName frames
	fn := runtime.FuncForPC(pc[0])
	return fn.Name()
}

// ConvertUrlWithAuth 把仓库httpURL转换为带认证的URL
func ConvertUrlWithAuth(url, username, password string) string {
	parts := strings.Split(url, "://")
	URL := parts[0] + "://" + username + ":" + password + "@" + parts[1]
	return URL
}

func RemoveHostFromURL(rawURL string) (string, error) {
	// 解析原始 URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// 创建一个新的 URL 对象，不包含 host
	newURL := &url.URL{
		Scheme:     parsedURL.Scheme,
		Path:       parsedURL.Path,
		RawPath:    parsedURL.RawPath,
		ForceQuery: parsedURL.ForceQuery,
		RawQuery:   parsedURL.RawQuery,
		Fragment:   parsedURL.Fragment,
	}

	// 返回新的 URL 字符串
	return newURL.String(), nil
}

// HasImgUrl check has img url
func HasImgUrl(content string) bool {
	re := regexp.MustCompile(imgUrlRegexp)
	return re.MatchString(content)
}

// MatchMarkdownImgUrl Match markdown iamge url
func MatchMarkdownImgeUrl(md string) map[string]string {
	re := regexp.MustCompile(imgUrlRegexp)
	matches := re.FindAllStringSubmatch(md, -1)
	imgUrls := make(map[string]string)

	// 打印所有匹配到的图片URL
	for _, match := range matches {
		if len(match) > 1 {
			index := strings.Index(match[1], ";")
			logger.Logger.Debugf("Matched img url: %s", match[1][:index])
			imgUrls[match[1][:index]] = match[1]
		}
	}
	return imgUrls
}

// ExtractAttachments  从 Markdown 内容中提取附件和图片的名称和 URL
func ExtractAttachments(markdown string) (attachments map[string]string, images map[string]string, exists bool) {
	// 正则表达式匹配 Markdown 链接（附件）
	linkRe := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	linkMatches := linkRe.FindAllStringSubmatch(markdown, -1)
	attachments = make(map[string]string)
	// 遍历链接匹配项并填充附件 map
	for _, match := range linkMatches {
		if len(match) == 3 {
			attachmentName := match[1] // 附件名
			attachmentURL := match[2]  // 附件 URL
			// 排除图片 URL
			if !isImageURL(attachmentURL) {
				attachments[attachmentName] = attachmentURL
			}
		}
	}

	// 正则表达式匹配 Markdown 图片，包括文件扩展名
	imageRe := regexp.MustCompile(`!\[(.*?)\]\((.*?)(\.[^.]+)?\)`)
	imageMatches := imageRe.FindAllStringSubmatch(markdown, -1)
	images = make(map[string]string)
	// 遍历图片匹配项并填充图片 map
	for _, match := range imageMatches {
		if len(match) == 4 {
			imageName := match[1] + match[3] // 图片名加上文件扩展名
			imageURL := match[2] + match[3]  // 图片 URL
			images[imageName] = imageURL
		}
	}

	// 检查是否找到任何附件或图片
	exists = len(attachments) > 0 || len(images) > 0

	return attachments, images, exists
}

// 检查 URL 是否以图片格式结尾
func isImageURL(url string) bool {
	imageExts := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}
	for _, ext := range imageExts {
		if strings.HasSuffix(strings.ToLower(url), ext) {
			return true
		}
	}
	return false
}

func isImageLink(markdownText, link string) bool {
	// 检查链接前是否有 '!'
	return regexp.MustCompile(`!\[` + regexp.QuoteMeta(link) + `\]`).MatchString(markdownText)
}

func RemoveImageFileExtension(filename string) string {
	// 定义常见的图片文件后缀
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg"}

	// 遍历后缀列表，检查并移除匹配的后缀
	for _, ext := range imageExtensions {
		if strings.HasSuffix(filename, ext) {
			return strings.TrimSuffix(filename, ext)
		}
	}
	// 如果没有找到匹配的后缀，返回原始文件名
	return filename
}

func GetFileNameFromURL(urlStr string) (string, error) {
	// 解析URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("parse URL failed: %w", err)
	}

	// 获取路径的最后一部分作为文件名
	fileName := path.Base(u.Path)

	// URL解码
	decodedFileName, err := url.QueryUnescape(fileName)
	if err != nil {
		return "", fmt.Errorf("decode filename failed: %w", err)
	}

	return decodedFileName, nil
}

func GetGoroutineID() int {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.Atoi(string(b))
	return n
}
