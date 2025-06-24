package models

import (
	"context"
	"time"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/db"
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID        bson.ObjectID `json:"_id" bson:"_id"`
	Name      string        `json:"name" bson:"name" form:"name" binding:"required,min=3,max=30"`
	Email     string        `json:"email" bson:"email" form:"email" binding:"required,email"`
	Avatar    string        `json:"avatar" bson:"avatar"`
	Password  string        `json:"password,omitempty" bson:"password" form:"password" binding:"required,min=6"`
	CreatedAt time.Time     `json:"createdAt,omitempty" bson:"createdAt"`
}

type LoginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (user *User) Save() error {
	hashedPw, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPw
	user.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = db.Users.InsertOne(ctx, user)

	return err
}

func UserExists(filter bson.M) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var user User
	err := db.Users.FindOne(ctx, filter).Decode(&user)

	switch err {
	case nil:
		return true, nil
	case mongo.ErrNoDocuments:
		return false, nil
	default:
		return true, err
	}
}

func FindUser(filter bson.M) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var user User
	err := db.Users.FindOne(ctx, filter).Decode(&user)

	return user, err
}

func SearchUsers(term string) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	regex := bson.M{"$regex": term, "$options": "i"} // i = case-insensitive

	filter := bson.M{
		"$or": []bson.M{
			{"name": regex},
			{"email": regex},
		},
	}

	projection := bson.M{
		"name":   1,
		"email":  1,
		"avatar": 1,
	}

	opts := options.Find().SetProjection(projection)

	var users []User
	cursor, err := db.Users.Find(ctx, filter, opts)
	if err != nil {
		return users, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &users)

	return users, err
}

func UpdateUser(userID bson.ObjectID, update bson.D) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.Users.UpdateByID(ctx, userID, update)

	return err
}
