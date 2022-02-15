package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

const (
	// Id кластера
	clusterID = "test-cluster"
	// Id клиента, публикующего данные в канал
	clientID = "publisher"
	// URL для подключения к NATS
	URL = nats.DefaultURL
	// Название канала
	subj = "foo"
)

// Генерирует и возвращает строку в формате json, которая будет передаваться в канал.
// Уникальные id заказа (order_uid), id транзакции (transaction) генерируются рандомно.
func generateRandomJson() string {
	itemsCount := rand.Intn(10) + 1
	items := generateRandomItems(itemsCount)
	json := fmt.Sprintf("{\"order_uid\": \"test%d\", \"track_number\": \"WBILMTESTTRACK\", "+
		"\"entry\": \"WBIL\", \"delivery\": {\"name\": \"Test Testov\", \"phone\": \"+9720000000\",\n"+
		"\"zip\": \"2639809\", \"city\": \"Kiryat Mozkin\", \"address\": \"Ploshad Mira 15\", "+
		"\"region\": \"Kraiot\", \"email\": \"test@gmail.com\"}, \"payment\": {\"transaction\": \"test%d\", "+
		"\"request_id\": \"\", \"currency\": \"USD\", \"provider\": \"wbpay\", \"amount\": %d, "+
		"\"payment_dt\": %d, \"bank\": \"alpha\", \"delivery_cost\": 1500, \"goods_total\": %d, "+
		"\"custom_fee\": 0}, \"items\": [", rand.Int(), rand.Int(), 1500+itemsCount*317, time.Now().Unix(), itemsCount*317) +
		items + fmt.Sprintf("], \"locale\": \"en\", \"internal_signature\": \"\", \"customer_id\": \"test\", "+
		"\"delivery_service\": \"meest\", \"shardkey\": \"9\", \"sm_id\": 99, "+
		"\"date_created\": \"%s\", \"oof_shard\": \"1\"}", time.Now().Format(time.RFC3339))

	return json
}

// Генерирует и возвращает строку со списком из count (1<=count<=10) позиций заказа в формате json.
// Уникальные id позиций (rid) генерируются рандомно.
func generateRandomItems(count int) string {
	var items string

	for i := 0; i < count; i++ {
		if i == 0 {
			items += fmt.Sprintf("{\"chrt_id\": 9934930, \"track_number\": \"WBILMTESTTRACK\", "+
				"\"price\": 453, \"rid\": \"testrid%d\", \"name\": \"Mascaras\", \"sale\": 30, \"size\": \"0\", "+
				"\"total_price\": 317, \"nm_id\": 2389212, \"brand\": \"Vivienne Sabo\", \"status\": 202}", rand.Int())
			continue
		}
		items += fmt.Sprintf(", {\"chrt_id\": 9934930, \"track_number\": \"WBILMTESTTRACK\", "+
			"\"price\": 453, \"rid\": \"testrid%d\", \"name\": \"Mascaras\", \"sale\": 30, \"size\": \"0\", "+
			"\"total_price\": 317, \"nm_id\": 2389212, \"brand\": \"Vivienne Sabo\", \"status\": 202}", rand.Int())
	}

	return items
}

func main() {
	// Соединение с NATS
	log.Println("Connecting to NATS...")
	nc, err := nats.Connect(URL)
	if err != nil {
		log.Fatalf("Can't connect to NATS. Error: %s\n", err)
	}
	log.Print("Connection to NATS is established\n")
	defer nc.Close()

	// Соединение с NATS Streaming server
	log.Println("Connecting to NATS Streaming server...")
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc), stan.MaxPubAcksInflight(1),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost. Reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect to NATS Streaming. Error: %s\n", err)
	}
	log.Print("Connection to NATS Streaming server is established...\n")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Printf("Received an interrupt, closing NATS Streaming connection and stop publishing...\n\n")
			sc.Close()
		}
	}()

	// Публикация сообщений в канал каждые 8 секунд
	// Цикл прерывается, когда sc.Publish() возвращает не пустую ошибку,
	// то есть когда публикация в канал становится невозможна
	for i := 1; ; i++ {
		json := generateRandomJson()
		err := sc.Publish(subj, []byte(json))
		if err != nil {
			log.Fatalf("Unable to write data to the channel. Reason: %s", err)
		}
		log.Printf("Message №%d send to channel %s\n", i, subj)

		time.Sleep(8 * time.Second)
	}
}
