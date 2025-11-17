package cst

// OrderStatus 订单状态
type OrderStatus int

// OrderFrom 订单来源
type OrderFrom int

const (
	// OrderInit 订单初始化(未完成支付）
	OrderInit OrderStatus = iota
	// OrderPaid 订单已付款
	OrderPaid
	// OrderMaking 订单制作中
	OrderMaking
	// OrderProduceCompleted 订单制作完成
	OrderProduceCompleted
	// OrderCompleted 订单完成 （已经自提/堂食配送完成）
	OrderCompleted
	// 	OrderCancelRejected  拒绝取消
	OrderCancelRejected
	// OrderCancel 订单取消
	OrderCancel
	// OrderCancelApproved 订单取消 （商家同意）
	OrderCancelApproved
	// OrderCancelCompleted 订单取消（完成）已经退款
	OrderCancelCompleted
	// OrderTakeoutPending 外卖 （等待接单 ）
	OrderTakeoutPending
	// OrderTakeoutConfirmed 外卖（已经接单）
	OrderTakeoutConfirmed
	// OrderTakeoutUnConfirmed 外卖（放弃接单）
	OrderTakeoutUnConfirmed
	// OrderTakeoutTake 外卖（已经取货）
	OrderTakeoutTake
	// OrderTakeoutDone 外卖（已经送达）
	OrderTakeoutDone
	// OrderTakeoutComment 订单已经评论
	OrderTakeoutComment
	// OrderClosed 订单关闭 (超时未支付的订单会被关闭)
	OrderClosed
)

const (
	// OrderFromWxMini 微信小程序
	OrderFromWxMini OrderFrom = iota
	// OrderFromAliMini 支付宝小程序
	OrderFromAliMini
	// OrderFromMeiTuan 美团
	OrderFromMeiTuan
	// OrderFromPos 前台点餐
	OrderFromPos
	// OrderFromWeb 来自官网
	OrderFromWeb
)
