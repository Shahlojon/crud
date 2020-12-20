package app

import (
	"github.com/Shahlojon/crud/cmd/app/middleware"
	"encoding/json"

	"log"
	"net/http"

	"github.com/Shahlojon/crud/pkg/customers"
	"github.com/Shahlojon/crud/pkg/managers"
	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

//Server ...
type Server struct {
	mux         *mux.Router
	customerSvc *customers.Service
	managerSvc  *managers.Service
}

//NewServer ... создает новый сервер
func NewServer(m *mux.Router, cSvc *customers.Service, mSvc *managers.Service) *Server {
	return &Server{
		mux:         m,
		customerSvc: cSvc,
		managerSvc:  mSvc,
	}
}

// функция для запуска хендлеров через мукс
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//Init ... инициализация сервера
func (s *Server) Init() {
	customersAuthenticateMd := middleware.Authenticate(s.customerSvc.IDByToken)
	customersSubrouter := s.mux.PathPrefix("/api/customers").Subrouter()
	customersSubrouter.Use(customersAuthenticateMd)

	customersSubrouter.HandleFunc("", s.handleCustomerRegistration).Methods(POST)
	customersSubrouter.HandleFunc("/token", s.handleCustomerGetToken).Methods(POST)
	customersSubrouter.HandleFunc("/products", s.handleCustomerGetProducts).Methods(GET)

	managersAuthenticateMd := middleware.Authenticate(s.managerSvc.IDByToken)
	managersSubRouter := s.mux.PathPrefix("/api/managers").Subrouter()
	managersSubRouter.Use(managersAuthenticateMd)
	managersSubRouter.HandleFunc("", s.handleManagerRegistration).Methods(POST)
	managersSubRouter.HandleFunc("/token", s.handleManagerGetToken).Methods(POST)
	managersSubRouter.HandleFunc("/sales", s.handleManagerGetSales).Methods(GET)
	managersSubRouter.HandleFunc("/sales", s.handleManagerMakeSales).Methods(POST)
	managersSubRouter.HandleFunc("/products", s.handleManagerGetProducts).Methods(GET)
	managersSubRouter.HandleFunc("/products", s.handleManagerChangeProducts).Methods(POST)
	managersSubRouter.HandleFunc("/products/{id:[0-9]+}", s.handleManagerRemoveProductByID).Methods(DELETE)
	managersSubRouter.HandleFunc("/customers", s.handleManagerGetCustomers).Methods(GET)
	managersSubRouter.HandleFunc("/customers", s.handleManagerChangeCustomer).Methods(POST)
	managersSubRouter.HandleFunc("/customers/{id:[0-9]+}", s.handleManagerRemoveCustomerByID).Methods(DELETE)

}


//это фукция для записывание ошибки в responseWriter или просто для ответа с ошиками
func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	//печатаем ошибку
	log.Print(err)
	//отвечаем ошибку с помошю библиотеке net/http
	http.Error(w, http.StatusText(httpSts), httpSts)
}

//это функция для ответа в формате JSON (он принимает интерфейс по этому мы можем в нем передат все что захочется)
func respondJSON(w http.ResponseWriter, iData interface{}) {

	//преобразуем данные в JSON
	data, err := json.Marshal(iData)

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//поставить хедер "Content-Type: application/json" в ответе
	w.Header().Set("Content-Type", "application/json")
	//пишем ответ
	_, err = w.Write(data)
	//если получили ошибку
	if err != nil {
		//печатаем ошибку
		log.Print(err)
	}
}

/*
// хендлер метод для извлечения всех клиентов
func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {

	//берем все клиенты
	items, err := s.customerSvc.All(r.Context())

	//если ест ошибка
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, items)
}

// хендлер метод для извлечения всех активных клиентов
func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request) {

	//берем все активные клиенты
	items, err := s.customerSvc.AllActive(r.Context())

	//если ест ошибка
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, items)
}

//хендлер который верет по айди
func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	//idP := r.URL.Query().Get("id")
	idP := mux.Vars(r)["id"]

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//получаем баннер из сервиса
	item, err := s.customerSvc.ByID(r.Context(), id)

	//если ошибка равно на notFound то вернем ошибку не найдено
	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, item)
}

//хендлер для блокировки
func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	//idP := r.URL.Query().Get("id")
	idP := mux.Vars(r)["id"]

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//изменяем статус клиента на фалсе
	item, err := s.customerSvc.ChangeActive(r.Context(), id, false)
	//если ошибка равно на notFound то вернем ошибку не найдено
	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, item)
}

//хенндлер для разблокировки
func (s *Server) handleUnBlockByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	//idP := r.URL.Query().Get("id")
	idP := mux.Vars(r)["id"]

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//получаем баннер из сервиса
	item, err := s.customerSvc.ChangeActive(r.Context(), id, true)
	//если ошибка равно на notFound то вернем ошибку не найдено
	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//вызываем функцию для ответа в формате JSON
	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, item)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	//idP := r.URL.Query().Get("id")
	idP := mux.Vars(r)["id"]

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//удаляем клиента из базу
	item, err := s.customerSvc.Delete(r.Context(), id)
	//если ошибка равно на notFound то вернем ошибку не найдено
	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//вызываем функцию для ответа в формате JSON
	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, item)
}

//хендлер для сохранения и обновления
func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {

	//обявляем структура клиента для запраса
	var item *customers.Customer

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//Генерируем bcrypt хеш от реалного пароля
	hashed, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//и поставляем хеш в поле парол
	item.Password = string(hashed)

	//сохроняем или обновляем клиент
	customer, err := s.customerSvc.Save(r.Context(), item)

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//вызываем функцию для ответа в формате JSON
	respondJSON(w, customer)
}

func (s *Server) handleCreateToken(w http.ResponseWriter, r *http.Request) {

	//обявляем структуру для запроса
	var item *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	//извелекаем данные из запраса
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	//взываем из сервиса  securitySvc метод AuthenticateCustomer
	token, err := s.securitySvc.TokenForCustomer(r.Context(), item.Login, item.Password)

	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//вызываем функцию для ответа в формате JSON
	respondJSON(w, map[string]interface{}{"status": "ok", "token": token})
}

//хендлер для валидации
func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	//создаем структуру для извлечения запроса
	var item *struct {
		Token string `json:"token"`
	}
	//из json извелекаем нужные данные
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//взываем из сервиса  securitySvc метод AuthenticateCustomer
	id, err := s.securitySvc.AuthenticateCustomer(r.Context(), item.Token)

	//если получили ошибу
	if err != nil {
		//то генурируем интернал еррор
		status := http.StatusInternalServerError
		text := "internal error"
		//если ошибка ErrNoSuchUser то изменим ответ на не найдено
		if err == security.ErrNoSuchUser {
			status = http.StatusNotFound
			text = "not found"
		}
		//если ошибка ErrExpireToken то изменим ответ на не expired
		if err == security.ErrExpireToken {
			status = http.StatusBadRequest
			text = "expired"
		}
		//отвечаем с кодом статуса
		respondJSONWithCode(w, status, map[string]interface{}{"status": "fail", "reason": text})
		return
	}

	//генерируем структуру ответа
	res := make(map[string]interface{})
	res["status"] = "ok"
	res["customerId"] = id

	//отвечаем с кодом статуса
	respondJSONWithCode(w, http.StatusOK, res)
}
*/
/*
+
+
+
+
+
+
+
*/
