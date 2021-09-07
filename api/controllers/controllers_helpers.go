package controllers

import (
	"os"
	"strconv"
	"strings"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/services"
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

func TokenCreateMiddleWare(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCreate := os.Getenv("ENABLE_CREATE")

		if tokenCreate == "OPEN" {
			c.Next()
		} else if tokenCreate == "SCOPED" {
			isValid := AuthenticateToken(c, requiredScope)

			if isValid {
				c.Next()
			} else {
				tokenCreateErr(c)
			}
		} else {
			tokenCreateErr(c)
		}
	}
}

func AuthenticateMiddleware(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isValid := AuthenticateToken(c, requiredScope)

		if isValid {
			c.Next()
		} else {
			err := error_utils.ForbiddenError("Invalid or missing token")
			c.JSON(err.Status(), err)
			c.Abort()
		}
	}
}

func AuthenticateToken(c *gin.Context, requiredScope string) bool {
	jwtToken := ExtractToken(c)

	return services.AuthService.ValidateToken(&domain.Token{Token: jwtToken}, requiredScope)
}

func ExtractToken(c *gin.Context) string {
	tokenString := c.GetHeader("Authorization")

	tokenString = stripTokenPrefix(tokenString)

	return tokenString
}

func stripTokenPrefix(tok string) string {
	// split token to 2 parts
	tokenParts := strings.Split(tok, " ")

	if len(tokenParts) < 2 {
		return tokenParts[0]
	}

	return tokenParts[1]
}

func tokenCreateErr(c *gin.Context) {
	err := error_utils.NotImplementedError("Token creation not enabled for this app")
	c.JSON(err.Status(), err)
	c.Abort()
}
