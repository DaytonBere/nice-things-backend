package main

import (
	"Nice-Things-Backend/controllers"
	"Nice-Things-Backend/initializers"
	"Nice-Things-Backend/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	route := gin.Default()
	route.Use(middleware.CORSMiddleware())

	route.POST("/api/user/signUp", middleware.RequireAuth, controllers.SignUp)
	route.POST("/api/user/signIn", controllers.SignIn)
	route.POST("/api/user/signOut", controllers.SignOut)
	route.GET("/api/user/getUsersTest", controllers.GetUsersTest)
	route.PATCH("/api/user/changePassword", middleware.RequireAuth, controllers.ChangePassword)
	route.GET("/api/user/validate", middleware.RequireAuth, controllers.Validate)

	route.GET("/api/niceThings/getUsers", middleware.RequireAuth, controllers.GetUsers)
	route.POST("/api/niceThings/createNiceThing",  middleware.RequireAuth, controllers.CreateNiceThing)
	route.PATCH("/api/niceThings/editNiceThing",  middleware.RequireAuth, controllers.EditNiceThing)
	route.GET("/api/niceThings/getUserNiceThings",  middleware.RequireAuth, controllers.GetUserNiceThings)

	route.Run()
}