package platform

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"gitlab.com/anthrozone/blog/models"
	"io/ioutil"
	"net/http"
	"strings"
)

func (b *Platform) GetUser(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("users")

	var userID bson.ObjectId

	if c.Param("id") == "@me" {
		user := b.userFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
		}
		userID = user.ID
	} else {
		userID = bson.ObjectIdHex(c.Param("id"))
	}

	var user models.User
	err = collection.FindId(userID).One(&user)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.Resp{Code: http.StatusNotFound, Result: "User not found"})
	}

	if !(c.Param("id") == "@me") {
		user.Email = ""
		user.Token = ""
	}

	// Never send the password hash
	user.Password = ""

	return c.JSON(http.StatusOK, models.Resp{Code: http.StatusOK, Result: user})
}

func (b *Platform) UpdateUser(c echo.Context) (err error) {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("users")

	authUser := b.userFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.Resp{Code: http.StatusUnauthorized, Result: "Bad user"})
	}

	var user models.User

	if err = c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Bad JSON"})
	}

	if user.Email != "" {
		authUser.Email = user.Email
	}

	if user.Password != "" {
		authUser.Password, _ = hashPassword(user.Password)
	}

	if user.FirstName != "" {
		authUser.FirstName = user.FirstName
	}

	if user.LastName != "" {
		authUser.LastName = user.LastName
	}

	if user.Picture != "" {
		imageData, err := base64.StdEncoding.DecodeString(user.Picture)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Image must be base64 encoded"})
		}
		fileType := http.DetectContentType(imageData)
		if !strings.HasPrefix(fileType, "image") {
			return c.JSON(http.StatusBadRequest, models.Resp{Code: http.StatusBadRequest, Result: "Profile picture must be an image"})
		}
		// Generate filename
		h := sha256.New()
		h.Write(imageData)
		fileHash := hex.EncodeToString(h.Sum(nil))
		// TODO implement image saving to cloud storage
		err = ioutil.WriteFile(fileHash, imageData, 0644)
		if err != nil {
			return c.JSON(http.StatusInsufficientStorage, models.Resp{Code: http.StatusInsufficientStorage, Result: fmt.Sprintf("Well, this is embarrassing: %s", err)})
		}
		authUser.Picture = fileHash
	}

	err = collection.UpdateId(authUser.ID, authUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Resp{Code: http.StatusInternalServerError, Result: "An unknown error occurred"})
	}

	authUser.Password = ""

	return c.JSON(http.StatusOK, models.Resp{Code: http.StatusOK, Result: authUser})
}
