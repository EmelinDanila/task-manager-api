package models

// TokenResponse represents a successful login response
type TokenResponse struct {
	Token string `json:"token"` // JWT-токен
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"` // Сообщение об ошибке
}

// TaskResponse represents a single task response
type TaskResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      uint   `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// TaskListResponse represents a list of tasks response
type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}
