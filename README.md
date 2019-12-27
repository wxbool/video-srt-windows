# VideoSrt简介

`VideoSrt` 是用 `Golang`语言，基于 [lxn/walk](https://github.com/lxn/walk) Windows-GUI 工具包开发。

这是一个可以识别视频语音自动生成字幕SRT文件的开源软件工具。<br />适用于快速、批量的为视频创建中/英文字幕文件的业务场景。

本项目使用了阿里云的[OSS对象存储](https://www.aliyun.com/product/oss?spm=5176.12825654.eofdhaal5.13.e9392c4aGfj5vj&aly_as=K11FcpO8)、[录音文件识别](https://ai.aliyun.com/nls/filetrans?spm=5176.12061031.1228726.1.47fe3cb43I34mn) 以及 [百度翻译](http://api.fanyi.baidu.com/api/trans/product/index) 的相关业务接口。

CLI（命令行）版本：[https://github.com/wxbool/video-srt](https://github.com/wxbool/video-srt)

帮助文档/使用教程：[https://www.yuque.com/viggo-t7cdi/videosrt](https://www.yuque.com/viggo-t7cdi/videosrt)

<a name="0b884e4f"></a>
## 界面预览

![](https://cdn.nlark.com/yuque/0/2019/gif/695280/1577093439760-4c4ff7f0-25f9-4bd9-b7fd-b89655564a89.gif#align=left&display=inline&height=440&originHeight=440&originWidth=784&size=0&status=done&style=none&width=784)

## 应用场景

- 识别**视频/音频**的语音生成字幕文件（支持中英互译，双语字幕）
- 提取**视频/音频**的语音文本


<a name="b89d37d3"></a>
## 软件优势

- 使用阿里云语音识别接口，准确度高，标准普通话/英语识别率95%以上
- 视频识别无需上传原视频，方便且节省时间
- 支持多任务多文件批量处理
- 支持视频、音频常见多种格式文件
- 支持输出字幕文件、普通文本两种类型
- 支持字幕中英互译、双语字幕输出

<a name="Download"></a>
## Download

<a name="e66a66f1"></a>
##### 下载地址:(v0.2.5)

- .zip（含ffmpeg依赖） [点我下载](http://file.viggo.site/video-srt/0.2.5/video-srt-gui-ffmpeg-0.2.5-x64.zip)
- .zip（不含ffmpeg依赖） [点我下载](http://file.viggo.site/video-srt/0.2.5/video-srt-gui-0.2.5-x64.zip)

你也可以到 [release](https://github.com/wxbool/video-srt-windows/releases) 页面下载其他版本

<a name="1bbbb204"></a>
## 注意事项

- 软件目录下的 `data`目录为数据存储目录，请勿删除。否则可能会导致配置丢失
- 项目使用了 [ffmpeg](http://ffmpeg.org/) 依赖，除非您的电脑已经安装了`ffmpeg`环境，否则请下载包含`ffmpeg`依赖的软件包

<a name="9a751511"></a>
## 升级说明

- 先下载最新版本的软件包
- 然后用旧版本软件的 `data` 文件夹覆盖新版软件的 `data` 文件夹

<a name="f3dc992e"></a>
## 交流&联系

- QQ：2012210812
- QQ交流群：[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO)

    ![image.png](https://cdn.nlark.com/yuque/0/2019/png/695280/1577104071489-4cc85009-29a0-42d6-8901-0cf3b45dee68.png#align=left&display=inline&height=177&name=image.png&originHeight=177&originWidth=172&size=17846&status=done&style=none&width=172)

<a name="AyJ3E"></a>
## 捐赠&支持

![](https://pic2.superbed.cn/item/5e00b93476085c3289dd2dc0.png)