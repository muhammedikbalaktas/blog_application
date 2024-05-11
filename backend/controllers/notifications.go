package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	token := c.GetHeader("authorization")
	var messages = make([]string, 0)
	if token == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid header"})
		return
	}
	userId, err := parseToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "invalid token"})
		return
	}
	db, err := createDb()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "error opennig db"})
		return
	}

	// 	select username, title, act_time
	// from (select * from users
	//       inner join notification on users.id = notification.like_source) as res
	// inner join blogs on res.blog_id = blogs.id where res.like_dest=5 order by act_time desc;
	query := "select username, title , act_time from " +
		"(select* from users inner join notification on users.id=notification.like_source) as res " +
		"inner join blogs on res.blog_id=blogs.id " +
		"where res.like_dest=? order by act_time desc"
	rows, err := db.Query(query, userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on executing query"})
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var (
			username string
			title    string
			act_time string
		)
		err := rows.Scan(&username, &title, &act_time)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error getting rows"})
			fmt.Println(err)
			return
		}
		message := fmt.Sprintf("%s liked your %s blog at time %s", username, title, act_time)
		messages = append(messages, message)
	}
	c.IndentedJSON(http.StatusOK, messages)
	db.Close()
}
