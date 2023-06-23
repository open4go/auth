package auth

import (
	"context"

	"github.com/r2day/collections"
)

func setMenu(ctx context.Context, accessAPIList []collections.APIInfo, hideSidebarKey string,
	keyPath2Name string, path2roles map[string][]string, roleID string) (err error) {

	for _, apiInfo := range accessAPIList {
		// 默认是false
		// 如果是true则忽略本条规则
		if apiInfo.Disable {
			continue
		}

		err = RDB.HSet(ctx, keyPath2Name, apiInfo.Path, apiInfo.Name).Err()
		if err != nil {
			continue
		}

		// 判断是否需要在导航menu中展示
		// 部分接口列表access和profile 是在个人中心展示的
		// 所以需要设置为true
		if apiInfo.HideOnSidebar {
			err = RDB.HSet(ctx, hideSidebarKey, apiInfo.Path, true).Err()
			if err != nil {
				continue
			}
		}

		path2roles[apiInfo.Path] = append(path2roles[apiInfo.Path], roleID)
		// 如果开启
		if apiInfo.CanViewDetail {
			pathForDetail := apiInfo.Path + "/:_id"
			path2roles[pathForDetail] = append(path2roles[pathForDetail], roleID)
		}
	}
	return nil
}
