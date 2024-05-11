package controllers

import (
	m "blog/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func ChangeNotif(c *gin.Context) {
	token := c.GetHeader("authorization")

	if token == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})

		return
	}
	userId, err := parseToken(token)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		fmt.Println(err)
		return
	}
	status := c.Query("status")
	db, err := createDb()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error opening db"})
		fmt.Println(err)
		return
	}
	if status == "0" {
		//set notif false
		query := "update users set has_notif =false where id=?"
		_, err := db.Exec(query, userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error setting notif"})
			fmt.Println(err)
			db.Close()
			return
		}

	} else if status == "1" {
		//set notif true
		query := "update users set has_notif =true where id=?"
		_, err := db.Exec(query, userId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error setting notif"})
			fmt.Println(err)
			db.Close()
			return
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		fmt.Println(err)
		db.Close()
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"success": "changed"})
	db.Close()

}
func GetPPImage(c *gin.Context) {
	token := c.GetHeader("authorization")
	if token == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})

		return
	}
	userId, err := parseToken(token)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		fmt.Println(err)
		return
	}
	db, err := createDb()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error opening db"})
		fmt.Println(err)
		return
	}
	query := "select picture_name from pp where user_id=?"
	var rowCount int = 0
	rows, err := db.Query(query, userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error geting data on db"})
		fmt.Println(err)
		return
	}
	var imagePath string

	if rows.Next() {
		err := rows.Scan(&imagePath)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error geting data from rows"})
			fmt.Println("error on getting pp")
			fmt.Println(err)
			return
		}

		rowCount++
	}

	if rowCount == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"success": "temp.png"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"success": imagePath})
	}
	db.Close()
}
func CheckNotif(c *gin.Context) {
	token := c.GetHeader("authorization")
	if token == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})

		return
	}
	userId, err := parseToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		fmt.Println(err)
		return
	}
	db, err := createDb()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error opening db"})
		fmt.Println(err)
		return
	}
	query := "select has_notif from users where id=?"
	rows, err := db.Query(query, userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error geting rows"})
		fmt.Println(err)
		return
	}
	var hastNotif bool
	if rows.Next() {
		err := rows.Scan(&hastNotif)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error scanning rows"})
			fmt.Println(err)
			return
		}
	}

	c.IndentedJSON(http.StatusOK, map[string]bool{"success": hastNotif})
	db.Close()

}
func GetUser(c *gin.Context) {
	var params m.GetUserParams
	err1 := c.BindJSON(&params)
	validator := validator.New()

	err2 := validator.Struct(params)

	if err1 != nil || err2 != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "some parameters are invalid"})
		fmt.Println(err1, err2)
		return
	}
	db, err := createDb()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error opening db"})
		fmt.Println(err)
		return
	}

	query := "select id ,password from users where username =? "
	rows, err := db.Query(query, params.Username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error getting user"})
		fmt.Println(err)
		return
	}
	var (
		id             int
		hashedPassword string
	)
	if rows.Next() {
		err = rows.Scan(&id, &hashedPassword)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error parsing user"})
			fmt.Println(err)
			return
		}
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no such user"})

		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(params.Password))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "wrong password"})

		return
	}
	token, err := generateToken(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on tokenize"})
		fmt.Println(err)
		return
	}
	fmt.Println("user id when getting user is ", id)
	c.IndentedJSON(http.StatusOK, map[string]string{"success": token})
	db.Close()
}
func CreateUser(c *gin.Context) {

	var user m.User
	err1 := c.BindJSON(&user)
	validator := validator.New()
	err2 := validator.Struct(user)
	if err1 != nil || err2 != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		fmt.Println(err1, err2)
		return
	}

	db, err := createDb()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on init database"})
		return
	}
	query := "insert into users (name,surname,username,email,password) values(?,?,?,?,?)"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on generating hashed password"})
		return
	}

	result, err := db.Exec(query, user.Name, user.Surname, user.Username, user.Email, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user exist"})
			return
		}
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "something happened when adding new user"})

		fmt.Println(err)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "something happened when adding new user on result"})
		return
	}
	fmt.Println("rows effected ", rows)
	c.IndentedJSON(http.StatusOK, gin.H{"success": "user created succesfully"})
	db.Close()
}
