package cst

// Permission is a type representing various permission flags.
type Permission uint

// Define the constants associated with the Permission type.
const (
	Create  Permission = 1 << iota // 1 << 0 = 1
	Delete                         // 1 << 1 = 2
	Update                         // 1 << 2 = 4
	Read                           // 1 << 3 = 8
	List                           // 1 << 4 = 16
	Show                           // 1 << 5 = 32
	Disable                        // 1 << 6 = 64
)

// PermissionNames 用于映射权限的二进制位到权限名称
var PermissionNames = map[Permission]string{
	Create:  "Create",
	Delete:  "Delete",
	Update:  "Update",
	Read:    "Read",
	List:    "List",
	Show:    "Show",
	Disable: "Disable",
}
