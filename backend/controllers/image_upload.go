package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadImage(c *gin.Context) {
	token := c.GetHeader("authorization")
	if token == "" {
		c.String(http.StatusBadRequest, "invalid token")
		return
	}
	userId, err := parseToken(token)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, "invalid token")
		return
	}
	fmt.Println(userId)

	file, err := c.FormFile("image")

	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "No image uploaded"})
		return
	}
	if !(strings.HasSuffix(file.Filename, ".jpg") ||
		strings.HasSuffix(file.Filename, ".jpeg") ||
		strings.HasSuffix(file.Filename, ".png")) {
		c.IndentedJSON(400, gin.H{"error": "invalid image"})
		fmt.Println("the file was not valid")
		return
	}
	db, err := createDb()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "cannot open db"})
		fmt.Println(err)
		return
	}
	imageId := uuid.New()
	imageName := imageId.String() + ".png"

	//check user exist
	query1 := "select count(*) from pp where user_id=?"
	var userCount int
	rows, err := db.Query(query1, userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "cannot open db"})
		fmt.Println(err)
		return
	}
	if rows.Next() {
		err := rows.Scan(&userCount)
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "cannot open db"})
			return
		}
	}
	//--------------------------------------------
	if userCount == 0 {
		query := "insert into pp values(?,?)"
		_, err = db.Exec(query, userId, imageName)
		if err != nil {
			fmt.Println("burda")
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error adding image to db"})
			return
		}
		fmt.Println("new image added")
	} else {
		query := "update pp set picture_name=? where user_id=?"
		_, err = db.Exec(query, imageName, userId)
		if err != nil {
			fmt.Println("burda1")
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error adding image to db"})
			return
		}
		fmt.Println("current user  added new image")
	}
	err = c.SaveUploadedFile(file, "../pp_images/"+imageName)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"success": "image added succesfully"})
	db.Close()
}
