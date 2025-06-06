package test

import (
	"github.com/luolayo/geocn-go/updater"
	"net/http"
	"testing"
)

// 测试 3 个下载地址是否可用（HEAD 取 200）
func TestDownloadURLsReachable(t *testing.T) {
	for _, item := range updater.DownloadList() { // 新增导出函数
		resp, err := http.Head(item.URL)
		if err != nil {
			t.Fatalf("HEAD %s error: %v", item.URL, err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("HEAD %s got status %d", item.URL, resp.StatusCode)
		}
	}
}
