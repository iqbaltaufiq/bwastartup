package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(localhost:3306)/bwastartup?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	campaigns, _ := campaignRepository.FindByUserId(4)
	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println("debug")
	for _, campaign := range campaigns {
		fmt.Println(campaign.Name)
		fmt.Println("length: ", len(campaign.CampaignImages))
		if len(campaign.CampaignImages) > 0 {
			fmt.Println(campaign.CampaignImages[0].FileName)
		}
	}

	userService := user.NewService(userRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check for bearer token in header 'Authorization'
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Bearer not found", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// get the token from the Bearer auth
		var jwtToken string
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 {
			jwtToken = bearerToken[1]
		}

		// validate the token
		token, errVal := authService.ValidateToken(jwtToken)
		if errVal != nil {
			response := helper.APIResponse("Invalid token", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// get the user id from inside the claim/payload
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Failed parsing claim", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userId := int(claim["user_id"].(float64))

		// fetch user using id from the token
		user, errFind := userService.GetUserById(userId)
		if errFind != nil {
			response := helper.APIResponse("User not found", http.StatusNotFound, "error", nil)
			c.AbortWithStatusJSON(http.StatusNotFound, response)
			return
		}

		// pass the user from middleware to be used in the handler
		c.Set("currentUser", user)
	}
}
