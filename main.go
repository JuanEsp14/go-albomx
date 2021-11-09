package main

import (
	_ "github.com/JuanEsp14/go-albomx/albomx-comics/docs"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg/repository"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"os"
)

// @title albomx-comics API
// @version 1.0.0
// @description albomx-comics API.

// @tag.name Albomxcomics
// @tag.description All the available albomx-comics operations

// @contact.name Juan Espinoza
// @contact.email juanmesp@hotmail.com

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	log := logrus.New()
	repository := repository.NewAlbomxComicsRepository(log)

	marvelService := services.NewMarvelService(log, new(http.Client), &repository)
	c := cron.New()
	c.AddFunc("*/10 * * * *", marvelService.RefreshDataBase)
	c.Start()

	router := gin.Default()
	controller := pkg.NewAlbomxComicsController(marvelService, logrus.Logger{})
	controller.SetupRouter(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" { port = "8080"}

	router.Run(":" + port)
}