package presentation

import "time"

type ResponseBase struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

func NewResponseBase() ResponseBase {
	return ResponseBase{
		Timestamp: time.Now().UnixMilli(),
	}
}

func (ResponseBase) Failed(status int, message string, requestID string) ResponseBase {
	return ResponseBase{
		Status:    status,
		Message:   message,
		Timestamp: time.Now().UnixMilli(),
		RequestID: requestID,
	}
}

func (ResponseBase) Success(message string, data interface{}, requestID string) ResponseBase {
	return ResponseBase{
		Status:    200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
		RequestID: requestID,
	}
}

// PaginationMetadata contains pagination information
type PaginationMetadata struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	ResponseBase
	Metadata *PaginationMetadata `json:"metadata,omitempty"`
}

func (ResponseBase) SuccessPaginated(message string, data interface{}, meta *PaginationMetadata, requestID string) PaginatedResponse {
	return PaginatedResponse{
		ResponseBase: ResponseBase{
			Status:    200,
			Message:   message,
			Data:      data,
			Timestamp: time.Now().UnixMilli(),
			RequestID: requestID,
		},
		Metadata: meta,
	}
}

type ExampleResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
