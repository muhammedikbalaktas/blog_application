package controllers

import (
	m "blog/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func GetSingleBlog(c *gin.Context) {
	blogId, isValid := c.GetQuery("blog_id")
	if !isValid {
		fmt.Println("unsatisfied url")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "inavlid url"})
		return
	}
	likeCount, err := GetLikes(blogId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("like count is ", likeCount)

	var blog m.Blog
	blog.LikeCount = likeCount
	db, err := createDb()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error opening database"})
		fmt.Println(err)
		return
	}
	query := "select username, title, content, image_name, div_id " +
		"from users left join blogs on users.id=user_id where div_id=? "
	rows, err := db.Query(query, blogId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error getting rows"})
		fmt.Println(err)
		return
	}
	if rows.Next() {
		err := rows.Scan(&blog.Username, &blog.Title, &blog.Content, &blog.ImageName, &blog.DivId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error scanning result"})
			fmt.Println(err)
			return
		}
	}
	c.IndentedJSON(http.StatusOK, blog)
	db.Close()
}
func GetAllBlogs(c *gin.Context) {
	blogs := make([]m.Blog, 0)

	db, err := createDb()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on opening database"})
		return
	}
	query := "select username, title, content, image_name, div_id " +
		"from users left join blogs on users.id=user_id;"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on executing query"})
		return
	}
	for rows.Next() {
		var blog m.Blog
		err := rows.Scan(&blog.Username, &blog.Title, &blog.Content, &blog.ImageName, &blog.DivId)
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on parsing rows"})
			return
		}
		likeCount, err := GetLikes(blog.DivId)
		if err != nil {
			fmt.Println(err)
			return
		}

		blog.LikeCount = likeCount
		blogs = append(blogs, blog)
	}
	c.IndentedJSON(http.StatusOK, blogs)
	db.Close()

}
func CreateBlog(c *gin.Context) {
	var blog m.Blog
	err1 := c.Bind(&blog)
	validator := validator.New()
	err2 := validator.Struct(blog)
	if err1 != nil || err2 != nil {
		fmt.Println(err1, err2)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}
	db, err := createDb()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on opening database"})
		return
	}
	userId, err := parseToken(blog.UserToken)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	imageId := uuid.New()
	imageName := imageId.String() + ".png"
	divId := uuid.New()
	query := "insert into blogs (title,content,user_id,publish_date,image_name,div_id) values(?,?,?,now(),?,?)"
	_, err = db.Exec(query, blog.Title, blog.Content, userId, imageName, divId)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error on adding blog"})
		return
	}
	file, err := c.FormFile("image")
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "No image uploaded"})
		return
	}
	err = c.SaveUploadedFile(file, "../images/"+imageName)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"succes": "blog added succesfully"})
	db.Close()
}
