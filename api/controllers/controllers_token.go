package controllers

import (
	"net/http"
	"strconv"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/services"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/gin-gonic/gin"
)

// CreateToken godoc
// @Summary Create a new token
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param token body domain.Token true "Create Token"
// @Success 200 {object} domain.Token
// @Failure 403 {object} error_utils.MessageErrStruct
// @Failure 404 {object} error_utils.MessageErrStruct
// @Failure 422 {object} error_utils.MessageErrStruct
// @Failure 500 {object} error_utils.MessageErrStruct
// @Failure 501 {object} error_utils.MessageErrStruct
// @Security ApiKeyAuth
// @Security OAuth2Application[token:create]
// @Router /token [post]
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

// GetToken godoc
// @Summary Fetches an existing token by ID
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param id path int true "Token ID"
// @Success 200 {object} domain.Token
// @Failure 403 {object} error_utils.MessageErrStruct
// @Failure 404 {object} error_utils.MessageErrStruct
// @Failure 422 {object} error_utils.MessageErrStruct
// @Failure 500 {object} error_utils.MessageErrStruct
// @Failure 501 {object} error_utils.MessageErrStruct
// @Security ApiKeyAuth
// @Security OAuth2Application[token:read]
// @Router /token/{id} [post]
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

// GetTokens godoc
// @Summary Fetches a list of all available tokens
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Token
// @Failure 403 {object} error_utils.MessageErrStruct
// @Failure 404 {object} error_utils.MessageErrStruct
// @Failure 422 {object} error_utils.MessageErrStruct
// @Failure 500 {object} error_utils.MessageErrStruct
// @Failure 501 {object} error_utils.MessageErrStruct
// @Security ApiKeyAuth
// @Security OAuth2Application[token:read]
// @Router /token/list [post]
func GetTokens(c *gin.Context) {
	result, err := services.AuthService.List()

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ResetToken godoc
// @Summary Resets the token specified by the provided ID
// @Description Using the given ID, a new token is generated with either the default timing
// @Description Or the timing specified in the request payload
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param id path int true "Token ID"
// @Success 200 {array} domain.Token
// @Failure 403 {object} error_utils.MessageErrStruct
// @Failure 404 {object} error_utils.MessageErrStruct
// @Failure 422 {object} error_utils.MessageErrStruct
// @Failure 500 {object} error_utils.MessageErrStruct
// @Failure 501 {object} error_utils.MessageErrStruct
// @Security ApiKeyAuth
// @Security OAuth2Application[token:read]
// @Router /token/{id} [put]
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

// DeleteToken godoc
// @Summary Deletes the specified token
// @Tags Tokens
// @Accept  json
// @Produce  json
// @Param id path int true "Token ID"
// @Success 200 {object} object "{message: "success"}"
// @Failure 403 {object} error_utils.MessageErrStruct
// @Failure 404 {object} error_utils.MessageErrStruct
// @Failure 422 {object} error_utils.MessageErrStruct
// @Failure 500 {object} error_utils.MessageErrStruct
// @Failure 501 {object} error_utils.MessageErrStruct
// @Security ApiKeyAuth
// @Security OAuth2Application[token:read]
// @Router /token/{id} [delete]
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
