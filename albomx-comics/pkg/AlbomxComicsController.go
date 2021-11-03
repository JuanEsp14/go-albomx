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
	server.GET("/hello", e.AlbomxcomicsHelloHandler) //each GET /hello request will be handled by AlbomxcomicsHelloHandler function
}


// AlbomxcomicsHelloHandler godoc
// @Tags Albomxcomics
// @Summary Example
// @Description Example
// @Accept  json
// @Produce  json
// @Router /hello [get]
// @Success 200 {object} dto.HelloWorldResponse
// @Failure 400 {object} string
// @Param name query string false "To who am i saying hi?"
// @x-amazon-apigateway-integration { "uri": "${lambda_arn}", "passthroughBehavior": "when_no_match", "httpMethod": "POST", "type": "aws_proxy" }
func (e *AlbomxComicsController) AlbomxcomicsHelloHandler(context *gin.Context) {
	request := new(dto.HelloWorldRequest)

	err := context.ShouldBindQuery(request) //Unmarshall query params to HelloWorldRequest struct
	if err != nil {
		_ = context.Error(fmt.Errorf("")) //Error will be handled by the ErrorHandlerMiddleware
		return
	}

	response, err := e.AlbomxComicsService.Hello(request)
	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusOK, response)
}
