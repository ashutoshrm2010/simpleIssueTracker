package services
import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
	"encoding/json"
	"github.com/tb/simpleIssueTracker/model"
	"github.com/pkg/errors"
	"github.com/tb/simpleIssueTracker/system"
	"fmt"
)

func SignUp(userDetails model.User) ([]byte, error) {
	fmt.Println("userdata ",userDetails)
	db, _ := system.InitMysql()
	name := userDetails.Name
	if name == "" {
		return nil, errors.New("name is missing")
	}
	userName := userDetails.UserName
	if userName == "" {
		return nil, errors.New("username is missing")
	}
	email := userDetails.Email
	if email == "" {
		return nil, errors.New("emailId is missing")
	}
	password := userDetails.Password
	if password == "" {
		return nil, errors.New("password is missing")
	}
	var user string
	if name != ""&&password != "" {
		err := db.QueryRow("SELECT username FROM" + " " + "user" + " " + "WHERE UserName=?", userName).Scan(&user)
		//fmt.Println(err)
		if err == nil {
			fmt.Println(err)
			return nil, errors.New("user alraedy exist")
		} else if err == sql.ErrNoRows {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			createdOn := time.Now()
			_, err = db.Exec("INSERT INTO" + " " + "user" + "(Name,UserName,Email,Password,CreatedOn) VALUES( ?, ?, ?, ?,?)", name, userName, email, hashedPassword,createdOn)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			response := make(map[string]interface{})
			response["message"] = "User Created"
			finalResponse, _ := json.Marshal(response)
			return finalResponse, nil
		}

	}

	return nil, nil

}
