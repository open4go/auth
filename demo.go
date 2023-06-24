package auth

import (
	"context"
	log "github.com/sirupsen/logrus"
)

const (
	demoAccountID = "1234"
)

func Demo() {
	sa := &SimpleAuth{}
	// 调用其他的任何服务前都需要调用BindKey 来完成数据的keys绑定
	sa = sa.BindKey(demoAccountID)
	// 查询用户账号
	// 通过手机号查询数据库信息
	err := sa.SignIn(context.TODO(), demoAccountID)
	if err != nil {
		log.Fatal(err)
	}

	// 读取角色列表

	//	读取应用列表

}
