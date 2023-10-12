package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	mytypes "github.com/moikk-app/types"
)

type FollowsDB struct{}

// init
func FollowsDBInit() *FollowsDB {
	return &FollowsDB{}
}

// CAN WE ADD HERE GET FOLLOWERS AND FOLLOWED USERS  ?
// for followers
func (x *FollowsDB) GetUserFollowers(username string) ([]*mytypes.GetFollowers, error) {
	//context
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//Query
	query := `
	SELECT
		users.username
	FROM
		follows
		LEFT JOIN users ON users.id = follower_user_id
	WHERE
		follows.followed_user_id IN(
			SELECT
				id FROM users
			WHERE
				username = $1);
			`
	rows, err := conn.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}

	//add array list
	var followers []*mytypes.GetFollowers
	for rows.Next() {
		var follower mytypes.GetFollowers
		err := rows.Scan(
			&follower.Username,
		)
		if err != nil {
			return nil, err
		}
		followers = append(followers, &follower)
	}
	return followers, nil
}

// for followeds
func (x *FollowsDB) GetUserFolloweds(username string) ([]*mytypes.GetFollowers, error){

	//context
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//Query
	query := `
	SELECT
		users.username
	FROM
		follows
		LEFT JOIN users ON users.id = followed_user_id
	WHERE
		follows.follower_user_id IN(
			SELECT
				id FROM users
			WHERE
				username = $1);
			`
	rows, err := conn.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}

	//add array list
	var followeds []*mytypes.GetFollowers
	for rows.Next() {
		var followed mytypes.GetFollowers
		err := rows.Scan(
			&followed.Username,
		)
		if err != nil {
			return nil, err
		}
		followeds = append(followeds, &followed)
	}
	return followeds, nil

}

// FOLLOW A USER
func (x *FollowsDB) FollowUser(userID, targetUsername string) error {
	//CTX
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//QUERY
	//find user id
	var targetID string
	query := `select id from users where username=$1`
	err := conn.QueryRow(ctx, query, targetUsername).Scan(&targetID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("there is no user with this username, check the username")
		}
		return err
	}
	//self control
	fmt.Println("user", userID, "target", targetID)
	if userID == targetID {
		return errors.New("it's you man, you can't follow yourself")
	}
	//
	//check user following target or not
	var isExist bool
	query = `
	SELECT EXISTS (SELECT 1 FROM follows WHERE follower_user_id=$1 and followed_user_id=$2)
	`
	err = conn.QueryRow(ctx, query, userID, targetID).Scan(&isExist)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("user already following the target user")
	}
	//if not save it
	var newFollowedUser mytypes.GetFollows
	query = `INSERT INTO follows (follower_user_id,followed_user_id) VALUES($1,$2) returning *`
	err = conn.QueryRow(ctx, query, userID, targetID).Scan(
		&newFollowedUser.ID,
		&newFollowedUser.FollowerID,
		&newFollowedUser.FollowedID,
		&newFollowedUser.CreatedAt,
	)
	if err != nil {
		return err
	}
	//RETURN
	return nil
}

// UNFOLLOW A USER
func (x *FollowsDB) UnfollowUser(userID, targetUsername string) error {
	//CTX
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//QUERY
	//find user id
	var targetID string
	query := `select id from users where username=$1`
	err := conn.QueryRow(ctx, query, targetUsername).Scan(&targetID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("there is no user with this username, check the username")
		}
		return err
	}
	//self control
	fmt.Println("user", userID, "target", targetID)
	if userID == targetID {
		return errors.New("it's you man, you can't follow yourself")
	}
	//check user following him or not?
	var isExist bool
	query = `SELECT EXISTS (SELECT 1 FROM follows WHERE follower_user_id=$1 and followed_user_id=$2)`
	err = conn.QueryRow(ctx, query, userID, targetID).Scan(&isExist)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("user doesn't follow the target user")
	}
	//IF NOT DELETE USER'S DATA ON DB
	query = `
			DELETE FROM 
				follows 
			WHERE  
			follower_user_id=$1 
		AND 
			followed_user_id=$2
			`
	var rowCount int
	err = conn.QueryRow(ctx, query, userID, targetID).Scan(&rowCount)
	if rowCount != 0 {
		return err
	}
	//RETURN
	return nil

}
