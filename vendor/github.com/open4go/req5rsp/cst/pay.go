package cst

// ChannelType 支付渠道类型
type ChannelType int

const (
	// WeChatPay 微信支付
	WeChatPay ChannelType = iota
	// BalancePay 余额支付
	BalancePay
	// KalaPay 卡拉卡支付
	KalaPay
	// AliPay 支付宝
	AliPay
	// CashPay 现金支付
	CashPay
	// UniPay 银联支付
	UniPay
)

type PayMethod int

const (
	// WxMiniPay 微信小程序支付
	WxMiniPay PayMethod = iota
	// WxScanQRCode 微信扫码支付
	WxScanQRCode
	// AliMiniPay 支付宝小程序支付
	AliMiniPay
	// AliScanQRCode 支付宝扫码支付
	AliScanQRCode
)

func (p PayMethod) String() string {
	return [...]string{"微信小程序支付", "微信扫码支付", "支付宝小程序支付", "支付宝扫码支付"}[p]
}

type PayStatus int

const (
	// UnPaid 未支付
	UnPaid PayStatus = iota
	// PaidDone 已支付
	PaidDone
	// Refunding 退款中
	Refunding
	// Refunded 已退款
	Refunded
)

func (p PayStatus) String() string {
	return [...]string{"未支付", "已支付", "退款中", "已退款"}[p]
}

type PayUnit int

const (
	// PayByFen 以分支付
	PayByFen PayUnit = 1
	// PayByYuan 以元支付
	PayByYuan PayUnit = 100
)
