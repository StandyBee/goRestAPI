package item

import "errors"

type CreateItemRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (request *CreateItemRequest) Validate() error {
	if len(request.Title) < 2 {
		return errors.New("title is too short")
	}
	return nil
}
