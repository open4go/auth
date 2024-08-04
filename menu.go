package auth

// MenuTree 主菜单sidebar 列表
type MenuTree struct {
	// 主菜单名称
	Name string `json:"name"`
	// 子菜单列表
	SubMenu []string `json:"sub_menu"`
}

// MyAccessData 访问数据结构定义
type MyAccessData struct {
	ID         string     `json:"id"`
	PathsAllow []string   `json:"paths_allow"`
	MenuTree   []MenuTree `json:"menu_tree"`
}
