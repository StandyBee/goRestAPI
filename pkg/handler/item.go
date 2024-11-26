package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorestAPI/pkg/request/item"
	"net/http"
	"strconv"
)

func (h *Handler) getItems(c *gin.Context) {
	listId, err := strconv.Atoi(c.Param("list_id"))
	if err != nil {
		return
	}

	items, err := h.services.TodoItem.GetListItems(listId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}

func (h *Handler) getItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		return
	}

	item, err := h.services.TodoItem.GetItem(itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    item,
	})
}

func (h *Handler) createItem(c *gin.Context) {
	var request item.CreateItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"success": false,
		})
	}

	err := request.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	listId, err := strconv.Atoi(c.Param("list_id"))
	if err != nil {
		return
	}

	itemId, err := h.services.TodoItem.CreateItem(request, listId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"success": false,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"itemId":  itemId,
	})
}

func (h *Handler) updateItem(c *gin.Context) {
	var request item.UpdateItemRequest

	listId, err := strconv.Atoi(c.Param("list_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err = c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err = request.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ERROR: " + err.Error(),
		})
		return
	}

	err = h.services.TodoItem.UpdateItem(&request, listId, itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update list: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("successfully updated item %d", itemId),
	})
}

func (h *Handler) deleteItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		return
	}

	err = h.services.TodoItem.DeleteItem(itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"success": false,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("successfully deleted item %d", itemId),
	})
}
