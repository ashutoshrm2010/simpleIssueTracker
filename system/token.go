package system

import (
	"net/http"
	"strings"
	"github.com/zenazn/goji/web"
	"database/sql"
	"errors"
	"github.com/tb/simpleIssueTracker/model"
	"fmt"
)

func (application *Application) ApplyAuth(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Applied authorization filter!")
		accessToken, err := getAccessToken(r)
		if (err != nil) {

			c.Env["AuthFailed"] = true
		} else {

			userContext, _ := GetUserContext(accessToken)
			if (err != nil) {
				c.Env["AuthFailed"] = true
			} else {
				c.Env["AuthFailed"] = false

				c.Env["userContext"] = userContext

			}

		}
		h.ServeHTTP(w, r)

	}
	return http.HandlerFunc(fn)
}

func getAccessToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	var token string

	if authHeader == "" {
		return "", errors.New("Unauthorized request")
	}

	s := strings.SplitN(authHeader, " ", 2)
	if (len(s) != 2 || strings.ToLower(s[0]) != "bearer") && token == "" {
		return "", errors.New("Unauthorized request")
	}
	if (len(s) > 0 && strings.ToLower(s[0]) == "bearer") {
		token = s[1]
	}

	return token, nil
}
func GetUserContext(accessToken string) (*UserContext, error) {

	userName, err := getUserName(accessToken)
	userInfo, err := getUserData1(userName)
	if (err != nil) {
		fmt.Println("err ",err)
	}

	var userContext UserContext

	userContext.UserName = userInfo.UserName
	userContext.Name = userInfo.Name
	userContext.EmailId = userInfo.Email

	userContext.UserId = userInfo.ID
	fmt.Println("userContext.UserId ",userContext.UserId)

	return &userContext, errors.New("Invalid Access Token for given company")

}

func getUserName(accessToken string) (string, error) {
	db, _ := InitMysql()
	defer db.Close()
	var username string

	err := db.QueryRow("SELECT UserName FROM" + " " + "accessToken" + " " + "WHERE AccessToken=?", accessToken).Scan(&username)
	if err != nil {
		return "", errors.New("Invalid accessToken")
	}
	switch {
	case err == sql.ErrNoRows:
		return "", err

	case err != nil:
		return "", err
	default:

		return username, nil
	}

}
func getUserData1(userName string) (*model.User, error) {
	db, _ := InitMysql()
	defer db.Close()
	var userdetails model.User
	var (
		name string
		email string
		id int
		password string
		createdOn string
	)
	fmt.Println("userName ",userName)
	rows, err := db.Query("select * from" + " " + "user" + " " + "where UserName = ?", userName)
	if err != nil {
		fmt.Println("err ",err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name,&userName, &email, &password,&createdOn)
		if err != nil {
			fmt.Println("err ",err)
		}

	}
	fmt.Println("id ",id)
	userdetails.ID = id
	userdetails.Name = name
	userdetails.Email = email
	userdetails.Password = password
	err = rows.Err()
	if err != nil {
		fmt.Println("err ",err)
	}

	return &userdetails, err

}

