package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	clientToken := c.Request.Header.Get("token")

	user, _ := ValidateToken(clientToken)
	userType := user.UserType
	err = nil
	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	clientToken := c.Request.Header.Get("token")

	user, _ := ValidateToken(clientToken)
	userType := user.UserType
	uid := user.Uid
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}
