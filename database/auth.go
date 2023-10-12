package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	mytypes "github.com/moikk-app/types"
)

type AuthDB struct{}

//INIT

func AuthDBInit() *AuthDB {
	return &AuthDB{}
}

// REGISTER
func (x *AuthDB) Register(userInfo *mytypes.Register) (string, error) {
	//
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//check email is exists
	var isEmailExists bool
	conn.QueryRow(ctx, "select exists(select * from users where email=$1)", userInfo.Email).Scan(&isEmailExists)
	if isEmailExists {
		return "", errors.New("this e-mail in already in use please try another e-mail")
	}
	//check username is exists
	var isUsernameExists bool
	conn.QueryRow(ctx, "select exists(select * from users where username=$1)", userInfo.Username).Scan(&isUsernameExists)
	if isUsernameExists {
		return "", errors.New("this username in already in use please try another username")
	}
	//query
	query := `INSERT INTO 
				users(name,email,password,location,username)
			VALUES($1,$2,$3,$4,$5)
			RETURNING id
				`
	var lastInsertedID string
	err := conn.QueryRow(ctx, query,
		&userInfo.Name,
		&userInfo.Email,
		&userInfo.Password,
		&userInfo.Location,
		&userInfo.Username,
	).Scan(&lastInsertedID)

	if err != nil {
		return "", err
	}
	//return
	return lastInsertedID, nil
}

// LOGIN
func (x *AuthDB) Login(loginInfo *mytypes.Login) (*mytypes.UserInfos, error) {
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	query := `SELECT * FROM users WHERE email=$1`
	var userInfoFromDB mytypes.UserInfos
	err := conn.QueryRow(ctx, query, loginInfo.Email).Scan(
		&userInfoFromDB.ID,
		&userInfoFromDB.Name,
		&userInfoFromDB.Username,
		&userInfoFromDB.Email,
		&userInfoFromDB.Password,
		&userInfoFromDB.Location,
		&userInfoFromDB.ProfilePhoto,
	)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &userInfoFromDB, nil
}

//UPDATEUSER

func (x *AuthDB) UpdateUserInfo(userID string, newUserInfos *mytypes.UserInfos) error {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//check jwt's userid and id coming json if not equal send an error.
	if userID != newUserInfos.ID {
		return errors.New("only the authorized user or user himself/hersef can do that this operation")
	}
	//
	query := `UPDATE users
				SET
				name=$1,
				email=$2,
				password=$3,
				location=$4
				username=$5
			WHERE id = $6
				`
	_, err := conn.Exec(ctx, query,
		&newUserInfos.Name,
		&newUserInfos.Email,
		&newUserInfos.Password,
		&newUserInfos.Location,
		&newUserInfos.Username,
		&newUserInfos.ID,
	)
	if err != nil {
		return errors.New("something went wrong try again later")
	}
	return nil
}

//update profile photo

func (x *AuthDB) UpdateProfilePhoto(userID, photoPath string) (string, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//op
	var currentPhotoPath string
	query := `select profile_photo from users where id=$1`
	err := conn.QueryRow(ctx, query, userID).Scan(&currentPhotoPath)
	if currentPhotoPath != "nophoto" {
		getImgName := strings.Split(currentPhotoPath, "/")
		currentPhotoPath = fmt.Sprintf("./public/profileImages/%v", getImgName[len(getImgName)-1])
	}
	if err != nil {
		fmt.Println("error here",userID,err)
		return "", err
	}
	//main query for update
	query = `update users set profile_photo=$1 where id=$2`
	_, err = conn.Exec(ctx, query, photoPath, userID)
	if err != nil {
		fmt.Println("error here 2")
		return "", err
	}
	return currentPhotoPath, nil
}
