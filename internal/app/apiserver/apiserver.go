package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lda_api/internal/app/entity"
	"github.com/lda_api/internal/app/lda"
	"github.com/lda_api/internal/app/repository"
	logrus "github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	db     *repository.DataBase
	model  *lda.LDA
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
	s.router.HandleFunc("/initData/{item}", s.handleFitModel).Methods("GET")
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
func (s *APIServer) handleFitModel(w http.ResponseWriter, r *http.Request) {
	// TODO -> подключить библиотеку для расчета собственных значений и векторов
	// 		-> дописать алгоритм расчета lda,
	//		-> + реализовать вывод параметров модели
	//		-> + начать реализовывать Post запрос на предикт от модели
	//
	// инициализация модели -> отправка в виде сериализация в json и дальнейшая отправка на клиент
	//

	s.model = lda.New(s.db)
	if err := s.model.FitModel(); err != nil {
		logrus.Fatal(err)
	}
	w.Header().Set("Content-type", "application/json")
	routerVariable := " " + mux.Vars(r)["item"]
	fmt.Println(routerVariable)
	data, err := s.db.Data().SelectingDataByClass(routerVariable)
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
