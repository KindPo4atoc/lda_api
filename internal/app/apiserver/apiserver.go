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

// Инициализация структуры
type APIServer struct {
	config  *Config
	logger  *logrus.Logger
	router  *mux.Router
	db      *repository.DataBase
	model   *lda.LDA
	predict *lda.PredictModel
}

// конструктор
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Метод запуска API. Инициализирует все поля структуры
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info("Начало конфигурации маршрутизатора.")
	s.configureRouter()
	s.logger.Info("Начало конфигурация базы данных.")
	if err := s.configureDB(); err != nil {
		return err
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500"}, // Разрешенные домены
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Разрешить куки и заголовки авторизации
		Debug:            true, // Логирование (опционально)
	})
	s.logger.Info("Конфигурация политики CORS завершена.")
	s.model = lda.New(s.db)
	s.logger.Info("Начало инициализация модели.")
	if err := s.model.FitModel(); err != nil {
		return err
	}
	handler := c.Handler(s.router)
	s.logger.Info("API сервер запущен.")
	return http.ListenAndServe(s.config.BindAddr, handler)
}

// инициализация логгера
func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)

	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

// инициализация бд
func (s *APIServer) configureDB() error {
	database := repository.New(s.config.DBConfig)
	if err := database.Open(); err != nil {
		return err
	}

	s.db = database
	s.db.Data()
	s.logger.Info("Конфигурация базы данных завершена.")
	return nil
}

// ручка на получение обучающих данных по маршруту -> /selectLearnData
func (s *APIServer) handleSelectLearnData(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Использован GET-запрос по маршруту http://localhost:8000/selectLearnData")
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

// ручка на получение параметров модели по маршруту -> /initModel
func (s *APIServer) handleGetParamModel(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Использован GET-запрос по маршруту http://localhost:8000/getParamModel")
	w.Header().Set("Content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(s.model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ручка на получение предикта post запрос, тело должно содержать userData, маршрут -> /predict
func (s *APIServer) handlePredictModel(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Использован POST-запрос по маршруту http://localhost:8000/predict")
	var dataForPredict entity.UserData
	json.NewDecoder(r.Body).Decode(&dataForPredict)
	var tmp []float64
	tmp = append(tmp, float64(dataForPredict.IncomeAnnum))
	tmp = append(tmp, float64(dataForPredict.LoanAmount))
	tmp = append(tmp, float64(dataForPredict.LoanTerm))
	tmp = append(tmp, float64(dataForPredict.CibilScore))
	var data [][]float64
	data = append(data, tmp)
	predict, distance, dataLDA, err := s.model.Predict(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	s.predict = lda.NewPredict(predict, distance, dataLDA)
	w.Header().Set("Content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(s.predict)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// инициализация роутера
func (s *APIServer) configureRouter() {

	s.router.HandleFunc("/selectLearnData", s.handleSelectLearnData).Methods("GET")
	s.router.HandleFunc("/getParamModel", s.handleGetParamModel).Methods("GET")
	s.router.HandleFunc("/predict", s.handlePredictModel).Methods("POST")
	s.logger.Info("Конфигурация маршрутизатора завершена.")
}
