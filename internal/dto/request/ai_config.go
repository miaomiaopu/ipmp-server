package request

type UpdateAIConfigRequest struct {
	Provider string `json:"provider" binding:"omitempty,oneof=deepseek openai claude custom"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
	BaseURL  string `json:"base_url"`
}

type TestAIConnectionRequest struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
	BaseURL  string `json:"base_url"`
}
