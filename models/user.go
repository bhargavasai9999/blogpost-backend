package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleUser    = "user"
	RoleCreator = "creator"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"password" validate:"required"`
	FirstName string             `bson:"firstname" json:"firstname" validate:"required"`
	LastName  string             `bson:"lastname" json:"lastname" validate:"required"`
	Role      string             `bson:"role" json:"role" validate:"required"`
	Followers []string           `bson:"followers,omitempty" `
	Following []string           `bson:"following,omitempty" `
}

func NewUser(username, email, password, firstName, lastName, role string) (*User, error) {
	if role != RoleCreator && role != RoleUser {
		return nil, errors.New("Invalid Role specified " + role)
	}

	user := &User{
		Id:        primitive.NewObjectID(),
		Username:  username,
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		Following: []string{},
	}
	if role == RoleCreator {
		user.Followers = []string{}

	}
	return user, nil

}
