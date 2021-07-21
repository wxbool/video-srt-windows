package ffmpeg

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"videosrt/app/tool"
)

type Ffmpeg struct {
	Os string //ffmpeg 文件目录
}


//提取视频音频
func ExtractAudio (video string , tmpAudio string) (error) {
	//校验
	if e := VailFfmpegLibrary();e != nil {
		return e
	}
	cmd := exec.Command("ffmpeg" , "-i" , video , "-ar" , "16000" , tmpAudio)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}


//校验依赖库
func VailFfmpegLibrary() error {
	ts := exec.Command("ffmpeg" , "-version")
	ts.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := ts.Run() ; err != nil {
		return errors.New("请先安装 ffmpeg 依赖 ，并设置环境变量")
	}
	return nil
}


//校验ffmpeg并加入临时环境遍历
func VailTempFfmpegLibrary(rootDir string)  {
	ffmpegDir := tool.WinDir(rootDir + "/ffmpeg")

	if tool.DirExists(ffmpegDir) {
		//临时加入用户环境变量
		path := os.Getenv("PATH")
		ps := strings.Split(path , ";")
		ps = append(ps , ffmpegDir)
		path = strings.Join(ps , ";")

		_ = os.Setenv("PATH", path)
	}
}