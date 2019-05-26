package platform

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"gitlab.com/anthrozone/blog/models"
	"golang.org/x/crypto/bcrypt"
)

func (b *Platform) userFromToken(c echo.Context) models.User {
	sessionCopy := b.Mongo.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB("blog").C("users")

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	var u models.User
	err := collection.FindId(bson.ObjectIdHex(claims["id"].(string))).One(&u)
	if err != nil {
		return models.User{}
	}
	return u
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
