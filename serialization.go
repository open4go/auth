package auth

import (
	"encoding/base64"
	"encoding/json"
)

// DumpLoginInfo 登陆信息
func DumpLoginInfo(namespace string, userId string, avatar string,
	loginType string, userName string, accountId string, loginLevel string) (string, error) {
	// step 01 转换为json
	loginInfo := LoginInfo{
		Namespace:  namespace,
		AccountId:  accountId,
		UserId:     userId,
		UserName:   userName,
		Avatar:     avatar,
		LoginType:  loginType,
		LoginLevel: loginLevel,
	}
	payload, err := json.Marshal(loginInfo)
	if err != nil {
		return "", err
	}
	sEnc := base64.StdEncoding.EncodeToString([]byte(payload))
	return sEnc, nil
}

// LoadLoginInfo 解析登陆信息
func LoadLoginInfo(payload string) (*LoginInfo, error) {
	// step 01 转换为bytes
	sDec, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}
	loginInfo := &LoginInfo{}
	err = json.Unmarshal(sDec, loginInfo)
	if err != nil {
		return nil, err
	}
	return loginInfo, nil
}
