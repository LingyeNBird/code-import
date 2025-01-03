package http_client

import (
	"ccrctl/pkg/logger"
	"fmt"
	"io/ioutil"
	"net/http"
)

func DownloadFromUrl(fileUrl string) (data []byte, err error) {
	logger.Logger.Debugf("Get file url: %s", fileUrl)

	// 创建请求并添加认证信息
	req, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to download file: %s, status code: %d\n", fileUrl, resp.StatusCode)
		return nil, fmt.Errorf("failed to download file: %s, status code: %d", fileUrl, resp.StatusCode)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Errorf("Read file error: %v", err)
		return nil, err
	}

	return data, nil
}
