package main

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
	"log"
	"orders_service/internal/httpserver"
	"os"
	"os/signal"
	"sync"
)

import (
	"orders_service/internal/db"
	"orders_service/internal/model"
	"orders_service/internal/nats"
)

// streamingServerMsgHandler - обработчик сообщения, полученного из канала.
// Здесь осуществляется запись данных в бд, а также в память программы (в слайс orders)
func streamingServerMsgHandler(msg *stan.Msg, orders *[]model.Order, mu *sync.Mutex, database *sqlx.DB, i int) {
	log.Printf("№%d message recieved.\n", i)
	log.Println(string(msg.Data))

	order := model.Order{}
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Printf("Wrong data format was delivered. Error: %s\n", err)
		return
	}
	if err := order.CheckData(); err != nil {
		log.Printf("Invalid order fields. Error: %s\n", err)
		return
	}

	for i, _ := range order.Items {
		order.Items[i].OrderUid = order.OrderUid
	}

	log.Println("Adding message data to the database...")

	tx := database.MustBegin()

	if deliveryId, err := db.InsertDelivery(database, order.Delivery); err != nil {
		return
	} else {
		order.Delivery.Id = deliveryId
		order.DeliveryId = deliveryId
	}

	if err := db.InsertPayment(tx, order.Payment); err != nil {
		return
	}
	order.PaymentId = order.Payment.Transaction

	if err := db.InsertOrder(tx, order); err != nil {
		return
	}

	if err := db.InsertItem(tx, order.Items); err != nil {
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error when commit the transaction: %s", err)
	} else {
		log.Println("Message data successfully added to the database")
	}

	mu.Lock()
	*orders = append(*orders, order)
	mu.Unlock()
	log.Println("Message data successfully added to the memory")
}

func main() {
	signalChan := make(chan os.Signal, 1)
	exitChan := make(chan bool, 2)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Printf("Received an interrupt, closing all connections and stop listening to the channel...\n\n")
			exitChan <- true
			exitChan <- true
		}
	}()

	// Подкючение к базе данных
	database := db.Connect()
	defer database.Close()

	// Слайс структур model.Order, где будут временно храниться данные, получаемые из канала
	orders := make([]model.Order, 0)
	ordersMutex := sync.Mutex{}

	// Выгрузка данных из бд в orders
	db.UploadOrdersData(database, &orders)

	// Запуск http сервера
	go httpserver.StartServer(&orders, &ordersMutex)

	// Подключение к NATS серверу и подписка на канал
	var i int
	streamingServerConn := nats.StreamingConnection{
		ClusterID:   "test-cluster",
		ClientID:    "reader",
		URL:         stan.DefaultNatsURL,
		Subject:     "foo",
		DurableName: "cool_reader",
		MessageHandler: func(msg *stan.Msg) {
			msg.Ack()
			i++
			streamingServerMsgHandler(msg, &orders, &ordersMutex, database, i)
		},
	}
	go streamingServerConn.SubscribeToChannel(exitChan)

	<-exitChan
}
