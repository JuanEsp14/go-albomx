package main

import (
	"database/sql"
	_ "github.com/JuanEsp14/go-albomx/albomx-comics/docs"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg/example"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gopkg.in/robfig/cron.v2"
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
	c := cron.New()
	//"*/1 * * * * *" -> un segundo
	c.AddFunc("*/1 * * * *", services.ExpireProviderInvitation)
	c.Start()

	router := gin.Default()
	serviceSpotify := example.NewAlbomxComicsService()
	controller := pkg.NewAlbomxComicsController(serviceSpotify, logrus.Logger{})
	controller.SetupRouter(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" { port = "8080"}


	query := "CREATE TABLE IF NOT EXISTS characters(tick timestamp)"
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logrus.Errorf("Error opening database: %q", err)
	}
	defer db.Close()
	res, err := db.Exec(query)
	if err != nil {
		logrus.Errorf("Error executing query")
		logrus.Error(err)
		return
	}
	logrus.Info("Connect DB successfully")
	logrus.Info(res)

	if _, err := db.Exec("INSERT INTO characters VALUES (now())"); err != nil {
		logrus.Errorf("Error executing query")
		logrus.Error(err)
		return
	}

	router.Run(":" + port)
}
