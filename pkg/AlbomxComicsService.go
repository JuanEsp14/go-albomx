package pkg

import "github.com/JuanEsp14/go-albomx/albomx-comics/pkg/dto"

type AlbomxComicsService interface {
	Hello(request *dto.HelloWorldRequest) (*dto.HelloWorldResponse, error)
}
