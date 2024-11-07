package controllers

import (
	"gotutorial/db_client"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ID        *int16    `json:"id"`
}

func CreatePost(c *gin.Context) {
	var reqBody Post

	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	response, err := db_client.DBClient.Exec(
		"INSERT INTO posts (title, content) VALUES($1, $2) RETURNING id;",
		reqBody.Title,
		reqBody.Content,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"body":       reqBody,
		"DBResponse": response,
		"error":      false,
	})
}

func GetPosts(c *gin.Context) {
	var posts []Post

	// NOTE: query is used for fetching multiple rows, queryRow is just one
	// context is a database transaction, row as well
	rows, err := db_client.DBClient.Query("SELECT id, title, content, created_at FROM posts ORDER BY id;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	for rows.Next() {
		var singlePost Post

		err = rows.Scan(&singlePost.ID, &singlePost.Title, &singlePost.Content, &singlePost.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		posts = append(posts, singlePost)
	}
	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	idSrt := c.Param("id")
	id, err := strconv.Atoi(idSrt)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	row := db_client.DBClient.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = $1;", id)
	var post Post
	err = row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "post does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, post)
}
