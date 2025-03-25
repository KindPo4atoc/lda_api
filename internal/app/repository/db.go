package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DataBase struct {
	config         *Config
	db             *sql.DB
	dataRepository *DataForLearnRepository
}

func New(config *Config) *DataBase {
	return &DataBase{
		config: config,
	}
}

func (data *DataBase) Open() error {
	db, err := sql.Open("postgres", data.config.DatabaseURL)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	data.db = db

	return nil
}

func (db *DataBase) Close() {
	db.db.Close()
}

func (data *DataBase) Data() *DataForLearnRepository {
	if data.dataRepository != nil {
		return data.dataRepository
	}

	data.dataRepository = &DataForLearnRepository{
		store: data,
	}
	return data.dataRepository
}
