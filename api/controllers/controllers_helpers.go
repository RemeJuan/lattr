package controllers

import (
	"strconv"

	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/gin-gonic/gin"
)

func GetIdFromParam(c *gin.Context) (*int64, error) {
	paramId := GetParam(c, "id")
	tkId, err := strconv.ParseInt(paramId, 10, 64)

	if err != nil {
		theErr := error_utils.UnprocessableEntityError("unable to parse ID")
		c.JSON(theErr.Status(), theErr)
		return nil, err
	}

	return &tkId, err
}

func GetParam(c *gin.Context, paramName string) string {
	return c.Params.ByName(paramName)
}
