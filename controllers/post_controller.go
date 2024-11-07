package controllers

import (
	"gotutorial/db_client"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID        *int16    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
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

	// Use NamedExec to directly bind struct fields to named parameters
	query := `INSERT INTO posts (title, content) VALUES (:title, :content) RETURNING id`
	rows, err := db_client.DBClient.NamedQuery(query, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Retrieve the last inserted ID
	var postID int64
	if rows.Next() {
		rows.Scan(&postID)
	}
	rows.Close()

	c.JSON(http.StatusCreated, gin.H{
		"body":    reqBody,
		"post_id": postID,
		"error":   false,
	})
}

func GetPosts(c *gin.Context) {
	var posts []Post

	err := db_client.DBClient.Select(&posts, "SELECT id, title, content, created_at FROM posts ORDER BY id;")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
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

	var post Post
	err = db_client.DBClient.Get(
		&post,
		"SELECT id, title, content, created_at FROM posts WHERE id = $1;", id,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "not found",
		})
		return
	}
	c.JSON(http.StatusOK, post)
}
