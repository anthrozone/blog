package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type (
	Blog struct {
		ID        bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
		Name      string          `bson:"name" json:"name"`
		Tagline   string          `bson:"tagline" json:"tagline"`
		Owner     bson.ObjectId   `bson:"owner" json:"owner"`
		CNAME     string          `bson:"cname,omitempty" json:"cname,omitempty"`
		Editors   []bson.ObjectId `bson:"editors,omitempty" json:"editors,omitempty"`
		CreatedAt time.Time       `bson:"timestamp" json:"timestamp"`
	}
	Post struct {
		ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
		Blog      bson.ObjectId `bson:"blog" json:"blog"`
		Title     string        `bson:"title" json:"title"`
		Author    bson.ObjectId `bson:"author" json:"author"`
		CreatedAt time.Time     `bson:"timestamp" json:"timestamp,omitempty"`
		EditedAt  time.Time     `bson:"last_edited,omitempty" json:"last_edited,omitempty"`
		Content   string        `bson:"content" json:"content"`
		Tags      []string      `bson:"tags" json:"tags"`
	}
	User struct {
		ID        bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
		Username  string          `bson:"username" json:"username"`
		Password  string          `bson:"password" json:"password,omitempty"`
		FirstName string          `bson:"firstname" json:"firstname"`
		LastName  string          `bson:"lastname" json:"lastname"`
		Email     string          `bson:"email" json:"email,omitempty"`
		Blogs     []bson.ObjectId `bson:"blogs,omitempty" json:"blogs,omitempty"`
		Token     string          `bson:"token,omitempty" json:"token,omitempty"`
		Picture   string          `bson:"profile_picture" json:"profile_picture,omitempty"`
	}
	Comment struct {
		ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
		Content   string        `bson:"content" json:"content"`
		Author    bson.ObjectId `bson:"author" json:"author"`
		Post      bson.ObjectId `bson:"post" json:"post"`
		CreatedAt time.Time     `bson:"timestamp" json:"timestamp"`
	}
	Resp struct {
		Code   int         `json:"code"`
		Result interface{} `json:"result"`
	}
)
