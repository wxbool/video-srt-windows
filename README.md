# VideoSrt

[![GitHub license](https://img.shields.io/github/license/wxbool/video-srt-windows?logo=Apache&style=flat-square)](https://github.com/wxbool/video-srt-windows/blob/master/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/wxbool/video-srt-windows?style=flat-square)](https://github.com/wxbool/video-srt-windows/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/wxbool/video-srt-windows?logo=github&style=flat-square)](https://github.com/wxbool/video-srt-windows/network)
![GitHub Action](https://github.com/wxbool/video-srt-windows/workflows/Go/badge.svg?branch=master)

## 简介

`VideoSrt` 是用 `Golang`语言，基于 [lxn/walk](https://github.com/lxn/walk) Windows-GUI 工具包开发。

这是一个可以识别视频语音自动生成字幕SRT文件的开源软件工具。<br />适用于快速、批量的为媒体（视频/音频）生成中/英文字幕、文本文件的业务场景。

0.3.5 版本将会使用以下接口：
- 阿里云 [OSS对象存储](https://www.aliyun.com/product/oss?spm=5176.12825654.eofdhaal5.13.e9392c4aGfj5vj&aly_as=K11FcpO8)
- 阿里云 [录音文件识别](https://ai.aliyun.com/nls/filetrans?spm=5176.12061031.1228726.1.47fe3cb43I34mn) 
- 百度翻译开放平台 [翻译API](http://api.fanyi.baidu.com/api/trans/product/index) 
- 腾讯云 [翻译API](https://cloud.tencent.com/product/tmt) 

软件帮助文档/使用教程看这个：[https://www.yuque.com/viggo-t7cdi/videosrt](https://www.yuque.com/viggo-t7cdi/videosrt)

B站Up主自制教程：[https://search.bilibili.com/all?keyword=videosrt](https://search.bilibili.com/all?keyword=videosrt)

##### Videosrt Pro 已发布，全新字幕生成神器，建议使用！产品说明详见：https://www.yuque.com/viggo-t7cdi/videosrt-pro/cpwtwl

线上“字幕生成/字幕翻译”解决方案：[字幕酱（付费）](https://www.zimujiang.com/aff?code=aannv4os)

<a name="0b884e4f"></a>
## 界面预览

![](https://pic.downk.cc/item/5f28c7ea14195aa594369a59.gif)

## 应用场景

- 识别**视频/音频**的语音生成字幕文件（支持中英互译，双语字幕）
- 提取**视频/音频**的语音文本
- 批量翻译、过滤处理/编码SRT字幕文件


<a name="b89d37d3"></a>
## 软件优势

- 使用阿里云语音识别接口，准确度高，标准普通话/英语识别率95%以上
- 视频识别无需上传原视频，方便快速且节省时间
- 支持多任务多文件批量处理
- 支持视频、音频常见多种格式文件
- 支持同时输出字幕SRT文件、LRC文件、普通文本3种类型
- 支持语气词过滤、自定义文本过滤、正则过滤等，使软件生成的字幕更加精准
- 支持字幕中英互译、双语字幕输出，及日语、韩语、法语、德语、西班牙语、俄语、意大利语、泰语等
- 支持多翻译引擎（百度翻译、腾讯云翻译）
- 支持批量翻译、编码SRT字幕文件

<a name="Download"></a>
## Download

<a name="e66a66f1"></a>
##### 下载地址:
- (v0.3.5)（含ffmpeg依赖） [点我下载](http://file.viggo.site/video-srt/0.3.5/video-srt-gui-ffmpeg-0.3.5-x64.zip)
- (v0.3.5)（不含ffmpeg依赖） [点我下载](http://file.viggo.site/video-srt/0.3.5/video-srt-gui-0.3.5-x64.zip)
- (v0.2.6)（含ffmpeg依赖） [点我下载](http://file.viggo.site/video-srt/0.2.6/video-srt-gui-ffmpeg-0.2.6-x64.zip)
- (v0.2.6)（不含ffmpeg依赖） [点我下载](http://file.viggo.site/video-srt/0.2.6/video-srt-gui-0.2.6-x64.zip)

你也可以到 [release](https://github.com/wxbool/video-srt-windows/releases) 页面下载其他版本

<a name="1bbbb204"></a>
## 注意事项

- 软件目录下的 `data`目录为数据存储目录，请勿删除。否则可能会导致配置丢失
- 项目使用了 [ffmpeg](http://ffmpeg.org/) 依赖，除非您的电脑已经安装了`ffmpeg`环境，否则请下载包含`ffmpeg`依赖的软件包

<a name="9a751511"></a>
## 升级说明

- 先下载最新版本的软件包
- 然后用旧版本软件的 `data` 文件夹覆盖新版软件的 `data` 文件夹
- 0.2.6 升级至 0.2.9 以上的版本时，由于翻译设置无法直接兼容低版本，可能需要重新在软件创建翻译引擎才能继续使用翻译功能

## 其他
今日，有小伙伴在淘宝、闲鱼上发现一些界面相似的盗版且额外收费的软件。由于`VideoSrt` 是开源软件，所以不排除个别不良用户通过微调软件源代码，然后进行二次售卖甚至传播病毒。 

目前已知的有：`大象字幕`

因此，软件的安装包不要在第三方平台渠道下载，尽量在 github/码云/帮助文档 等官方渠道的页面进行下载。

小伙伴们要擦亮眼睛，打击盗版，望周知！


## FAQ
##### 1.为什么Linux和Mac不能用？
因为`VideoSrt`的GUI是使用[lxn/walk](https://github.com/lxn/walk)开发的，仅支持Windows的GUI，如果您想在Linux上使用，可以体验[CLI版本](https://github.com/wxbool/video-srt)

##### 2.使用此软件会产生费用吗？
如果您适量使用本软件（各个API的免费使用额度可以自行查询），将不会产生费用。
如果您大量使用，建议根据自己的情况购买各个平台的资源包，以满足需求。

##### 3.难受，为什么我一直报错？
报错的原因有很多，软件配置错误、阿里云、腾讯云等账户权限问题都可能会导致软件显示错误。如果您遇到麻烦，建议加入`QQ群`与我们交流。


<a name="f3dc992e"></a>
## 交流&联系

- QQ交流群3：[742310571](https://jq.qq.com/?_wv=1027&k=OPKLP16H) 
- QQ交流群2：[1035527812](https://jq.qq.com/?_wv=1027&k=lzbkMgzJ) (已满)
- QQ交流群1：[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO) (已满)


<a name="AyJ3E"></a>
## 捐赠&支持

![](https://cdn.nlark.com/yuque/0/2021/png/695280/1633671491324-b3375ad6-027b-4297-bd7f-2f42784dccbf.png?x-oss-process=image%2Fwatermark%2Ctype_d3F5LW1pY3JvaGVp%2Csize_23%2Ctext_VmlkZW9TcnQ%3D%2Ccolor_FFFFFF%2Cshadow_50%2Ct_80%2Cg_se%2Cx_10%2Cy_10%2Fresize%2Cw_794%2Climit_0)