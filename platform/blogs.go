package platform

import (
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"gitlab.com/anthrozone/blog/models"
	"net/http"
	"time"
)

func (b *Platform) GetBlogs(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("blogs")

	var blogs []models.Blog
	err = collection.Find(nil).All(&blogs)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, models.Resp{Code: http.StatusOK, Result: blogs})
}

func (b *Platform) GetBlog(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("blogs")

	var blog models.Blog
	err = collection.FindId(bson.ObjectIdHex(c.Param("blog"))).One(&blog)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Blog not found"})
	}

	return c.JSON(http.StatusOK, models.Resp{Code: http.StatusOK, Result: blog})
}

func (b *Platform) CreateBlog(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("blogs")

	blog := models.Blog{
		ID: bson.NewObjectId(),
	}

	if err = c.Bind(&blog); err != nil {
		return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Bad JSON"})
	}

	owner := b.userFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
	}

	blog.CreatedAt = time.Now()
	blog.Owner = owner.ID

	err = collection.Insert(blog)
	if err != nil {
		return c.JSON(http.StatusConflict, models.Resp{Code: http.StatusConflict, Result: "Blog already exists"})
	}

	return c.JSON(http.StatusOK, blog)
}

func (b *Platform) DeleteBlog(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("blogs")

	owner := b.userFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
	}

	var blog models.Blog
	err = collection.FindId(bson.ObjectIdHex(c.Param("blog"))).One(&blog)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Blog not found"})
	}

	if blog.Owner != owner.ID {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Only the blog owner can delete a blog"})
	}

	err = collection.RemoveId(bson.ObjectIdHex(c.Param("blog")))
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Blog not found"})
	}

	return c.JSON(http.StatusOK, models.Resp{Code: http.StatusOK, Result: "Blog deleted successfully"})
}

func (b *Platform) UpdateBlog(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("blogs")

	var blog models.Blog
	if err = c.Bind(&blog); err != nil {
		return c.JSON(http.StatusBadRequest, "Bad JSON")
	}

	var oldBlog models.Blog
	err = collection.FindId(blog.ID).One(&oldBlog)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Blog not found"})
	}

	owner := b.userFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
	}

	if oldBlog.Owner != owner.ID {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Blog can only be edited by the owner"})
	}

	blog.CreatedAt = oldBlog.CreatedAt

	err = collection.UpdateId(bson.ObjectIdHex(c.Param("blog")), blog)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, blog)
}
