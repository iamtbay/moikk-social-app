package controllers

import (
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/moikk-app/helpers"
	mytypes "github.com/moikk-app/types"
)

type AuthController struct{}

// INIT
func AuthControllerInit() *AuthController {
	return &AuthController{}
}

// @Summary Register For New User
// @Description Register For New User
// @ID register-user
// @Accept json
// @Produce json
// @Param userInfo body mytypes.Register true "User register data"
// @Success 200 {object} map[string]interface{} "Succesful register"
// @Failure 400 {object} map[string]interface{} "Error"
// @Router /auth/register [post]
// @Tags auth
func (x *AuthController) Register(c *gin.Context) {
	var userInfo *mytypes.Register
	err := c.BindJSON(&userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check mail valid or not
	_, err = mail.ParseAddress(userInfo.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a valid email adress."})
		return
	}
	//hash pass
	hashedPass, err := helpers.PasswordHasher(userInfo.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//
	userInfo.Password = hashedPass
	registeredUserID, err := authDB.Register(userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Succesfully registered",
		"userID":  registeredUserID,
	})
}

// @Summary User Login
// @Description Login for User
// @ID login-user
// @Accept json
// @Produce json
// @Param userInfo body mytypes.Login true "User login data"
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /auth/login [post]
// @Tags auth
func (x *AuthController) Login(c *gin.Context) {
	var loginInfo *mytypes.Login
	err := c.BindJSON(&loginInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//

	userInfos, err := authDB.Login(loginInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = helpers.PasswordChecker(loginInfo.Password, userInfos.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//create jwt
	jwtToken, err := helpers.CreateJWT(userInfos.ID, userInfos.Name, userInfos.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"ot":    "5",
		})
		return
	}
	//set cookie
	err = helpers.CreateNewCookie(c, "accessToken", jwtToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//return
	c.JSON(200, gin.H{
		"success": true,
		"message": "Succesfully logged in.",
		"user": gin.H{
			"username": userInfos.Name,
			"email":    userInfos.Email,
			"location": userInfos.Location,
		},
	},
	)

}

// @Summary Logout for User
// @Description Logout For User
// @ID logout-user
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /auth/logout [post]
// @Tags auth
func (x *AuthController) Logout(c *gin.Context) {
	helpers.DeleteACookie(c, "accessToken")
	c.JSON(200, gin.H{
		"success": true,
		"message": "Succesfully logout",
	},
	)
}

// @Summary Update Infos For User
// @Description Update Infos For User
// @ID update-user
// @Accept json
// @Produce json
// @Param userInfo body mytypes.UserInfos true "User Update Info data"
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /auth/updateUser [patch]
// @Tags auth
func (x *AuthController) UpdateUserInfo(c *gin.Context) {
	//get jwt
	token, _ := c.Cookie("accessToken")
	//check jwt
	jwtUser, err := helpers.JWTParse(token)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Something went wrong, please logout and login again.",
		})
		c.Abort()
		return
	}
	//
	var newUserInfos *mytypes.UserInfos
	err = c.BindJSON(&newUserInfos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	//check mail valid or not
	_, err = mail.ParseAddress(newUserInfos.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a valid email adress."})
		return
	}
	//USERID SHOULD CHECK WITH JWT.
	hashedPass, err := helpers.PasswordHasher(newUserInfos.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	newUserInfos.Password = hashedPass
	//db operations
	err = authDB.UpdateUserInfo(jwtUser.ID, newUserInfos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Succesfully updated",
	})
}

// @Summary Upload Profile Photo For User
// @Description Upload Profile Photo For User
// @ID upload-profilePhoto-user
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /auth/profilePhoto [post]
// @Tags auth
func (x *AuthController) UpdateProfilePhoto(c *gin.Context) {
	url := c.Request.Host
	file, _ := c.FormFile("photo")
	ext := strings.Split(file.Filename, ".")
	filename := fmt.Sprintf("%v.%v", uuid.New(), ext[len(ext)-1])
	//upload file
	if err := c.SaveUploadedFile(file, fmt.Sprintf("./public/profileImages/%v", filename)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fullPhotoPath := fmt.Sprintf("%v/files/profileImages/%v", url, filename)
	//get userid from jwt
	userID, _ := helpers.GetUserIDFromJWT(c)
	//save to db
	oldPhotoPath, err := authDB.UpdateProfilePhoto(userID, fullPhotoPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//if user has old photo delete it
	if oldPhotoPath != "nophoto" {
		os.Remove(oldPhotoPath)
	}
	//return json
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "successful",
		"data":    fullPhotoPath,
	})

}
