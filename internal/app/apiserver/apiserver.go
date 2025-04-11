package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lda_api/internal/app/entity"
	"github.com/lda_api/internal/app/lda"
	"github.com/lda_api/internal/app/repository"
	"github.com/rs/cors"
	logrus "github.com/sirupsen/logrus"
)

type APIServer struct {
	config  *Config
	logger  *logrus.Logger
	router  *mux.Router
	db      *repository.DataBase
	model   *lda.LDA
	predict *lda.PredictModel
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
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500"}, // Разрешенные домены
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Разрешить куки и заголовки авторизации
		Debug:            true, // Логирование (опционально)
	})
	handler := c.Handler(s.router)
	return http.ListenAndServe(s.config.BindAddr, handler)
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

func (s *APIServer) handleSelectLearnData(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /selectLearnData: GET request")
	w.Header().Set("Content-type", "application/json")
	data, err := s.db.Data().SelectAllLearnData()

	if err != nil {
		logrus.Fatal(err)
	}
	dataForJson := entity.New(data.Data)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(entity.New(dataForJson.Data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (s *APIServer) handleSelectTestData(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /selectTestData: GET request")
	w.Header().Set("Content-type", "application/json")
	data, err := s.db.Data().SelectTestData()

	if err != nil {
		logrus.Fatal(err)
	}
	dataForJson := entity.New(data.Data)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(entity.New(dataForJson.Data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (s *APIServer) handleSelectClassData(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /selectClassData/{item}: GET request")
	w.Header().Set("Content-type", "application/json")
	data, err := s.db.Data().SelectingDataByClass(mux.Vars(r)["item"])

	if err != nil {
		logrus.Fatal(err)
	}
	dataForJson := entity.New(data.Data)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(entity.New(dataForJson.Data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (s *APIServer) handleInitModel(w http.ResponseWriter, r *http.Request) {
	// TODO ->
	// 		->
	//		->
	//		->

	logrus.Info("Route /initModel: GET request")
	s.model = lda.New(s.db)
	if err := s.model.FitModel(); err != nil {
		logrus.Fatal(err)
	}
	w.Header().Set("Content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(s.model)

	/*data, err := s.db.Data().SelectingDataByClass(mux.Vars(r)["item"])
	if err != nil {
		logrus.Fatal(err)
	}
	dataForJson := entity.New(data.Data)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(entity.New(dataForJson.Data))*/
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *APIServer) handlePredictModel(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /predict: POST request")
	var dataForPredict entity.UserData
	json.NewDecoder(r.Body).Decode(&dataForPredict)

	s.predict = lda.NewPredict(s.model.PredictModel(dataForPredict))

	w.Header().Set("Content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(s.predict)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *APIServer) handleGetConvData(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /getConvData: GET request")
	w.Header().Set("Content-type", "application/json")
	data, err := s.db.Data().SelectAllLearnData()

	if err != nil {
		logrus.Fatal(err)
	}
	dataConv := s.model.ConvertData(data)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(dataConv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) configureRouter() {

	s.router.HandleFunc("/selectLearnData", s.handleSelectLearnData).Methods("GET")
	s.router.HandleFunc("/selectTestData", s.handleSelectTestData).Methods("GET")
	s.router.HandleFunc("/getConvData", s.handleGetConvData).Methods("GET")
	s.router.HandleFunc("/selectClassData/{item}", s.handleSelectClassData).Methods("GET")
	s.router.HandleFunc("/initModel", s.handleInitModel).Methods("GET")
	s.router.HandleFunc("/predict", s.handlePredictModel).Methods("POST")

}
