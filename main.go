package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Response struct {
	Success bool    `json:"Success"`
	ErrCode string  `json:"ErrCode"`
	Value1  float64 `json:"Value1"`
	Value2  float64 `json:"Value2,omitempty"`
}

func ErrorResponse(w http.ResponseWriter, err string) {
	response := Response{
		Success: false,
		ErrCode: err,
	}
	jsonResponse(w, response)
}
func jsonResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Main page ")
}

func OperationHandler(w http.ResponseWriter, r *http.Request) {

	as := r.URL.Query().Get("a") //получение параметра a
	bs := r.URL.Query().Get("b") //получение параметра b

	vars := mux.Vars(r)
	operation := vars["operation"]

	if len(r.URL.Query()) == 2 {
		a, b, err := CheckDigit(as, bs)
		if err != nil {
			ErrorResponse(w, err.Error())
			return
		}
		response := Ariphmetical(w, operation, a, b)
		jsonResponse(w, response)
	} else {
		err := errors.New("укажите 2 парметра")
		ErrorResponse(w, err.Error())
	}

}

func CheckDigit(a, b string) (float64, float64, error) {
	if len(a) <= 0 || len(b) <= 0 {
		return 0, 0, errors.New("неправильно названа переменная или пустое значение")
	} else {
		a, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return 0, 0, errors.New("ошибка в чтении числа a")
		}
		b, err := strconv.ParseFloat(b, 64)
		if err != nil {
			return 0, 0, errors.New("ошибка в чтении числа b")
		}
		return a, b, nil
	}
}

func Ariphmetical(w http.ResponseWriter, operation string, a, b float64) (response Response) {

	switch operation {
	case "add":
		response = Response{Success: true, Value1: a + b}
	case "sub":
		response = Response{Success: true, Value1: a - b, Value2: b - a}
	case "mul":
		response = Response{Success: true, Value1: a * b}
	case "div":
		if a != 0 && b != 0 {
			response = Response{Success: true, Value1: a / b, Value2: b / a}
		} else {
			err := errors.New("деление на ноль")
			response = Response{Success: false, ErrCode: err.Error()}
		}
	default:
		err := errors.New("неизвестная операция")
		response = Response{Success: false, ErrCode: err.Error()}
	}
	return
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	apiSubR := r.PathPrefix("/api").Subrouter()
	apiSubR.HandleFunc("/{operation}", OperationHandler).Methods(http.MethodGet)
	http.Handle("/", r)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)

}
