package settings

type Wechat struct {
	AppID     string
	AppSecret string
	DevMock   bool
}

var WechatSettings = &Wechat{
	DevMock: true,
}
