package main

import (
	"github.com/globalsign/mgo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/anthrozone/blog/platform"
	"log"
	"os"
)

func main() {
	e := echo.New()
	blog := new(platform.Platform)

	blog.Key = "mysuperawesometestkey"

	var err error

	blog.Mongo, err = mgo.Dial(os.Getenv("DB_HOST"))
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(blog.Key),
	}))

	e.GET("/blog/:blog/posts", blog.GetPosts)
	e.GET("/blog/:blog/posts/:id", blog.GetPost)
	e.POST("/blog/:blog/posts", blog.CreatePost)
	e.DELETE("/blog/:blog/posts/:id", blog.DeletePost)
	e.PATCH("/blog/:blog/posts", blog.UpdatePost)
	e.GET("/blog", blog.GetBlogs)
	e.GET("/blog/:blog", blog.GetBlog)
	e.POST("/blog", blog.CreateBlog)
	e.DELETE("/blog/:blog", blog.DeleteBlog)
	e.PATCH("/blog", blog.UpdateBlog)
	e.GET("/users/:id", blog.GetUser)
	e.PATCH("/users", blog.UpdateUser)
	e.GET("/blog/:blog/posts/:post/comments", blog.GetComments)
	e.POST("/blog/:blog/posts/:post/comments", blog.CreateComment)
	e.PATCH("/blog/:blog/posts/:post/comments", blog.UpdateComment)
	e.DELETE("/blog/:blog/posts/:post/comments/:id", blog.DeleteBlog)

	e.Logger.Fatal(e.Start(":8080"))
}
