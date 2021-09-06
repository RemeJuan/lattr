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

func TokenCreateMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		isEnabled, _ := strconv.ParseBool(os.Getenv("ENABLE_CREATE"))

		if isEnabled {
			c.Next()
		} else {
			err := error_utils.NotImplementedError("Token creation not enabled for this app")
			c.JSON(err.Status(), err)
			c.Abort()
		}
	}
}

func AuthenticateJWTMiddleware(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isValid := AuthenticateJWTToken(c, requiredScope)

		if isValid {
			c.Next()
		} else {
			err := error_utils.ForbiddenError("Invalid or missing token")
			c.JSON(err.Status(), err)
			c.Abort()
		}
	}
}

func AuthenticateJWTToken(c *gin.Context, requiredScope string) bool {
	jwtToken := ExtractJWTToken(c)

	return services.AuthService.ValidateToken(&domain.Token{Token: jwtToken}, requiredScope)
}

func ExtractJWTToken(c *gin.Context) string {
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
