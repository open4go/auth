package cst

// WSMessageType 消息类型
type WSMessageType int

const (
	// DailyOrderNumber 订单笔数
	DailyOrderNumber WSMessageType = iota
	// DailyOrderPriceSum 日交易金额
	DailyOrderPriceSum
	// WeeklyOrderPriceSum 周交易金额
	WeeklyOrderPriceSum
	// MonthlyOrderPriceSum 月交易金额
	MonthlyOrderPriceSum
	// TotalOrderPriceSum 月交易金额
	TotalOrderPriceSum
	// DailyNewUser 每日新增用户
	DailyNewUser
	// WeeklyNewUser 每周新增用户
	WeeklyNewUser
	// MonthlyNewUser 每月新增用户
	MonthlyNewUser
	// PendingOrderNews 新的待支付订单信息
	PendingOrderNews
	// FeedbackReviews 反馈评论数据
	FeedbackReviews
)

const (
	// WSDataChannel 数据推送
	WSDataChannel = "ws:data:channel:0"
)
