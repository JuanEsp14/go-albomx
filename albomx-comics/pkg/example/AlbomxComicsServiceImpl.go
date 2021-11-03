package example

import (
	"fmt"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg/dto"
)

type AlbomxComicsServiceImpl struct {

}

func NewAlbomxComicsService() pkg.AlbomxComicsService {
	return &AlbomxComicsServiceImpl{}
}

func (h *AlbomxComicsServiceImpl) Hello(request *dto.HelloWorldRequest) (*dto.HelloWorldResponse, error) {
	response := &dto.HelloWorldResponse{
		Message: fmt.Sprintf("Hello %s", request.Name),
	}

	return response, nil
}

