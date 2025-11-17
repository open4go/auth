package cst

// ReviewStatus 订单状态
type ReviewStatus int

const (
	// ReviewInit 初始化
	ReviewInit ReviewStatus = iota
	// ReviewApproved 通过
	ReviewApproved
	// ReviewReject 拒绝
	ReviewReject
	// ReviewAccept 接受
	ReviewAccept
)
