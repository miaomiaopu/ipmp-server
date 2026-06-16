package request

type GenerateAIReportRequest struct {
	Content string `json:"content"`
}

type SummarizeRequest struct {
	Content string `json:"content" binding:"required"`
}
