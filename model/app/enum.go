package app

// PickUpType 取餐方式
type PickUpType int

// SourceType 订单来源
type SourceType int

const (
	// HelpYourSelf 自提
	HelpYourSelf PickUpType = iota
	// Takeout 外卖
	Takeout
	// DineIn 店内就餐
	DineIn
	// Express 快递
	Express
	// NoNeed 无需取货（虚拟商品，或补差价）
	NoNeed
)

const (
	// WxMiniApp 微信小程序
	WxMiniApp SourceType = iota
	// DouYinMiniApp 抖音小程序
	DouYinMiniApp
	// AlipayMiniApp 支付宝小程序
	AlipayMiniApp
	// MeiTuanMiniApp 美团小程序
	MeiTuanMiniApp
	// PosClient 点餐系统
	PosClient
	// WebClient 网页系统
	WebClient
)
