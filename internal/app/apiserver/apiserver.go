package apiserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lda_api/internal/app/entity"
	"github.com/lda_api/internal/app/repository"
	logrus "github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	db     *repository.DataBase
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureDB(); err != nil {
		return nil
	}

	s.logger.Info("Starting api server")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)

	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *APIServer) configureDB() error {
	database := repository.New(s.config.DBConfig)
	if err := database.Open(); err != nil {
		return err
	}

	s.db = database
	s.db.Data()

	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/selectUser", s.handleSelectUser)
	s.router.HandleFunc("/fitModel", s.FitModel()).Methods("GET")
}

func (s *APIServer) handleSelectUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	data, err := s.db.Data().SelectAllData()

	if err != nil {
		logrus.Fatal(err)
	}
	dataForJson := entity.New(data)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(entity.New(dataForJson.Data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (s *APIServer) FitModel() http.HandlerFunc {
	// TODO -> реализовать алгоритм lda
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hui")
	}
}
