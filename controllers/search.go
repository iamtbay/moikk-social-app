package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchStruct struct{}

func SearchControllerInit() *SearchStruct {
	return &SearchStruct{}
}

// @Summary Search Posts by Query
// @Description Search Posts by providing a Query
// @ID search-posts-by-query
// @Produce json
// @Success 200 {object} any
// @Failure 404 {string} Not Found
// @Router /search/posts [get]
// @Tags search
func (x *SearchStruct) SearchPosts(c *gin.Context) {
	//get search query from
	keyword := c.Query("q")
	pageStr := c.Query("page")
	fmt.Println(pageStr)
	page, _ := strconv.Atoi(pageStr)
	//db operations
	posts, pagination, err := searchDB.SearchPosts(keyword, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//
	msg := "Succesful"
	if len(posts) < 1 {
		msg = "Oh couldn't find anything"
	}
	//return
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    msg,
		"pagination": pagination,
		"data":       posts,
	})
}

// @Summary Search Users by Query
// @Description Search Users by providing a Query
// @ID search-users-by-query
// @Produce json
// @Success 200 {object} any
// @Failure 404 {string} Not Found
// @Router /search/users [get]
// @Tags search
func (x *SearchStruct) SearchUsers(c *gin.Context) {
	//get url
	//get search query from
	keyword := c.Query("q")
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	//db operations
	users, pagination, err := searchDB.SearchUsers(keyword, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	msg := "Succesful"
	if len(users) < 1 {
		msg = "Oh couldn't find anyone"
	}
	//return
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    msg,
		"pagination": pagination,
		"data":       users,
	})

}
