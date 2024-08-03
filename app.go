package auth

import (
	"context"
	"fmt"
	"github.com/open4go/auth/model/app"
	r2mongo "github.com/open4go/db/mongo"
	"github.com/open4go/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetDBHandler 获取数据库handler 这里定义一个方法
func GetDBHandler(ctx context.Context) *mongo.Database {
	handler, err := r2mongo.DBPool.GetHandler("sys_auth")
	if err != nil {
		log.Log(ctx).Fatal(err)
	}
	return handler
}

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
	err := GetRedisAuthHandler(a.Ctx).HSet(a.Ctx, a.GlobalAppKey, path, name).Err()
	if err != nil {
		return err
	}
	return nil
}

// getNameByPath 通过请求路径快速找到应用名称
func (a *MyApp) getNameByPath(path, name string) string {
	name, err := GetRedisAuthHandler(a.Ctx).HGet(a.Ctx, a.GlobalAppKey, path).Result()
	if err != nil {
		return ""
	}
	return name
}

// GetAllPath 通过请求路径快速找到应用名称
func (a *MyApp) GetAllPath() map[string]string {
	name, err := GetRedisAuthHandler(a.Ctx).HGetAll(a.Ctx, a.GlobalAppKey).Result()
	if err != nil {
		return nil
	}
	return name
}

// setNameWithPath 通过请求路径快速设定应用属性
func (a *MyApp) setApiAttribute(path, name string, value interface{}) error {
	secondKey := fmt.Sprintf("%s:%s", path, name)
	err := GetRedisAuthHandler(a.Ctx).HSet(a.Ctx, a.GlobalAppAttr, secondKey, value).Err()
	if err != nil {
		return err
	}
	return nil
}

// getApiAttribute 通过请求路径快速获取应用属性
func (a *MyApp) getApiAttribute(path, name string) string {
	secondKey := fmt.Sprintf("%s:%s", path, name)
	val, err := GetRedisAuthHandler(a.Ctx).HGet(a.Ctx, a.GlobalAppAttr, secondKey).Result()
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
	//
	//for _, i := range app.AccessAPI {
	//	// 设置大的应用名称
	//	err := a.setNameWithPath(i.Path, app.Name)
	//	if err != nil {
	//		log.Log(a.Ctx).WithField("path", i.Path).WithField("name", app.Name).
	//			Error(err)
	//		continue
	//	}
	//
	//	err = a.setApiAttribute(i.Path, "name", i.Name)
	//	if err != nil {
	//		log.Log(a.Ctx).WithField("path", i.Path).WithField("name", i.Name).
	//			Error(err)
	//		continue
	//	}
	//
	//	err = a.setApiAttribute(i.Path, "desc", i.Desc)
	//	if err != nil {
	//		log.Log(a.Ctx).WithField("path", i.Path).WithField("desc", i.Desc).
	//			Error(err)
	//		continue
	//	}
	//
	//	err = a.setApiAttribute(i.Path, "disable", i.Disable)
	//	if err != nil {
	//		log.Log(a.Ctx).WithField("path", i.Path).WithField("disable", i.Disable).
	//			Error(err)
	//		continue
	//	}
	//
	//	err = a.setApiAttribute(i.Path, "can_view_detail", i.CanViewDetail)
	//	if err != nil {
	//		log.Log(a.Ctx).WithField("path", i.Path).WithField("can_view_detail", i.CanViewDetail).
	//			Error(err)
	//		continue
	//	}
	//
	//	err = a.setApiAttribute(i.Path, "hide_on_sidebar", i.HideOnSidebar)
	//	if err != nil {
	//		log.Log(a.Ctx).WithField("path", i.Path).WithField("hide_on_sidebar", i.HideOnSidebar).
	//			Error(err)
	//		continue
	//	}
	//
	//}
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
	handler := m.Init(context.TODO(), GetDBHandler(a.Ctx), m.CollectionName())
	_, err := handler.GetList(bson.D{}, &apps)
	if err != nil {
		panic(err)
	}
	a.LoadAllAppInfo(apps)
}
