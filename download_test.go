package main

import (
	. "GoBingWallpaper/bing"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestDownlad(t *testing.T) {
	t.Log("start ---------------")
	if get, err := http.Get("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=1&n=10"); err == nil {
		body := get.Body
		defer body.Close()
		data := new(MetaData)
		if err := json.NewDecoder(body).Decode(data); err != nil {
			fmt.Println("错误", err)
		} else {
			FixData(data)
			EnsureDir(DEFAULT_PATH)
			wg := new(sync.WaitGroup)
			for _, item := range data.Images {
				targeetName := filepath.Join(DEFAULT_PATH, item.Name+".jpg")
				if _, err := os.Stat(targeetName); os.IsNotExist(err) {
					wg.Add(1)
					go DownloadUrlFile(targeetName, item.Url, wg)
				}
			}
			wg.Wait()
			fmt.Println("下载完毕")
		}
	}
	t.Log("end ---------------")
}
