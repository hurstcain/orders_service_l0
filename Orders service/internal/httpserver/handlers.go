package httpserver

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"sync"
)

import (
	"orders_service/internal/model"
)

// StartServer - настраивает маршрутизацию и запускает сервер
func StartServer(orders *[]model.Order, mu *sync.Mutex) {
	r := mux.NewRouter()

	// Настройка стартевой страницы с формой ввода uid заказа
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tmpl := template.Must(template.ParseFiles("static/start_page.html"))
		tmpl.Execute(writer, nil)
	})

	// Настройка страницы /order, на которой выводится информация о заказе
	r.HandleFunc("/order", func(writer http.ResponseWriter, request *http.Request) {
		orderUid := request.FormValue("order_uid")
		order := model.Order{}

		mu.Lock()
		for _, tempOrder := range *orders {
			if tempOrder.OrderUid == orderUid {
				order = tempOrder
				break
			}
		}
		mu.Unlock()

		if order.OrderUid == "" {
			tempStruct := struct {
				UidInput string
			}{
				UidInput: orderUid,
			}
			tmpl := template.Must(template.ParseFiles("static/error.html"))
			tmpl.Execute(writer, tempStruct)
		} else {
			tmpl := template.Must(template.ParseFiles("static/order_info_page.html"))
			tmpl.Execute(writer, order)
		}
	})

	http.Handle("/", r)

	log.Println("Start server on 127.0.0.1:8000...")
	if err := http.ListenAndServe("127.0.0.1:8000", nil); err != nil {
		log.Fatalf("Server failed to start. Error: %s\n", err)
	}
}
