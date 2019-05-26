package platform

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"gitlab.com/anthrozone/blog/models"
	"net/http"
	"time"
)

func isEditor(u models.User, b models.Blog) bool {
	for i := 0; i < len(b.Editors); i++ {
		if u.ID == b.Editors[i] {
			return true
		}
	}
	return false
}

func (b *Platform) GetPosts(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("posts")
	var posts []models.Post
	err = collection.Find(bson.M{"blog": bson.ObjectIdHex(c.Param("blog"))}).All(&posts)
	if err != nil {
		return c.String(http.StatusNotFound, "Blog not found")
	}

	return c.JSON(http.StatusOK, posts)
}

func (b *Platform) GetPost(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("posts")
	var post models.Post
	err = collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&post)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Unable to find post: %s", err))
	}

	if post.Blog != bson.ObjectIdHex(c.Param("blog")) {
		return c.String(http.StatusNotFound, "Unable to find post")
	}

	return c.JSON(http.StatusOK, post)
}

func (b *Platform) CreatePost(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("posts")
	blogCollection := sessionCopy.DB("blog").C("blogs")

	post := models.Post{
		ID: bson.NewObjectId(),
	}

	if err = c.Bind(&post); err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Unable to parse JSON: %s", err))
	}

	user := b.userFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
	}

	var blog models.Blog
	err = blogCollection.FindId(bson.ObjectIdHex(c.Param("blog"))).One(&blog)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Blog not found"})
	}

	if blog.Owner != user.ID && !isEditor(user, blog) {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "You do not have permission to post to this blog"})
	}

	post.CreatedAt = time.Now()
	post.Blog = blog.ID
	post.Author = user.ID

	err = collection.Insert(post)
	if err != nil {
		return c.String(http.StatusConflict, "Post already exists with that ID")
	}
	return c.JSON(http.StatusOK, post)
}

func (b *Platform) DeletePost(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("posts")
	blogCollection := sessionCopy.DB("blog").C("blogs")

	user := b.userFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
	}

	var blog models.Blog
	err = blogCollection.FindId(bson.ObjectIdHex(c.Param("blog"))).One(&blog)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Blog not found"})
	}

	if blog.Owner != user.ID && !isEditor(user, blog) {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "You do not have permission to post to this blog"})
	}

	var post models.Post
	err = collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&post)
	if err != nil || post.Blog != blog.ID {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Post not found"})
	}

	err = collection.RemoveId(bson.ObjectIdHex(c.Param("id")))
	if err != nil {
		return c.String(http.StatusNotFound, "Post not found")
	}

	return c.String(http.StatusOK, "Post deleted successfully")
}

func (b *Platform) UpdatePost(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("posts")
	blogCollection := sessionCopy.DB("blog").C("blogs")

	var post models.Post
	if err = c.Bind(&post); err != nil {
		return c.String(http.StatusBadRequest, "Unable to parse JSON")
	}

	user := b.userFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
	}

	var blog models.Blog
	err = blogCollection.FindId(bson.ObjectIdHex(c.Param("blog"))).One(&blog)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Blog not found"})
	}

	if blog.Owner != user.ID && !isEditor(user, blog) {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "You do not have permission to post to this blog"})
	}

	var oldPost models.Post
	err = collection.FindId(post.ID).One(&oldPost)
	if err != nil || oldPost.Blog != blog.ID {
		return c.String(http.StatusNotFound, "Post not found")
	}

	// Don't allow users to change the createdAt time
	post.CreatedAt = oldPost.CreatedAt
	post.Blog = oldPost.Blog

	post.EditedAt = time.Now()

	_, err = collection.Upsert(bson.M{"_id": post.ID}, post)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)
}
