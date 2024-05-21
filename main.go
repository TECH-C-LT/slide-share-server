package main

import (
	"fmt"
	"os"
	"slide-share/infrastructure/adapter"
	"slide-share/infrastructure/repository/firebase"
	rstorage "slide-share/infrastructure/repository/storage"
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
	firestore, storage := lib.InitFirebase()
	driveService, _ := lib.InitDrive()
	slidesService, _ := lib.InitSlides()
	defer firestore.Close()

	userRepository := firebase.NewUserRepository(firestore)
	speakerRepository := firebase.NewSpeakerRepository(firestore)
	slideRepository := firebase.NewSlideRepository(firestore)
	thumbnailRepository := rstorage.NewThumbnailRepository(storage)

	driveAdapter := adapter.NewDriveAdapter(driveService)
	slidesAdapter := adapter.NewSlidesAdapter(slidesService)

	authUsecase := auth_usecase.NewAuthUsecase(userRepository)
	userUsecase := user_usecase.NewUserUsecase(userRepository)
	speakerUsecase := speaker_usecase.NewSpeakerUsecase(speakerRepository)
	slideUsecase := slide_usecase.NewSlideUsecase(slideRepository, driveAdapter, slidesAdapter, thumbnailRepository)

	authController := auth_http.NewAuthController(authUsecase)
	userController := user_http.NewUserController(userUsecase)
	speakerController := speaker_http.NewSpeakerController(speakerUsecase)
	slideController := slide_http.NewSlideController(slideUsecase)

	e := echo.New()
	e.Use(middleware.Logger())

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Cloud Run が自動的に設定する PORT 環境変数のデフォルト値
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
