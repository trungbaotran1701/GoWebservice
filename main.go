package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Contributor struct {
	OwnerName string `json:"ownerName`
	RepoName  string `json:"repoName`
}

type User struct {
	Login string `json:"login"`
}

type UserName struct {
	UserName string `json:"userName"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

func main() {
	r := gin.Default()

	// Set up CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.POST("/contributors", getAllUser)

	r.Run(":3001")
}

func getAllUser(ctx *gin.Context) {
	var contributor Contributor

	if err := ctx.BindJSON(&contributor); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error(), Code: http.StatusBadRequest})
		return
	}

	if contributor.OwnerName == "" || contributor.RepoName == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "OwnerName and RepoName are required fields", Code: http.StatusBadRequest})
		return
	}

	// Call to github API
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", contributor.OwnerName, contributor.RepoName)

	res, err := http.Get(url)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	defer res.Body.Close()

	//handle data return from github
	var users []User
	err = json.NewDecoder(res.Body).Decode(&users)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	var logins []UserName
	for _, user := range users {
		logins = append(logins, UserName{UserName: user.Login})
	}

	ctx.JSON(http.StatusOK, SuccessResponse{Data: logins, Code: http.StatusOK})
}
