package auth

import (
	"context"
	"fmt"
	authModel "github.com/open4go/auth/model"
	"github.com/open4go/auth/model/app"
	"github.com/open4go/log"
	"go.mongodb.org/mongo-driver/bson"
)

type MyApp struct {
	Ctx           context.Context
	GlobalAppKey  string
	GlobalAppAttr string
}

func NewMyApp(ctx context.Context) MyApp {
	return MyApp{
		Ctx:           ctx,
		GlobalAppKey:  "global:app:name",
		GlobalAppAttr: "global:app:attr",
	}
}

// setNameWithPath 通过请求路径快速设定应用名称
func (a *MyApp) setNameWithPath(path, name string) error {
	err := RDB.HSet(a.Ctx, a.GlobalAppKey, path, name).Err()
	if err != nil {
		return err
	}
	return nil
}

// getNameByPath 通过请求路径快速找到应用名称
func (a *MyApp) getNameByPath(path, name string) string {
	name, err := RDB.HGet(a.Ctx, a.GlobalAppKey, path).Result()
	if err != nil {
		return ""
	}
	return name
}

// GetAllPath 通过请求路径快速找到应用名称
func (a *MyApp) GetAllPath() map[string]string {
	name, err := RDB.HGetAll(a.Ctx, a.GlobalAppKey).Result()
	if err != nil {
		return nil
	}
	return name
}

// setNameWithPath 通过请求路径快速设定应用属性
func (a *MyApp) setApiAttribute(path, name string, value interface{}) error {
	secondKey := fmt.Sprintf("%s:%s", path, name)
	err := RDB.HSet(a.Ctx, a.GlobalAppAttr, secondKey, value).Err()
	if err != nil {
		return err
	}
	return nil
}

// getApiAttribute 通过请求路径快速获取应用属性
func (a *MyApp) getApiAttribute(path, name string) string {
	secondKey := fmt.Sprintf("%s:%s", path, name)
	val, err := RDB.HGet(a.Ctx, a.GlobalAppAttr, secondKey).Result()
	if err != nil {
		return val
	}
	return val
}

// LoadAppInfo 加载应用信息
// "name": "货车",
// "desc": "货车信息",
// "disable": false,
// "can_view_detail": true,
// "hide_on_sidebar": false
func (a *MyApp) LoadAppInfo(app app.Model) {

	for _, i := range app.AccessAPI {
		// 设置大的应用名称
		err := a.setNameWithPath(i.Path, app.Name)
		if err != nil {
			log.Log().WithField("path", i.Path).WithField("name", app.Name).
				Error(err)
			continue
		}

		err = a.setApiAttribute(i.Path, "name", i.Name)
		if err != nil {
			log.Log().WithField("path", i.Path).WithField("name", i.Name).
				Error(err)
			continue
		}

		err = a.setApiAttribute(i.Path, "desc", i.Desc)
		if err != nil {
			log.Log().WithField("path", i.Path).WithField("desc", i.Desc).
				Error(err)
			continue
		}

		err = a.setApiAttribute(i.Path, "disable", i.Disable)
		if err != nil {
			log.Log().WithField("path", i.Path).WithField("disable", i.Disable).
				Error(err)
			continue
		}

		err = a.setApiAttribute(i.Path, "can_view_detail", i.CanViewDetail)
		if err != nil {
			log.Log().WithField("path", i.Path).WithField("can_view_detail", i.CanViewDetail).
				Error(err)
			continue
		}

		err = a.setApiAttribute(i.Path, "hide_on_sidebar", i.HideOnSidebar)
		if err != nil {
			log.Log().WithField("path", i.Path).WithField("hide_on_sidebar", i.HideOnSidebar).
				Error(err)
			continue
		}

	}
}

func (a *MyApp) LoadAllAppInfo(apps []app.Model) {
	for _, i := range apps {
		a.LoadAppInfo(i)
	}
}

func (a *MyApp) InitApp() {
	// 查询该角色下的应用
	// 获取当前用户角色的应用列表
	m := &app.Model{}
	apps := make([]app.Model, 0)
	handler := m.Init(context.TODO(), authModel.MDB, m.CollectionName())
	_, err := handler.GetList(bson.D{}, &apps)
	if err != nil {
		panic(err)
	}
	a.LoadAllAppInfo(apps)
}
