package pkg

import (
	"fmt"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AlbomxComicsController struct {
	AlbomxComicsService AlbomxComicsService
	logger 				logrus.Logger
}


//NewAlbomxComicsController Constructor
func NewAlbomxComicsController(AlbomxComicsService AlbomxComicsService, logger logrus.Logger) *AlbomxComicsController {
	return &AlbomxComicsController {
			AlbomxComicsService: AlbomxComicsService,
			logger: logger,
	}
}


func (e *AlbomxComicsController) SetupRouter(server *gin.Engine)  {
	server.GET("/characters", e.AlbomxcomicsCharactersHandler)
	server.GET("/collaborators", e.AlbomxcomicsCollaboratorsHandler)
}


// AlbomxcomicsCharactersHandler godoc
// @Tags Albomxcomics
// @Summary get
// @Description Get characters info
// @Accept  json
// @Produce  json
// @Router /characters [get]
// @Success 200 {object} dto.CharactersResponse
// @Failure 400 {object} string
// @Param avengerId query string false "To who am i getting?"
func (e *AlbomxComicsController) AlbomxcomicsCharactersHandler(context *gin.Context) {
	request := new(dto.ComicRequest)

	err := context.ShouldBindQuery(request) //Unmarshall query params to HelloWorldRequest struct
	if err != nil {
		_ = context.Error(fmt.Errorf("")) //Error will be handled by the ErrorHandlerMiddleware
		return
	}

	response, err := e.AlbomxComicsService.GetCharacters(request)
	if err != nil {
		_ = context.Error(err)
		return
	}
	if response.LastSync == ""{
		context.JSON(http.StatusNotFound, fmt.Sprintf("Id %s not found in database", request.AvengerId))
		return
	}

	context.JSON(http.StatusOK, response)
}

// AlbomxcomicsCollaboratorsHandler godoc
// @Tags Albomxcomics
// @Summary get
// @Description Get collaborators info
// @Accept  json
// @Produce  json
// @Router /collaborators [get]
// @Success 200 {object} dto.CollaboratorsResponse
// @Failure 400 {object} string
// @Param avengerId query string false "To who am i getting?"
func (e *AlbomxComicsController) AlbomxcomicsCollaboratorsHandler(context *gin.Context) {
	request := new(dto.ComicRequest)

	err := context.ShouldBindQuery(request) //Unmarshall query params to HelloWorldRequest struct
	if err != nil {
		_ = context.Error(fmt.Errorf("")) //Error will be handled by the ErrorHandlerMiddleware
		return
	}

	response, err := e.AlbomxComicsService.GetCollaborators(request)
	if err != nil {
		_ = context.Error(err)
		return
	}
	if response.LastSync == ""{
		context.JSON(http.StatusNotFound, fmt.Sprintf("Id %s not found in database", request.AvengerId))
		return
	}

	context.JSON(http.StatusOK, response)
}

