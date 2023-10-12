package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/moikk-app/helpers"
	mytypes "github.com/moikk-app/types"
)

type PostsController struct{}

// INIT
func PostControllerinit() *PostsController {
	return &PostsController{}
}

// @Summary Create a New Post
// @Description Create a New Post
// @ID create-new-post
// @Accept json
// @Produce json
// @Param postInfo body mytypes.NewPost true "New Post data"
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /posts/new [post]
// @Tags posts
func (x *PostsController) CreatePost(c *gin.Context) {
	var postInfos *mytypes.NewPost
	//get posts' json.
	err := c.BindJSON(&postInfos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//get user's id from jwt
	userIDFromJWT, err := helpers.GetUserIDFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	postInfos.UserID = userIDFromJWT
	//send it to db.
	newPost, err := postsDB.CreatePost(postInfos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "succesfully inserted",
		"data":    newPost,
	})
}

// @Summary Update a Post By ID
// @Description Update a Post
// @ID update-post
// @Accept json
// @Produce json
// @Param postInfo body mytypes.UpdatePost true "Update Post data"
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /posts/{postID} [patch]
// @Tags posts
func (x *PostsController) UpdatePost(c *gin.Context) {
	//vars
	var updatePostInfo *mytypes.UpdatePost
	err := c.BindJSON(&updatePostInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//get post id.
	postID := c.Param("id")
	//get userid
	userID, err := helpers.GetUserIDFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//db operations
	updatedPost, err := postsDB.UpdatePost(postID, userID, updatePostInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//return

	c.JSON(200, gin.H{
		"success": true,
		"message": "succesfully updated",
		"data":    updatedPost,
	})
}

// / @Summary Delete a Post By ID
// @Description Delete a Post
// @ID delete-post
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /posts/{postID} [delete]
// @Tags posts
func (x *PostsController) DeletePost(c *gin.Context) {
	//get post id
	postID := c.Param("id")
	//get userid from jwt
	userID, err := helpers.GetUserIDFromJWT(c)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//send it to db and delete
	err = postsDB.DeletePost(userID, postID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//return

	c.JSON(200, gin.H{
		"success":   true,
		"message":   "Post deleted succesfully",
		"deletedID": postID,
	})
}

// @Summary Upload a File
// @Description Upload file.
// @ID upload-file
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /posts/upload [post]
// @Tags posts
func (x *PostsController) UploadFile(c *gin.Context) {
	url := c.Request.Host
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	files := form.File["files"]
	var uploadedFiles []string
	for _, file := range files {
		ext := strings.Split(file.Filename, ".")
		filename := fmt.Sprintf("%v.%v", uuid.New(), ext[len(ext)-1])
		if err := c.SaveUploadedFile(file, fmt.Sprintf("./public/posts/%v", filename)); err != nil {
			return
		}
		uploadedFiles = append(uploadedFiles, fmt.Sprintf("%v/files/posts/%v", url, filename))
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": fmt.Sprintf("%v file uploaded", len(files)),
		"data":    uploadedFiles,
	})
}

// @Summary Get Single Post by id
// @Description get single Post by id
// @ID get-a-post
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /posts/single/{id} [get]
// @Tags posts
func (x *PostsController) GetAPost(c *gin.Context) {
	//get post id
	postID := c.Param("id")
	//send it to db
	post, err := postsDB.GetAPost(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "Succesfuly got post",
		"data":    post,
	})
}

// @Summary Get All Post By Followed users
// @Description Get All Post By User's followed users
// @ID get-all-post
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /posts/ [get]
// @Tags posts
func (x *PostsController) GetAllPosts(c *gin.Context) {
	//get user's id from jwt
	userID, _ := helpers.GetUserIDFromJWT(c)
	//page
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	//query for followed user's
	posts, pagination, err := postsDB.GetAllPosts(userID, page)
	//check count of coming data
	var msg string = "Succesfully got posts"
	if len(posts) < 1 {
		msg = "oh no!, there is no post to show."
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//return
	c.JSON(http.StatusOK, gin.H{
		"message":    msg,
		"pagination": pagination,
		"data":       posts,
	})
}

// @Summary Get All Post By User name
// @Description Get User's all posts
// @ID get-all-post-by-username
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /posts/user/{username} [get]
// @Tags posts
func (x *PostsController) GetAllPostsByUserID(c *gin.Context) {
	//get user's id
	username := c.Param("username")
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	//send it to db
	userPosts, pagination, err := postsDB.GetAllPostsByUserID(page, username)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//check count of coming data
	var msg string = "Succesfully got posts"
	if len(userPosts) < 1 {
		msg = "there is no post to show."
	}
	//return data
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": msg,
		"pages":   pagination,
		"data":    userPosts,
	})

}
