package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

//User user model
type User struct {
	UserID           bson.ObjectId `json:"userId" bson:"_id,omitempty"`
	FullName         string        `json:"name" bson:"name"`
	Gender           string        `json:"gender" bson:"gender"`
	Mobile           string        `json:"mobile" bson:"mobile"`
	Email            string        `json:"email" bson:"email"`
	Password         string        `json:"password" bson:"password"`
	DOB              time.Time     `json:"dob" bson:"dob"`
	CreateDate       time.Time     `json:"createDate" bson:"createDate"`
	LastModifiedDate time.Time     `json:"lastModifiedDate" bson:"lastModifiedDate"`
	School           School        `json:"school" bson:"school"`
	PrevSchools      School        `json:"prevSchools" bson:"prevSchools"`
	Roles            []Role        `json:"roles" bson:"roles"`
}
