package services

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg/dto"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg/repository"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	ironMan		= "iron man"
	capAmerica 	= "captain america"
	apiKey		= "API_KEY"
	privateKey  = "PRIVATE_KEY"
	baseUrl		= "http://gateway.marvel.com/v1/public/comics"
)

type marvelService struct {
	logger 	*logrus.Logger
	client 	dto.HTTPClient
	db 		*repository.AlbomxComicsRepository
}

func NewMarvelService(log *logrus.Logger, client dto.HTTPClient, repository *repository.AlbomxComicsRepository) marvelService {
	return marvelService{
		logger: log,
		client: client,
		db: repository,
	}
}

func (s *marvelService) RefreshDataBase() {
	s.logger.Info("Active cron")
	s.processCharacter(ironMan)
	s.processCharacter(capAmerica)
}

func (s *marvelService) processCharacter(characterName string) {
	characterComicsResponse, err := s.refreshComic(characterName)
	if err != nil {
		s.logger.Errorf("Error processing character %s", characterName, err)
	}
	if characterComicsResponse.Data.Results == nil{
		s.logger.Errorf("No data to process")
		return
	}
	s.processMarvelResponse(characterComicsResponse, characterName)
}

func (s *marvelService) refreshComic(comicName string) (*dto.MarvelComicsResponse, error) {
	s.logger.Info("Getting request")
	request, err := getRequest(comicName)
	response, err := s.client.Do(request)

	responseType := new(dto.MarvelComicsResponse)
	err = json.NewDecoder(response.Body).Decode(responseType)
	if err != nil {
		s.logger.Error("Error decoder response", err)
		return nil, err
	}
	s.logger.Infof("Chapter %s obtained", comicName)
	return responseType, nil
}

func (s *marvelService) processMarvelResponse(response *dto.MarvelComicsResponse, characterName string) {
	s.db.InitializeDataBase()
	s.db.UpdateDatabase(response, characterName)
}

func getRequest(comicName string) (*http.Request, error) {
	timeStamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	apiKey := os.Getenv(apiKey)
	privateKey := os.Getenv(privateKey)
	auth :=  timeStamp + privateKey + apiKey
	hasher := md5.New()
	hasher.Write([]byte(auth))
	hash := hex.EncodeToString(hasher.Sum(nil))

	request, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		logrus.Errorf("Error generating Request: %s", err)
		return nil, err
	}
	q := request.URL.Query()
	q.Add("title", comicName)
	q.Add("ts", timeStamp)
	q.Add("apikey", apiKey)
	q.Add("hash", hash)
	q.Add("limit", "1")
	request.URL.RawQuery = q.Encode()

	return request, nil
}