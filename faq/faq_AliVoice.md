# FAQ_阿里云录音文件识别常见问题
> By EldersJavas

### 请根据自己的错误码对号入座
### 如果还是无法解决请到QQ群中询问，群号[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO)
- ###### 21050003_SUCCESS_WITH_NO_VALID_FRAGMENT
识别结果查询接口调用成功，但是没有识别到语音
检查录音文件是否有语音，或者语音时长太短。
- ###### 41050001_USER_BIZDURATION_QUOTA_EXCEED
单日时间超限
- ###### 41050002_FILE_DOWNLOAD_FAILED
文件下载失败
检查录音文件路径是否正确，是否可以外网访问和下载
- ###### 41050003_FILE_CHECK_FAILED
文件格式错误
检查录音文件是否是单轨/双轨的WAV格式、MP3格式
- ###### 41050004_FILE_TOO_LARGE
文件过大
检查录音文件大小是否超过512MB
- ###### 41050005_FILE_NORMALIZE_FAILED
文件归一化失败
检查录音文件是否有损坏，是否可以正常播放
- ###### 41050006_FILE_PARSE_FAILED
文件解析失败
检查录音文件是否有损坏，是否可以正常播放
- ###### 41050007_MKV_PARSE_FAILED
MKV解析失败
检查录音文件是否有损坏，是否可以正常播放
- ###### 41050008_UNSUPPORTED_SAMPLE_RATE
采样率不匹配
检查调用时设置的采样率和管控台上appkey绑定的ASR模型采样率是否一致
- ###### 41050009_UNSUPPORTED_ASR_GROUP
ASR分组不支持
确认下ak和appkey是否一致
###### 41050010_FILE_TRANS_TASK_EXPIRED
录音文件识别任务过期
TaskId不存在，或者已过期
- ###### 41050011_REQUEST_INVALID_FILE_URL_VALUE
请求file_link参数非法
请确认file_link参数格式是否正确
- ###### 41050012_REQUEST_INVALID_CALLBACK_VALUE
请求callback_url参数非法
请确认callback_url参数格式是否正确，是否为空
- ###### 41050013_REQUEST_PARAMETER_INVALID
请求参数无效
确认下请求task值为有效JSON格式字符串
- ###### 41050014_REQUEST_EMPTY_APPKEY_VALUE
请求参数appkey值为空
请确认是否设置了appkey参数值
- ###### 41050015_REQUEST_APPKEY_UNREGISTERED
请求参数appkey未注册
请确认请求参数appkey值是否设置正确，或者是否与阿里云账号的AccessKey ID同一个账号
- ###### 41050021_RAM_CHECK_FAILED
RAM检查失败
请检查您的RAM用户是否已经授权调用语音服务的API，请阅读开通服务—RAM用户鉴权配置
- ###### 41050023_CONTENT_LENGTH_CHECK_FAILED
content-length 检查失败
请检查下载文件时，http response中的content-length与文件实际大小是否一致.
- ###### 41050024_FILE_404_NOT_FOUND
需要下载的文件不存在
请检查需要下载的文件是否存在.
- ###### 41050025_FILE_403_FORBIDDEN
没有权限下载需要的文件
请检查是否有权限下载录音文件.
- ###### 41050026_FILE_SERVER_ERROR
请求的文件所在的服务不可用
请检查请求的文件所在的服务是否可用.
- ###### 51050000_INTERNAL_ERROR
内部通用错误
如果偶现可以忽略，重复出现请到QQ群中咨询，群号[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO)
- ###### 51050001_VAD_FAILED
VAD失败
如果偶现可以忽略，重复出现请到QQ群中咨询，群号[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO)
- ###### 51050002_RECOGNIZE_FAILED
内部alisr识别失败
如果偶现可以忽略，重复出现请到QQ群中咨询，群号[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO)
- ###### 51050003_RECOGNIZE_INTERRUPT
内部alisr识别中断
如果偶现可以忽略，重复出现请到QQ群中咨询，群号[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO)
- ###### 51050004_OFFER_INTERRUPT
内部写入队列中断
如果偶现可以忽略，重复出现请到QQ群中咨询，群号[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO)
- ###### 51050005_FILE_TRANS_TIMEOUT
内部整体超时失败
如果偶现可以忽略，重复出现请到QQ群中咨询，群号[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO)
- ###### 51050006_FRAGMENT_FAILED
内部分断失败
如果偶现可以忽略，重复出现请到QQ群中咨询，群号[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO)
