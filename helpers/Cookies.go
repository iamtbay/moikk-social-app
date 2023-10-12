package helpers

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// SET COOKIE
func CreateNewCookie(c *gin.Context, name, value string) error {
	cookie, err := c.Cookie(name)
	if err != nil {
		cookie = "not set"
		c.SetCookie(name, value, 3600, "/", "localhost", false, true)
		return nil
	}
	fmt.Printf("Cookie value: %s \n", cookie)
	return errors.New("cookie already defined")
}

// DELETE COOKIE
func DeleteACookie(c *gin.Context, name string) error {
	c.SetCookie(name, "", -1, "/", "localhost", false, true)
	return nil
}
