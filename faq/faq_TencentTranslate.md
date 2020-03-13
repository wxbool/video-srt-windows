
# FAQ_腾讯云翻译API常见问题
> By EldersJavas



## 常见错误

#### 公共错误码

- ###### AuthFailure.InvalidSecretId	
密钥非法（不是云 API 密钥类型）。
- ###### AuthFailure.MFAFailure	MFA 
错误。
- ###### AuthFailure.SecretIdNotFound	
密钥不存在。 请在控制台检查密钥是否已被删除或者禁用，如状态正常，请检查密钥是否填写正确，注意前后不得有空格。
- ###### AuthFailure.SignatureExpire	
签名过期。Timestamp 和服务器时间相差不得超过五分钟，请检查本地时间是否和标准时间同步。
- ###### AuthFailure.SignatureFailure	
签名错误。 签名计算错误，请对照调用方式中的接口鉴权文档检查签名计算过程。
- ###### AuthFailure.TokenFailure	token 
错误。
- ###### AuthFailure.UnauthorizedOperation	
请求未授权。请参考 CAM 文档对鉴权的说明。
- ###### DryRunOperation	DryRun 
操作，代表请求将会是成功的，只是多传了 DryRun 参数。
- ###### FailedOperation	
操作失败。
- ###### InternalError	
内部错误。
- ###### InvalidAction	
接口不存在。
- ###### InvalidParameter	
参数错误。
- ###### InvalidParameterValue	
参数取值错误。
- ###### LimitExceeded	
超过配额限制。
- ###### MissingParameter	
缺少参数错误。
- ###### NoSuchVersion	
接口版本不存在。
- ###### RequestLimitExceeded	
请求的次数超过了频率限制。
- ###### ResourceInUse	
资源被占用。
- ###### ResourceInsufficient	
资源不足。
- ###### ResourceNotFound	
资源不存在。
- ###### ResourceUnavailable	
资源不可用。
- ###### UnauthorizedOperation	
未授权操作。
- ###### UnknownParameter	
未知参数错误。
- ###### UnsupportedOperation	
操作不支持。
- ###### UnsupportedProtocol	
HTTP(S)请求协议错误，只支持 GET 和 POST 请求。
- ###### UnsupportedRegion	
接口不支持所传地域。
#### 业务错误码

- ###### FailedOperation.NoFreeAmount	
本月免费额度已用完，如需继续使用您可以在机器翻译控制台升级为付费使用。
- ###### FailedOperation.ServiceIsolate	
账号因为欠费停止服务，请在腾讯云账户充值。
- ###### FailedOperation.UserNotRegistered	
服务未开通，请在腾讯云官网机器翻译控制台开通服务。
- ###### InternalError	
内部错误。
- ###### InternalError.BackendTimeout	
后台服务超时，请稍后重试。
- ###### InternalError.ErrorUnknown	
未知错误。
- ###### InvalidParameter	
参数错误。
- ###### InvalidParameter.DuplicatedSessionIdAndSeq	
重复的SessionUuid和Seq组合。
- ###### InvalidParameter.SeqIntervalTooLarge	
Seq之间的间隙请不要大于2000。
- ###### LimitExceeded	
超过配额限制。
- ###### MissingParameter	
缺少参数错误。
- ###### UnauthorizedOperation.ActionNotFound	
请填写正确的Action字段名称。
- ###### UnsupportedOperation	
操作不支持。
- ###### UnsupportedOperation.TextTooLong	
单次请求text超过⻓长度限制，请保证单次请求⻓长度低于2000。
- ###### UnsupportedOperation.UnSupportedTargetLanguage	
不支持的目标语言，请参照语言列表。
- ###### UnsupportedOperation.UnsupportedLanguage	
不支持的语言，请参照语言列表。
- ###### UnsupportedOperation.UnsupportedSourceLanguage	
不支持的源语言，请参照语言列表。
### 如果还是无法解决请到QQ群中询问，群号[109695078](https://jq.qq.com/?_wv=1027&k=5Eco2hO)
