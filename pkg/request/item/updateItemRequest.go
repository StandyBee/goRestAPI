package item

import "errors"

type UpdateItemRequest struct {
	Title       *string `json:"title"`       // Указатель на строку
	Description *string `json:"description"` // Указатель на строку
	Done        *bool   `json:"done"`        // Указатель на bool
}

func (request *UpdateItemRequest) Validate() error {
	if request.Title != nil && len(*request.Title) < 2 {
		return errors.New("title is too short")
	}
	return nil
}
