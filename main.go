package main

import (
	"GoBingWallpaper/bing"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type MyMainWindow struct {
	*walk.MainWindow
}

var infoLabel *walk.TextEdit

const logo = "res/logo.ico"

//go build -ldflags "-H windowsgui -w -s"
func main() {
	mw := &MyMainWindow{}
	err := MainWindow{
		MenuItems: []MenuItem{
			Action{
				Text: "关于",
				OnTriggered: func() {
					walk.MsgBox(mw, "关于", "Bing壁纸下载器v1.0\n作者：Mainli", walk.MsgBoxIconQuestion)
				},
			},
		},
		AssignTo: &mw.MainWindow, //窗口重定向至mw，重定向后可由重定向变量控制控件
		Title:    "必应壁纸下载器",      //标题
		MinSize:  Size{Width: 150, Height: 200},
		Size:     Size{300, 400},
		Layout:   VBox{}, //样式，纵向
		Children: []Widget{ //控件组
			PushButton{
				Text: "打开下载目录",
				OnClicked: func() { //点击事件响应函数
					file, err := os.Open(bing.DEFAULT_PATH)
					defer file.Close()
					if err == nil {
						_, ReaddirErr := file.Readdir(1)
						if ReaddirErr == nil {
							abs, _ := filepath.Abs(bing.DEFAULT_PATH)
							cmd := exec.Command("explorer", abs)
							cmd.Start()
							return
						}
					}
					walk.MsgBox(mw, "提示", "还没有下载任何图片", walk.MsgBoxIconQuestion)

				},
			},
			PushButton{
				Text: "开始下载",
				OnClicked: func() { //点击事件响应函数
					go clickDownload()
				},
			},
			TextEdit{
				VScroll:  true,
				ReadOnly: true,
				AssignTo: &infoLabel,
			},
		},
	}.Create() //创建
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if file, error := walk.NewIconFromFile(logo); error == nil {
		mw.SetIcon(file)
	}
	mw.Fullscreen()
	//获取屏幕宽高居中显示
	mw.SetX((int(win.GetSystemMetrics(0)) - mw.Width()) / 2)
	mw.SetY((int(win.GetSystemMetrics(1)) - mw.Height()) / 2)
	mw.Run() //运行
}
func clickDownload() {
	//os.RemoveAll(bing.DEFAULT_PATH)
	if get, err := http.Get("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=1&n=10"); err == nil {
		body := get.Body
		defer body.Close()
		data := new(bing.MetaData)
		if err := json.NewDecoder(body).Decode(data); err != nil {
			infoLabel.AppendText("error: ")
			infoLabel.AppendText(err.Error())
		} else {
			bing.FixData(data)
			bing.EnsureDir(bing.DEFAULT_PATH)
			infoLabel.AppendText("--------start-------------\r\n")
			waitGroup := bing.DownloadAllData(bing.DEFAULT_PATH, data, infoLabel)
			infoLabel.AppendText("--------end-------------\r\n")
			waitGroup.Wait()
		}
	}
}
