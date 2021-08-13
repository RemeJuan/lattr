package error_formats

import (
	"fmt"
	"strings"

	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/jinzhu/gorm"
)

func ParseError(err error) error_utils.MessageErr {
	_, ok := err.(*gorm.Errors)
	if !ok {
		if strings.Contains(err.Error(), "no rows in result set") {
			return error_utils.NotFoundError("no record matching given id")
		}
		return error_utils.InternalServerError(fmt.Sprintf("error when trying to save message: %s", err.Error()))
	}

	return error_utils.InternalServerError(fmt.Sprintf("error when processing request: %s", err.Error()))
}
