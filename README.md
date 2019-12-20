## VideoSrt 简介
 
`VideoSrt` 是用 `Golang`语言，基于 [lxn/walk](https://github.com/lxn/walk) Windows-GUI 工具包开发。

这是一个可以识别视频语音自动生成字幕SRT文件的开源软件工具。

适用于快速、批量的为视频创建中/英文字幕文件的业务场景。

本项目使用了阿里云的[OSS对象存储](https://www.aliyun.com/product/oss?spm=5176.12825654.eofdhaal5.13.e9392c4aGfj5vj&aly_as=K11FcpO8)、[录音文件识别](https://ai.aliyun.com/nls/filetrans?spm=5176.12061031.1228726.1.47fe3cb43I34mn) 以及 [百度翻译](http://api.fanyi.baidu.com/api/trans/product/index) 的相关业务接口。

CLI（命令行）版本：[https://github.com/wxbool/video-srt](https://github.com/wxbool/video-srt)


## 界面预览
![界面预览](https://ae01.alicdn.com/kf/H79c2202bb1734f15a870059b27fc9519E.gif)


## 优势
* 使用阿里云语音识别接口，准确度高，标准普通话/英语识别率95%以上
* 视频识别无需上传原视频，方便且节省时间
* 支持多任务多文件并发处理
* 支持视频、音频常见多种格式文件
* 支持输出字幕文件、普通文本两种类型
* 支付字幕中英互译、双语字幕输出


## Download 

##### 下载地址:(v0.2.1)
* .zip（含ffmpeg依赖） [点我下载](http://file.viggo.site/video-srt/0.2.1/video-srt-gui-ffmpeg-0.2.1-x64.zip)
* .zip（不含ffmpeg依赖） [点我下载](http://file.viggo.site/video-srt/0.2.1/video-srt-gui-0.2.1-x64.zip)

你也可以到 [release](https://github.com/wxbool/video-srt-windows/releases) 页面下载其他版本


## 注意事项
* 软件目录下的 `data`目录为数据存储目录，请勿删除。否则可能会导致配置丢失
* 项目使用了 [ffmpeg](http://ffmpeg.org/) 依赖，除非您的电脑已经安装了`ffmpeg`环境，否则请下载包含`ffmpeg`依赖的软件包

## 升级说明
* 先下载最新版本的软件包
* 然后用旧版本软件的 `data` 文件夹覆盖新版软件的 `data` 文件夹

## FAQ
* 软件支持哪些语言？
    * 视频字幕文本识别的核心服务是由阿里云`录音文件识别`业务提供的接口进行的，设置好对应的语音引擎，可以支持汉语普通话、方言、欧美英语等语言
* 如何开通阿里云的相关服务？
    * 注册阿里云账号
    * 账号快速实名认证
    * 开通 `访问控制` 服务，并创建角色，设置开放 `OSS对象存储`、`智能语音交互` 的访问权限 
    * 开通 `OSS对象存储` 服务，并创建一个存储空间（Bucket）（读写权限设置为公共读）
    * 开通 `智能语音交互` 服务，并创建项目（根据使用场景选择识别语言以及偏好等）
    * 设置 软件配置
    
## 联系方式
* QQ：2012210812
* 群：`暂无`
