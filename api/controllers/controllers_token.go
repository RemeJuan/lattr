package controllers

import (
	"net/http"
	"strconv"

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
	paramId := GetParam(c, "id")
	tkId, err := strconv.ParseInt(paramId, 10, 64)

	if err != nil {
		theErr := error_utils.UnprocessableEntityError("unable to parse ID")
		c.JSON(theErr.Status(), theErr)
		return
	}

	result, getErr := services.AuthService.Get(tkId)

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
	paramId := GetParam(c, "id")
	tkId, err := strconv.ParseInt(paramId, 10, 64)

	if err != nil {
		theErr := error_utils.UnprocessableEntityError("unable to parse ID")
		c.JSON(theErr.Status(), theErr)
		return
	}

	var token domain.Token
	token.Id = tkId

	result, resErr := services.AuthService.Reset(&token)
	if resErr != nil {
		c.JSON(resErr.Status(), resErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func DeleteToken(c *gin.Context) {
	paramId := GetParam(c, "id")
	tkId, parseErr := strconv.ParseInt(paramId, 10, 64)

	if parseErr != nil {
		theErr := error_utils.UnprocessableEntityError("unable to parse ID")
		c.JSON(theErr.Status(), theErr)
		return
	}

	if err := services.AuthService.Delete(tkId); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
