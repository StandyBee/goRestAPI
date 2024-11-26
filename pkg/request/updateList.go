package reqs

import (
	"errors"
	"fmt"
)

type UpdateListRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ListId      int    `json:"-"`
	UserId      int    `json:"-"`
}

func (req UpdateListRequest) Validate() error {
	fmt.Println(req)
	if req.Title == "" && req.Description == "" || req.Title == "" {
		return errors.New("missing required fields")
	}
	return nil
}
