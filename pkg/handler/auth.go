package handler

import (
	"github.com/gin-gonic/gin"
	"gorestAPI"
	"net/http"
)

// @Summary SignUp
// @Tags Auth
// @Description Create user account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body gorestAPI.User true "account info"
// @Success 200 {integer} integer 1
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input gorestAPI.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	err := c.BindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	token, err := h.services.Authorization.Login(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"token":   token,
	})
}
