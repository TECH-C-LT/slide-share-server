package main

import (
	"fmt"
	"os"
	"slide-share/infrastructure/repository/firebase"
	"slide-share/lib"
	auth_http "slide-share/service/auth/http"
	auth_usecase "slide-share/service/auth/usecase"
	slide_http "slide-share/service/slides/http"
	slide_usecase "slide-share/service/slides/usecase"
	speaker_http "slide-share/service/speaker/http"
	speaker_usecase "slide-share/service/speaker/usecase"
	user_http "slide-share/service/user/http"
	user_usecase "slide-share/service/user/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	firestore, _ := lib.InitFirebase()
	defer firestore.Close()

	userRepository := firebase.NewUserRepository(firestore)
	speakerRepository := firebase.NewSpeakerRepository(firestore)
	slideRepository := firebase.NewSlideRepository(firestore)

	authUsecase := auth_usecase.NewAuthUsecase(userRepository)
	userUsecase := user_usecase.NewUserUsecase(userRepository)
	speakerUsecase := speaker_usecase.NewSpeakerUsecase(speakerRepository)
	slideUsecase := slide_usecase.NewSlideUsecase(slideRepository)

	authController := auth_http.NewAuthController(authUsecase)
	userController := user_http.NewUserController(userUsecase)
	speakerController := speaker_http.NewSpeakerController(speakerUsecase)
	slideController := slide_http.NewSlideController(slideUsecase)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))

	auth_http.NewAuthRouter(e, authController)
	user_http.NewUserRouter(e, userController)
	speaker_http.NewSpeakerRouter(e, speakerController)
	slide_http.NewSlideRouter(e, slideController)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", 8080)))
}
