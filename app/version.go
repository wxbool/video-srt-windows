package app

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/lxn/walk"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	VERSION_SOURCE string = "https://gitee.com/641453620/video-srt-windows/tags"
)

type AppVersion struct {

}

//根据码云查询新版本
func (v *AppVersion) GetVersion () (string , error) {

	timeout := time.Duration(2 * time.Second)
	client := &http.Client{
		Timeout:timeout,
	}
	//获取版本来源html
	response, err := client.Get(VERSION_SOURCE)
	if err != nil {
		return "" , err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "" , errors.New("status code error : " + strconv.Itoa(response.StatusCode))
	}

	//加载html
	html, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "" , err
	}

	//查找节点
	vs := ""
	html.Find("#git-tags-container .releases-tags-wrap .releases-tag-content .tag-list .tag-item").Each(func(i int, s *goquery.Selection) {
		if vs != "" {
			return
		}
		tag , is := s.Find(".tag-name a").Attr("title")
		if is && tag != "" {
			vs = strings.TrimSpace(tag)
		}
	})

	return vs,nil
}


//显示更新提醒
func (v *AppVersion) ShowVersionNotifyInfo (version string , own *MyMainWindow) error {
	mw, err := walk.NewMainWindow()
	if err != nil {
		return err
	}
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		return err
	}

	defer func() {
		time.Sleep(time.Second * 15)
		_ = ni.Dispose()
	}()

	if err := ni.SetVisible(true); err != nil {
		return err
	}
	if err := ni.ShowMessage("更新提醒" , "检测到VideoSrt的新版本（v"+version+"），请及时下载更新哦") ; err != nil {
		return err
	}
	return nil
}