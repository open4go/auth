package cst

// ChartType 图形类型
type ChartType int

const (
	// BarChartType 柱状图
	BarChartType ChartType = iota
	// LineChartType 折线图
	LineChartType ChartType = iota
)
