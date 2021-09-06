package controllers

import (
	"net/http"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/services"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/gin-gonic/gin"
)

func CreateToken(c *gin.Context) {
	var token domain.Token

	if err := c.ShouldBindJSON(&token); err != nil {
		thErr := error_utils.UnprocessableEntityError("invalid json body")
		c.JSON(thErr.Status(), thErr)
		return
	}

	result, err := services.AuthService.Create(&token)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetToken(c *gin.Context) {
	tkId, _ := GetIdFromParam(c)

	result, getErr := services.AuthService.Get(*tkId)

	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func GetTokens(c *gin.Context) {
	result, err := services.AuthService.List()

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func ResetToken(c *gin.Context) {
	tkId, _ := GetIdFromParam(c)

	var token domain.Token
	token.Id = *tkId

	result, err := services.AuthService.Reset(&token)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func DeleteToken(c *gin.Context) {
	tkId, _ := GetIdFromParam(c)

	if err := services.AuthService.Delete(*tkId); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
