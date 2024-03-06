package controllers

import (
	"crud_go_gin/initailizer"
	"crud_go_gin/models"
	"crud_go_gin/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type postBodyPayload struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

func postCreate(c *gin.Context) {
	// get data from request body
	var body postBodyPayload

	if !util.BindAndValidate(c, &body) {
		return
	}

	// create a post
	post := models.Post{Title: body.Title, Body: body.Body}

	// return it
	if createPostResult := initailizer.DB.Create(&post); createPostResult.Error != nil {
		util.ResponseWithError(c, http.StatusBadRequest, "create post failed")
		return
	}

	util.ResponseWithJson(c, http.StatusOK, post)
}

func getIndex(c *gin.Context) {
	// get the posts
	var posts []models.Post

	if queryResult := initailizer.DB.Find(&posts); queryResult.Error != nil {
		util.ResponseWithError(c, http.StatusNotFound, "did not find any posts")
		return
	}

	// return them
	util.ResponseWithJson(c, http.StatusOK, posts)
}

func getOnePost(c *gin.Context) {
	// get the post and id
	var post models.Post

	if queryResult := initailizer.DB.First(&post, c.Param("id")); queryResult.Error != nil {
		util.ResponseWithError(c, http.StatusNotFound, "Post not found")
		return
	}

	// return it
	util.ResponseWithJson(c, http.StatusOK, post)
}

func postUpdate(c *gin.Context) {
	// get the post and id
	var post models.Post

	if queryResult := initailizer.DB.First(&post, c.Param("id")); queryResult.Error != nil {
		util.ResponseWithError(c, http.StatusNotFound, "Post not found")
		return
	}

	// get data from request body
	var body postBodyPayload

	if !util.BindAndValidate(c, &body) {
		return
	}

	initailizer.DB.Model(&post).Updates(models.Post{Title: body.Title, Body: body.Body})

	// // update the post
	// post.Title = body.Title
	// post.Body = body.Body

	// // save it
	// initailizer.DB.Save(&post)

	// return it
	util.ResponseWithJson(c, http.StatusOK, post)
}

func postDelete(c *gin.Context) {
	// get the post and id
	var post models.Post

	if queryResult := initailizer.DB.First(&post, c.Param("id")); queryResult.Error != nil {
		util.ResponseWithError(c, http.StatusNotFound, "Post not found")
		return
	}

	// delete it
	if deleteResult := initailizer.DB.Delete(&post); deleteResult.Error != nil {
		util.ResponseWithError(c, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	// return it
	util.ResponseWithJson(c, http.StatusOK, post)
}

func SetupPostsRouter(r *gin.Engine) {
	postRouters := r.Group("/posts")
	{
		postRouters.POST("/", postCreate)
		postRouters.PUT("/update/:id", postUpdate)

		postRouters.GET("/", getIndex)
		postRouters.GET("/:id", getOnePost)

		postRouters.DELETE("/delete/:id", postDelete)
	}
}
