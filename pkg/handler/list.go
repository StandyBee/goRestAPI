package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorestAPI/pkg/request"
	"net/http"
	"strconv"
)

func (h *Handler) getLists(c *gin.Context) {
	userId := getUserId(c)

	if userId == 0 {
		return
	}

	lists, _ := h.services.TodoList.GetUserLists(userId)
	c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"user_lists": lists,
	})
}

type GetListResponse struct {
	Success bool        `json:"success"`
	List    interface{} `json:"list"`
}

func (h *Handler) getList(c *gin.Context) {

	userId := getUserId(c)

	if userId == 0 {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized",
		})
		return
	}

	listId, err := strconv.Atoi(c.Param("list_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	list, err := h.services.TodoList.GetListById(listId, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	response := GetListResponse{
		Success: true,
		List:    list,
	}
	c.JSON(http.StatusOK, response)
}

type CreateListRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func (h *Handler) createList(c *gin.Context) {
	var request CreateListRequest
	userId := getUserId(c)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	listId, err := h.services.TodoList.CreateList(request.Title, request.Description, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"list_id": listId,
	})
}

func (h *Handler) updateList(c *gin.Context) {
	var request reqs.UpdateListRequest

	userId := getUserId(c)
	listId, err := strconv.Atoi(c.Param("list_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err = c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	request.ListId = listId
	request.UserId = userId

	err = request.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "ERROR: " + err.Error(),
		})
		return
	}

	err = h.services.TodoList.UpdateList(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "failed to update list: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("successfully updated list %d", listId),
	})
}

func (h *Handler) deleteList(c *gin.Context) {
	userId := getUserId(c)
	listId, err := strconv.Atoi(c.Param("list_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err = h.services.TodoList.DeleteList(listId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "failed to delete list: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("successfully deleted list %d", listId),
	})
}

func getUserId(c *gin.Context) int {
	userIdRaw, ok := c.Get(userCtx)
	if !ok {
		return 0
	}
	return userIdRaw.(int)
}
