package services
import (
	"golang.org/x/crypto/bcrypt"
	"time"
	"encoding/json"
	"github.com/tb/simpleIssueTracker/model"
	"github.com/pkg/errors"
	"github.com/tb/simpleIssueTracker/system"
	"fmt"
	"database/sql"
	"github.com/pborman/uuid"
	"net/smtp"
)

func SignUp(userDetails model.User) ([]byte, error) {
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
		err := db.QueryRow("SELECT UserName FROM" + " " + "user" + " " + "WHERE UserName=?", userName).Scan(&user)
		if err == nil {
			fmt.Println("err ",err)
			return nil, errors.New("user already exist")
		} else if err == sql.ErrNoRows{
			stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS" + " " + "user" + "(id int NOT NULL AUTO_INCREMENT, Name varchar(50) NOT NULL,UserName varchar(50) NOT NULL,Email varchar(50) NOT NULL," +
				"Password varchar(120) NOT NULL,CreatedOn varchar(100) NOT NULL, PRIMARY KEY (id));")
			if err != nil {
				fmt.Println(err.Error())
			}
			_, err = stmt.Exec()
			if err != nil {
				fmt.Println(err.Error())
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			createdOn := time.Now().UTC()
			_, err = db.Exec("INSERT INTO" + " " + "user" + "(Name,UserName,Email,Password,CreatedOn) VALUES( ?, ?, ?, ?,?)", name, userName, email,hashedPassword,createdOn )
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
func ValidateUserNameandPassword(userDetails map[string]interface{}) ([]byte, error) {
	db, _ := system.InitMysql()
	defer db.Close()
	username := userDetails["userName"].(string)
	password := userDetails["password"].(string)

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT UserName, Password FROM" + " " + "user" + " " + "WHERE UserName=?", username).Scan(&databaseUsername, &databasePassword)
	if err != nil {
		return nil, errors.New("Invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))

	if err != nil {
		return nil, errors.New("Invalid username or password")
	} else {
		fmt.Println("login success")
		stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS"+" "+"accessToken"+"(AccessToken varchar(100) NOT NULL, UserName varchar(50) NOT NULL,PRIMARY KEY (AccessToken));")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec()
		if err != nil {
			fmt.Println(err.Error())
		}
		var accesstoken string
		err = db.QueryRow("select AccessToken from" + " " + "accessToken" + " " + "where UserName= ?", username).Scan(&accesstoken)
		if err == sql.ErrNoRows {
			accessToken := uuid.New()
			_, err = db.Exec("INSERT INTO" + " " + "accessToken" + "(accessToken,userName) VALUES(?, ?)", accessToken, username)
			if err != nil {
				return nil, err
			}
			response := make(map[string]interface{})
			response["accessToken"] = accessToken
			response["userName"] = username
			finalResponse, _ := json.Marshal(response)
			return finalResponse, nil

		} else {
			response := make(map[string]interface{})
			response["accessToken"] = accesstoken
			response["userName"] = username
			finalResponse, _ := json.Marshal(response)
			return finalResponse, nil
		}

	}

	return nil, nil
}
func CreateIssue(issue model.Issue,userContext *system.UserContext,) ([]byte, error) {
	db, _ := system.InitMysql()
	defer db.Close()
	name := issue.Title
	if name == "" {
		return nil, errors.New("title is missing")
	}
	description := issue.Description
	if description == "" {
		return nil, errors.New("description is missing")
	}

	assignedToUserId := issue.AssignedTo
	if assignedToUserId == 0 {
		return nil, errors.New("assigned to userId missing")
	}

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS" + " " + "issues" + "(id int NOT NULL AUTO_INCREMENT, Title varchar(50) NOT NULL,Description varchar(1000) NOT NULL,AssignedToUserId varchar(50) NOT NULL," +
		"Status varchar(120) NOT NULL,CreatedBy varchar(100) NOT NULL,CreatedOn varchar(100) NOT NULL, PRIMARY KEY (id));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = db.Exec("INSERT INTO" + " " + "issues" + "(Title,Description,AssignedToUserId,Status,CreatedBy,CreatedOn) VALUES(?, ?,?,?,?,?)", name, description,assignedToUserId,"active",userContext.UserId,time.Now())
	if err != nil {
		return nil, err
	}
	userDetailsToAssigned:=FetchUserDetails(issue.AssignedTo)

	_ = time.AfterFunc(12*time.Minute, func() {SendEmailAfter12MinuteCheck(userContext.EmailId,userDetailsToAssigned.Email,issue.Title,issue.Description)})

	response := make(map[string]interface{})
	response["message"] = "successfully created"
	finalResponse, _ := json.Marshal(response)
	return finalResponse, nil
}
func UpdateIssue(issue model.Issue,userContext *system.UserContext,) ([]byte, error) {
	db, _ := system.InitMysql()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE" + " " + "issues" + " " + "set Title=?,Description=?,AssignedToUserId=?,Status=?,CreatedBy=?,CreatedOn=? WHERE (id,CreatedBy)=(?,?)")
	if err!=nil{
		return nil,err
	}

	_, err = stmt.Exec(issue.Title,issue.Description,issue.AssignedTo,"active",userContext.UserId,time.Now().UTC(), issue.ID,userContext.UserId)
	if err!=nil{
		return nil,err
	}
	userDetailsToAssigned:=FetchUserDetails(issue.AssignedTo)

	_ = time.AfterFunc(12*time.Minute, func() {SendEmailAfter12MinuteCheck(userContext.EmailId,userDetailsToAssigned.Email,issue.Title,issue.Description)})

	response := make(map[string]interface{})
	response["message"] = "successfully updated"
	finalResponse, _ := json.Marshal(response)
	return finalResponse, nil
}
func DeleteIssue(issueId int,userContext *system.UserContext) ([]byte, error) {
	db, _ := system.InitMysql()
	defer db.Close()
	var status string
	err := db.QueryRow("SELECT status FROM" + " " + "issues" + " " + "WHERE id=?",issueId).Scan(&status)
	if err != nil {
		return nil, errors.New("already ")
	}

	if status == "deleted" {
		return nil, errors.New("already deleted")
	} else {

		stmt, err := db.Prepare("UPDATE" + " " + "issues" + " " + "set status=? WHERE (id,CreatedBy)=(?,?)")
		_, err = stmt.Exec("deleted", issueId,userContext.UserId)
		if err != nil {
			fmt.Println(err)
		}
		response := make(map[string]interface{})

		response["message"] = "Issue Deleted successful"
		finalResponse, _ := json.Marshal(response)
		return finalResponse, nil
	}
}
func ListIssue(userContext *system.UserContext) ([]byte, error) {
	db, _ := system.InitMysql()
	defer db.Close()
	var userData model.Issue
	var userAllData []model.Issue

	rows, err := db.Query("select * from" + " " + "issues" + " " + "where createdBy = ?", userContext.UserId)
	if err != nil {
		fmt.Println("err ",err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&userData.ID, &userData.Title, &userData.Description, &userData.AssignedTo, &userData.Status,&userData.CreatedBy, &userData.CreatedOn)
		if err != nil {
			fmt.Println("err ",err)
		}
		userAllData = append(userAllData, userData)

	}
	response := make(map[string]interface{})
	response["issueList"] = userAllData
	finalResponse, _ := json.Marshal(response)
	return finalResponse, nil
}

func SendEmailAfter12MinuteCheck(senderMail,reciverMail,issueTitle,issueDescription string){
	auth := smtp.PlainAuth("", "user@example.com", "password", "mail.example.com")

	fmt.Println("issue Title ",issueTitle)
	fmt.Println("issue Description ",issueDescription)

	to := []string{reciverMail}
	msg := []byte("To: "+reciverMail+"\r\n" +
		"Subject: Issue Assigned!\r\n" +
		"\r\n" +
		"This is the issue body.\r\n")
	err := smtp.SendMail("mail.example.com:25", auth, senderMail, to, msg)
	if err != nil {
		fmt.Println("err ",err)
	}
	fmt.Println("email sent after 12 minute")
}
func FetchUserDetails(id int)model.User{
	db, _ := system.InitMysql()
	defer db.Close()
	var userdetails model.User
	var (
		name string
		email string
		password string
		createdOn string
		userName string
	)
	rows, err := db.Query("select * from" + " " + "user" + " " + "where id = ?", id)
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
	userdetails.ID = id
	userdetails.Name = name
	userdetails.Email = email
	userdetails.Password = password
	userdetails.UserName = userName
	userdetails.CreatedOn = createdOn
	err = rows.Err()
	if err != nil {
		fmt.Println("err ",err)
	}

	return userdetails
}
func CronProcessToIntimateUserEvery24Hours()  {

}