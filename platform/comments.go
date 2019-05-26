package platform

import (
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"gitlab.com/anthrozone/blog/models"
	"net/http"
	"time"
)

func (b *Platform) GetComments(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("comments")

	var comments []models.Comment
	err = collection.Find(bson.M{"post": bson.ObjectIdHex(c.Param("post"))}).All(&comments)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "No post with that ID"})
	}

	return c.JSON(http.StatusOK, models.Resp{Code: http.StatusOK, Result: comments})
}

func (b *Platform) CreateComment(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("comments")

	comment := models.Comment{
		ID: bson.NewObjectId(),
	}

	if err = c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Bad JSON"})
	}

	user := b.userFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
	}

	comment.CreatedAt = time.Now()
	comment.Author = user.ID

	err = collection.Insert(comment)
	if err != nil {
		return c.JSON(http.StatusConflict, models.Resp{Code: http.StatusConflict, Result: "Comment already exists"})
	}

	return c.JSON(http.StatusOK, comment)
}

func (b *Platform) UpdateComment(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("comments")

	var comment models.Comment

	if err = c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Bad JSON"})
	}

	user := b.userFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
	}

	var oldComment models.Comment
	err = collection.FindId(comment.ID).One(&oldComment)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Post not found"})
	}

	if oldComment.Author != user.ID {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Only the author may update a comment"})
	}

	if comment.Content != "" {
		oldComment.Content = comment.Content
	}

	err = collection.UpdateId(oldComment.ID, oldComment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Resp{Code: http.StatusInternalServerError, Result: "Oops"})
	}

	return c.JSON(http.StatusOK, models.Resp{Code: http.StatusOK, Result: oldComment})
}

func (b *Platform) DeleteComment(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("comments")

	user := b.userFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
	}

	var comment models.Comment
	err = collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&comment)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Comment not found"})
	}

	if comment.Author == user.ID {
		err = collection.RemoveId(comment.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "Comment not found"})
		}
	} else {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Only the comment author may delete"})
	}

	return c.JSON(http.StatusOK, models.Resp{Code: http.StatusOK, Result: "Comment deleted"})
}
