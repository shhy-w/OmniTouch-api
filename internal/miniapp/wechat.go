package miniapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"git.uozi.org/uozi/cosy-example-api/model"
	appSettings "git.uozi.org/uozi/cosy-example-api/settings"
	cosySettings "github.com/uozi-tech/cosy/settings"
	"gorm.io/gorm"
)

type wxSessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func LoginByCode(db *gorm.DB, code string) (*model.MiniappUser, error) {
	session, err := fetchWxSession(code)
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	user := &model.MiniappUser{}
	err = db.Where("open_id = ?", session.OpenID).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &model.MiniappUser{
			OpenID:      session.OpenID,
			UnionID:     session.UnionID,
			SessionKey:  session.SessionKey,
			LastLoginAt: now,
		}
		return user, db.Create(user).Error
	}
	if err != nil {
		return nil, err
	}
	user.UnionID = session.UnionID
	user.SessionKey = session.SessionKey
	user.LastLoginAt = now
	return user, db.Save(user).Error
}

func fetchWxSession(code string) (*wxSessionResp, error) {
	conf := appSettings.WechatSettings
	if conf.DevMock && cosySettings.ServerSettings.RunMode == "debug" {
		if code == "" {
			code = "dev"
		}
		return &wxSessionResp{
			OpenID:     "dev-openid-" + code,
			SessionKey: "dev-session-key",
		}, nil
	}
	if conf.AppID == "" || conf.AppSecret == "" {
		return nil, errors.New("wechat appid or secret is empty")
	}
	endpoint := "https://api.weixin.qq.com/sns/jscode2session"
	query := url.Values{}
	query.Set("appid", conf.AppID)
	query.Set("secret", conf.AppSecret)
	query.Set("js_code", code)
	query.Set("grant_type", "authorization_code")
	resp, err := http.Get(endpoint + "?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data wxSessionResp
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if data.ErrCode != 0 {
		return nil, fmt.Errorf("wechat login failed: %s", data.ErrMsg)
	}
	if data.OpenID == "" {
		return nil, errors.New("wechat openid is empty")
	}
	return &data, nil
}
