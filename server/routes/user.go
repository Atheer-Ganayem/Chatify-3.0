package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/models"
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func register(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	exists, err := models.UserExists(bson.M{"email": user.Email})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, please try again later."})
		return
	}
	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email already exists, please choose another one."})
		return
	}

	filePath, code, err := utils.ExtractFileAndUpload(ctx.Request, "avatar")
	if err != nil {
		ctx.JSON(code, gin.H{"message": err.Error()})
		return
	}

	user.Avatar = filePath
	err = user.Save()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't create user."})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created Successfully."})
}

func login(ctx *gin.Context) {
	var body models.LoginBody
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := models.FindUser(bson.M{"email": body.Email})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password."})
		return
	}

	// MAYBE GENERATE A JWT TOKEN

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged in successfully.", "user": gin.H{
		"id":     user.ID,
		"email":  user.Email,
		"name":   user.Name,
		"avatar": user.Avatar,
	}})
}

func searchUsers(ctx *gin.Context) {
	searchTerm := ctx.Query("search")
	if searchTerm == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "\"search\" query parameter is required."})
		return
	}

	users, err := models.SearchUsers(searchTerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong, please try again later."})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Users fetched successfully.", "users": users})
}

func changeUserName(ctx *gin.Context) {
	type reqBody struct {
		Name string `json:"name" binding:"required,min=3,max=30"`
	}
	var body reqBody
	err := ctx.ShouldBindJSON(&body)
	body.Name = strings.TrimSpace(body.Name)
	if err != nil || len(body.Name) < 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Name must be at least 3 letters and 30 at most."})
		return
	}

	userID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User doesn't exist."})
		return
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: body.Name}}}}
	err = models.UpdateUser(userID, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't update your username."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Username updated successfully."})
}

func changePassword(ctx *gin.Context) {
	type reqBody struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword" binding:"required,min=6"`
	}
	var body reqBody
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "New password should be at least of length 6."})
		return
	}

	userID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User doesn't exist."})
		return
	}

	user, err := models.FindUser(bson.M{"_id": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try agaon later."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.CurrentPassword))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid current password."})
		return
	}

	newHashedPW, err := utils.HashPassword(body.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try agaon later."})
		return
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: newHashedPW}}}}
	err = models.UpdateUser(userID, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try agaon later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Your password has been changed successfully."})
}

func changeAvatar(ctx *gin.Context) {
	filePath, code, err := utils.ExtractFileAndUpload(ctx.Request, "file")
	if err != nil {
		ctx.JSON(code, gin.H{"message": err.Error()})
		return
	}

	userID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User doesn't exist."})
		return
	}

	user, err := models.FindUser(bson.M{"_id": userID})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User doesn't exist."})
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "avatar", Value: filePath}}}}
	err = models.UpdateUser(userID, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try agaon later."})
		return
	}

	go utils.DeleteFile(user.Avatar)

	ctx.JSON(http.StatusOK, gin.H{"message": "Your avatar has been changed successfully.", "avatar": filePath})
}
