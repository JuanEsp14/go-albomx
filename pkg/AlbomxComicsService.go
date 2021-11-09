package pkg

import "github.com/JuanEsp14/go-albomx/albomx-comics/pkg/dto"

type AlbomxComicsService interface {
	GetCharacters(request *dto.ComicRequest) (*dto.CharactersResponse, error)
	GetCollaborators(request *dto.ComicRequest) (*dto.CollaboratorsResponse, error)
	RefreshDataBase()
}
