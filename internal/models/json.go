package models

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type Response struct {
	Data       any             `json:"data,omitempty"`
	Error      *ErrorResponse  `json:"error,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

type PaginationMeta struct {
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
	Total         int64  `json:"total"`
	NextPageToken string `json:"next_page_token"`
	PrevPageToken string `json:"prev_page_token"`
	TotalPages    int    `json:"total_pages"`
	PageSize      int    `json:"page_size"`
	HasNext       bool   `json:"has_next"`
	HasPrev       bool   `json:"has_prev"`
}

// ERROR sends an error response in JSON format with the specified status code and message.
func ERROR(c *gin.Context, statusCode int, err error) {
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(statusCode)
	c.JSON(statusCode, Response{
		Data: nil,
		Error: &ErrorResponse{
			Error: struct {
				Code    int    "json:\"code\""
				Message string "json:\"message\""
			}{},
		},
	})

}

// JSON sends a JSON response with the specified status code and data.
func JSON(c *gin.Context, statusCode int, data Response) {
	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(statusCode)
	c.JSON(statusCode, data)
	c.Abort()
}

type WorkflowResponse struct {
	WorkflowID  string `json:"workflow_id"`
	RunID       string `json:"run_id"`
	ReferenceID string `json:"reference_id"`
	Message     string `json:"message"`
}
