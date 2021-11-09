package repository

import (
	"database/sql"
	"fmt"
	"github.com/JuanEsp14/go-albomx/albomx-comics/pkg/dto"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

const (
	ironMan		= "iron man"
	ironManId	= "ironman"
	capAmericaId = "capamerica"
	colorist	= "colorist"
	writer 		= "writer"
	editor		= "editor"
)

//GetConnection is a singleton than define the type of connection to invoke database.
func (r *AlbomxComicsRepository) GetConnection() *sql.DB{
	//once.Do(func() {
		db, errorConnection := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if errorConnection != nil {
			r.logger.Errorf("Error opening database: %q", errorConnection)
			panic("Error to connect database")
		}
	return db
	//})
}

type AlbomxComicsRepository struct {
	logger 	*logrus.Logger
}

func NewAlbomxComicsRepository(log *logrus.Logger) AlbomxComicsRepository {
	return AlbomxComicsRepository{
		logger: log,
	}
}


func (r *AlbomxComicsRepository) InitializeDataBase() {
	r.logger.Info("Initializing data base")
	if err := r.createComicsTable(); err != nil {
		return
	}
	if err := r.createCollaboratorsTable(); err != nil {
		return
	}
	if err := r.createCharactersTable(); err != nil {
		return
	}
}

func (r *AlbomxComicsRepository) GetCollaborators(id string) (*dto.CollaboratorsResponse, error){
	r.logger.Infof("Getting collaborators to id %s", id)
	db := r.GetConnection()
	rowsComic, err := db.Query("SELECT * FROM comics WHERE characterName = $1", id)
	if err != nil {
		r.logger.Error("Error getting comic. ", err)
		return nil, err
	}
	rowsCollaborator, err := db.Query("SELECT * FROM collaborators WHERE comicId = $1", id)
	if err != nil {
		r.logger.Error("Error getting comic. ", err)
		return nil, err
	}
	db.Close()
	return r.getCollaboratorsResponse(rowsComic, rowsCollaborator)
}

func (r *AlbomxComicsRepository) GetCharacters(id string) (*dto.CharactersResponse, error){
	r.logger.Infof("Getting collaborators to id %s", id)
	db := r.GetConnection()
	rowsComic, err := db.Query("SELECT * FROM comics WHERE characterName = $1", id)
	if err != nil {
		r.logger.Error("Error getting comic. ", err)
		return nil, err
	}
	rowsCharacter, err := db.Query("SELECT * FROM characters WHERE comicId = $1", id)
	if err != nil {
		r.logger.Error("Error getting comic. ", err)
		return nil, err
	}
	db.Close()
	return r.getCharactersResponse(rowsComic, rowsCharacter)
}

func (r *AlbomxComicsRepository) UpdateDatabase(response *dto.MarvelComicsResponse, characterName string) {
	r.logger.Info("Updating data base")
	collaborators := new(dto.CollaboratorsResponse)
	characters := new(dto.CharactersResponse)
	characters.Characters = make(map[string][]string)
	id := capAmericaId
	if strings.Contains(ironMan, characterName){
		id = ironManId
	}
	for _, result := range response.Data.Results {
		for _, collaborator := range  result.Creators.Items{
			switch collaborator.Role {
			case editor:
				r.logger.Info("Collaborator is editor")
				collaborators.Editors = append(collaborators.Editors, collaborator.Name)
			case writer:
				r.logger.Info("Collaborator is writer")
				collaborators.Writers = append(collaborators.Writers, collaborator.Name)
			case colorist:
				r.logger.Info("Collaborator is colorist")
				collaborators.Colorists = append(collaborators.Colorists, collaborator.Name)
			default:
				r.logger.Info("Collaborator isn't editor, writer or colorist")
			}
		}
		for _, character := range  result.Characters.Items{
			if comics := characters.Characters[character.Name]; comics != nil{
				characters.Characters[character.Name] = append(characters.Characters[character.Name], result.Title)
			}else{
				characters.Characters[character.Name] = []string{result.Title}
			}
		}
	}
	if err := r.saveData(id, collaborators, characters); err != nil {
		r.logger.Errorf("Error updating comic %s database. ", id, err)
		return
	}
}

func (r *AlbomxComicsRepository) createComicsTable() error {
	r.logger.Info("Creating comics database")
	query := `CREATE TABLE IF NOT EXISTS comics (
				  characterName varchar(45) PRIMARY KEY UNIQUE NOT NULL,
				  lastSync varchar(100)
				)
			`
	db := r.GetConnection()
	if _, err := db.Exec(query); err != nil {
		r.logger.Error("Error creating comics database. ", query, err)
		return err
	}
	db.Close()
	r.logger.Info("Created comics database")
	return nil
}

func (r *AlbomxComicsRepository) createCollaboratorsTable() error {
	r.logger.Info("Creating collaborators database")
	query := `CREATE TABLE IF NOT EXISTS collaborators (
				  collaboratorName varchar(45) NOT NULL,
				  rol varchar(20),
				  comicId varchar(20),
				  uniqueId varchar(200) PRIMARY KEY
		)
	`
	db := r.GetConnection()
	if _, err := db.Exec(query); err != nil {
		r.logger.Error("Error creating collaborators database. ", query, err)
		return err
	}
	db.Close()
	r.logger.Info("Created collaborators database")
	return nil
}

func (r *AlbomxComicsRepository) createCharactersTable() interface{} {
	r.logger.Info("Creating characters database")
	query := `CREATE TABLE IF NOT EXISTS characters (
				  characterName varchar(45) NOT NULL,
				  comicName varchar(100),
				  comicId varchar(20),
				  uniqueId varchar(200) PRIMARY KEY
		)
	`
	db := r.GetConnection()
	if _, err := db.Exec(query); err != nil {
		r.logger.Error("Error creating characters database. ", query, err)
		return err
	}
	db.Close()
	r.logger.Info("Created characters database")
	return nil
}

func (r *AlbomxComicsRepository) saveData(id string, collaborators *dto.CollaboratorsResponse, characters *dto.CharactersResponse) error {
	r.logger.Info("Saving data")
	if err := r.updateComicsTable(id); err != nil {
		return err
	}
	if err := r.updateCollaboratorsTable(id, collaborators); err != nil {
		return err
	}
	if err := r.updateCharactersTable(id, characters); err != nil {
		return err
	}
	return nil
}

func (r *AlbomxComicsRepository) updateComicsTable(id string) error {
	lastSync := "Last sync date " + time.Now().String()
	query := `INSERT INTO comics
				(characterName, lastSync)
			VALUES
				($1, $2)
			ON CONFLICT (characterName) DO UPDATE
			SET lastSync = $2`
	db := r.GetConnection()
	if _, err := db.Query(query, id, lastSync); err != nil {
		r.logger.Error("Error added new comicId. ", query, err)
		return err
	}
	db.Close()
	r.logger.Infof("Updated comic dataBase. Id %s", id)
	return nil
}

func (r *AlbomxComicsRepository) updateCollaboratorsTable(id string, collaborators *dto.CollaboratorsResponse) error {
	r.logger.Infof("Updating collaborators for Id %s", id)
	query := `INSERT INTO collaborators
				(collaboratorName, rol, comicId, uniqueId)
			VALUES
				($1, $2, $3, $4)
			ON CONFLICT (uniqueId) DO NOTHING`
	db := r.GetConnection()
	for _, collaborator := range collaborators.Colorists {
		uniqueId := fmt.Sprintf("%s%s%s", collaborator, colorist, id)
		if _, err := db.Query(query, collaborator, colorist, id, uniqueId); err != nil {
			r.logger.Error("Error added new comicId. ", query, err)
			return err
		}
	}

	for _, collaborator := range collaborators.Editors {
		uniqueId := fmt.Sprintf("%s%s%s", collaborator, editor, id)
		if _, err := db.Query(query, collaborator, editor, id, uniqueId); err != nil {
			r.logger.Error("Error added new comicId. ", query, err)
			return err
		}
	}

	for _, collaborator := range collaborators.Writers {
		uniqueId := fmt.Sprintf("%s%s%s", collaborator, writer, id)
		if _, err := db.Query(query, collaborator, writer, id, uniqueId); err != nil {
			r.logger.Error("Error added new comicId. ", query, err)
			return err
		}
	}
	db.Close()

	r.logger.Infof("Updating collaborators for Id %s", id)
	return nil
}

func (r *AlbomxComicsRepository) updateCharactersTable(id string, characters *dto.CharactersResponse) error {
	r.logger.Infof("Updating collaborators for Id %s", id)
	query := `INSERT INTO characters
				(characterName, comicName, comicId, uniqueId)
			VALUES
				($1, $2, $3, $4)
			ON CONFLICT (uniqueId) DO NOTHING`
	db := r.GetConnection()
	for characterName, comicsName := range characters.Characters {
		for _, comicName := range comicsName{
			uniqueId := fmt.Sprintf("%s%s%s", characterName, comicName, id)
			if _, err := db.Query(query, characterName, comicName, id, uniqueId); err != nil {
				r.logger.Error("Error added new comicId. ", query, err)
				return err
			}
		}
	}
	db.Close()

	r.logger.Infof("Updating characters for Id %s", id)
	return nil
}

func (r *AlbomxComicsRepository) getCollaboratorsResponse(rowsComic *sql.Rows, rowsCollaborator *sql.Rows) (*dto.CollaboratorsResponse, error) {
	response := dto.CollaboratorsResponse{}
	for rowsComic.Next() {
		var comic dto.ComicDb
		if err := rowsComic.Scan(&comic.CharacterName, &comic.LastSync); err != nil {
			r.logger.Error("Error getting comic. ", err)
			return nil, err
		}
		response.LastSync = comic.LastSync
		for rowsCollaborator.Next(){
			collaborator := new(dto.CollaboratorDb)
			if err := rowsCollaborator.Scan(&collaborator.CollaboratorName, &collaborator.Role, &collaborator.ComicId, &collaborator.UniqueId); err != nil {
				r.logger.Error("Error getting collaborator. ", err)
				return nil, err
			}
			switch collaborator.Role {
			case editor:
				r.logger.Info("Collaborator is editor")
				response.Editors = append(response.Editors, collaborator.CollaboratorName)
			case writer:
				r.logger.Info("Collaborator is writer")
				response.Writers = append(response.Writers, collaborator.CollaboratorName)
			case colorist:
				r.logger.Info("Collaborator is colorist")
				response.Colorists = append(response.Colorists, collaborator.CollaboratorName)
			default:
				r.logger.Info("Collaborator isn't editor, writer or colorist")
			}
		}

	}
	r.logger.Info("Get collaborators")
	return &response,	nil
}

func (r *AlbomxComicsRepository) getCharactersResponse(rowsComic *sql.Rows, rowsCharacter *sql.Rows) (*dto.CharactersResponse, error) {
	response := new(dto.CharactersResponse)
	response.Characters = make(map[string][]string)
	for rowsComic.Next() {
		var comic dto.ComicDb
		if err := rowsComic.Scan(&comic.CharacterName, &comic.LastSync); err != nil {
			r.logger.Error("Error getting comic. ", err)
			return nil, err
		}
		response.LastSync = comic.LastSync
		for rowsCharacter.Next(){
			character := new(dto.CharacterDb)
			if err := rowsCharacter.Scan(&character.CharacterName, &character.ComicName, &character.ComicId, &character.UniqueId); err != nil {
				r.logger.Error("Error getting character. ", err)
				return nil, err
			}
			if comics := response.Characters[character.CharacterName]; comics != nil{
				response.Characters[character.CharacterName] = append(response.Characters[character.CharacterName], character.ComicName)
			}else{
				response.Characters[character.CharacterName] = []string{character.ComicName}
			}
		}

	}
	r.logger.Info("Get characters")
	return response,	nil
}