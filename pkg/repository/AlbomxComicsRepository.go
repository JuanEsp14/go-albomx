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
				collaborators.Writers = append(collaborators.Editors, collaborator.Name)
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
		uniqueId := fmt.Sprintf("%s%s%s", collaborator, colorist, id)
		if _, err := db.Query(query, collaborator, colorist, id, uniqueId); err != nil {
			r.logger.Error("Error added new comicId. ", query, err)
			return err
		}
	}

	for _, collaborator := range collaborators.Writers {
		uniqueId := fmt.Sprintf("%s%s%s", collaborator, colorist, id)
		if _, err := db.Query(query, collaborator, colorist, id, uniqueId); err != nil {
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

/**
db := r.GetConnection()
	rows, err := db.Query("SELECT * FROM collaborators")
	if err != nil {
		r.logger.Error("Error added new comicId. ", err)
		return err
	}
	db.Close()
	for rows.Next() {
		var comic dto.CollaboratorDb
		if err = rows.Scan(&comic.CollaboratorName, &comic.Rol, &comic.ComicId, &comic.UniqueId); err != nil {
			r.logger.Error("Error added new comicId. ", err)
			return err
		}
		r.logger.Info("CollaboratorName ", comic.CollaboratorName)
		r.logger.Info("Rol ", comic.Rol)
		r.logger.Info("ComicId ", comic.ComicId)

	}
	r.logger.Info(rows)
 */