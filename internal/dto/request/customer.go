package request

// CreateCustomerRequest 创建客户请求
type CreateCustomerRequest struct {
	CustomerCode  string `json:"customer_code" binding:"required,min=1,max=32"`
	Name          string `json:"name" binding:"required,min=1,max=256"`
	ContactPerson string `json:"contact_person" binding:"omitempty,max=128"`
	ContactPhone  string `json:"contact_phone" binding:"omitempty,max=64"`
	ContactEmail  string `json:"contact_email" binding:"omitempty,max=256"`
	Address       string `json:"address" binding:"omitempty,max=512"`
	Notes         string `json:"notes" binding:"omitempty,max=2048"`
}

// UpdateCustomerRequest 更新客户请求
type UpdateCustomerRequest struct {
	Name          *string `json:"name" binding:"omitempty,min=1,max=256"`
	CustomerCode  *string `json:"customer_code" binding:"omitempty,min=1,max=32"`
	ContactPerson *string `json:"contact_person" binding:"omitempty,max=128"`
	ContactPhone  *string `json:"contact_phone" binding:"omitempty,max=64"`
	ContactEmail  *string `json:"contact_email" binding:"omitempty,max=256"`
	Address       *string `json:"address" binding:"omitempty,max=512"`
	Notes         *string `json:"notes" binding:"omitempty,max=2048"`
	Status        *string `json:"status" binding:"omitempty,oneof=active inactive"`
}
