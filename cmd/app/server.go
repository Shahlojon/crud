package app

import (
	"github.com/Shahlojon/crud/cmd/app/middleware"
	"github.com/Shahlojon/crud/pkg/customers/security"
	"github.com/gorilla/mux"
	"github.com/Shahlojon/crud/pkg/customers"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

)

//Server ...
type Server struct {
	mux         *mux.Router
	customerSvc *customers.Service
	securitySvc *security.Service
}

//NewServer ...
func NewServer(m *mux.Router, cSvc *customers.Service, sSvc *security.Service) *Server {
	return &Server{mux: m, customerSvc: cSvc, securitySvc:sSvc,}
}

func (s *Server)ServeHTTP(w http.ResponseWriter, r *http.Request){
	s.mux.ServeHTTP(w,r)
}

const (
	GET = "GET"
	POST = "POST"
	DELETE = "DELETE"
)

//Init ...
func (s *Server) Init() {
	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers", s.handleSave).Methods(POST)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	//s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnBlockByID).Methods(DELETE)
	//s.mux.HandleFunc("/customers.removeById", s.handleDelete)
	s.mux.HandleFunc("/customers/{id}", s.handleDelete).Methods(DELETE)
	//s.mux.HandleFunc("/customers.save", s.handleSave)

	s.mux.Use(middleware.Basic(s.securitySvc.Auth))

}

// хендлер метод для извлечения всех клиентов
func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {

	items, err :=s.customerSvc.All(r.Context())
	if err != nil{
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, items)
}

// хендлер метод для извлечения всех активных клиентов
func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request) {

	items, err :=s.customerSvc.AllActive(r.Context())
	if err != nil{
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	

	respondJSON(w, items)
}

func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	idP, ok := mux.Vars(r)["id"]
	if !ok{
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

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
	respondJSON(w, item)
}

func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	idP, ok := mux.Vars(r)["id"]
	if !ok{
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//получаем баннер из сервиса
	item, err := s.customerSvc.ChangeActive(r.Context(), id, false)

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
	respondJSON(w, item)
}

func (s *Server) handleUnBlockByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	idP, ok := mux.Vars(r)["id"]
	if !ok{
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

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
	respondJSON(w, item)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	idP, ok := mux.Vars(r)["id"]
	if !ok{
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//получаем баннер из сервиса
	item, err := s.customerSvc.Delete(r.Context(), id)

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
	respondJSON(w, item)
}


func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {

	// //получаем данные из параметра запроса
	// idP := r.FormValue("id")
	// name := r.FormValue("name")
	// phone := r.FormValue("phone")

	// id, err := strconv.ParseInt(idP, 10, 64)
	// //если получили ошибку то отвечаем с ошибкой
	// if err != nil {
	// 	//вызываем фукцию для ответа с ошибкой
	// 	errorWriter(w, http.StatusBadRequest, err)
	// 	return
	// }
	// //Здесь опционалная проверка то что если все данные приходит пустыми то вернем ошибку
	// if name == "" && phone == ""  {
	// 	//вызываем фукцию для ответа с ошибкой
	// 	errorWriter(w, http.StatusBadRequest, err)
	// 	return
	// }

	// item := &customers.Customer{
	// 	ID:id,
	// 	Name:name,
	// 	Phone:phone,
	// 	/* Active:true,
	// 	Created:time.Now() */
	// }

	// customer, err := s.customerSvc.Save(r.Context(), item)

	// //если получили ошибку то отвечаем с ошибкой
	// if err != nil {
	// 	//вызываем фукцию для ответа с ошибкой
	// 	errorWriter(w, http.StatusInternalServerError, err)
	// 	return
	// }
	// //вызываем функцию для ответа в формате JSON
	// respondJSON(w, customer)
	var item *customers.Customer
	err:=json.NewDecoder(r.Body).Decode(&item)
	if err!=nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err = s.customerSvc.Save(r.Context(), item)

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

//это фукция для записывание ошибки в responseWriter или просто для ответа с ошиками
func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	//печатаем ошибку
	log.Print(err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}


//это функция для ответа в формате JSON
func respondJSON(w http.ResponseWriter, iData interface{}) {

	//преобразуем данные в JSON
	data, err := json.Marshal(iData)

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		//печатаем ошибку
		log.Print(err)
	}
}
