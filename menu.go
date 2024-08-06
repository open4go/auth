package auth

import "sort"

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

// ByName implements sort.Interface for []MenuTree based on the Name field.
type ByName []MenuTree

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// SortMenuTree sorts the MenuTree slice by Name.
func (m *MyAccessData) SortMenuTree() {
	sort.Sort(ByName(m.MenuTree))
}
