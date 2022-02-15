package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
)

// StreamingConnection содержит в себе настройки для подключения к NATS Streaming, а также для подписки на канал
type StreamingConnection struct {
	// Id кластера
	ClusterID string
	// Id клиента, который читает данные из канала
	ClientID string
	// URL для соединения с NATS
	URL string
	// Название канала
	Subject string
	// Имя клиента для долгосрочной подписки на канал
	DurableName string
	// Функция для обработки полученных данных из канала
	MessageHandler func(msg *stan.Msg)
}

// SubscribeToChannel - функция подписки на канал
func (s StreamingConnection) SubscribeToChannel(exitCh chan bool) {
	// Соединение с NATS
	log.Println("Connecting to NATS...")
	nc, err := nats.Connect(s.URL)
	if err != nil {
		log.Fatalf("Can't connect to NATS. Error: %s\n", err)
	}
	defer nc.Close()
	log.Print("Connection to NATS is established\n")

	// Соединение с NATS Streaming server
	log.Println("Connecting to NATS Streaming server...")
	sc, err := stan.Connect(s.ClusterID, s.ClientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost. Reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect to NATS Streaming. Error: %s\n", err)
	}
	defer sc.Close()
	log.Printf("Connected to %s, clusterID: [%s], clientID: [%s]\n", s.URL, s.ClusterID, s.ClientID)

	// Подписка на канал
	log.Printf("Subscribing to the channel %s...\n", s.Subject)
	sub, err := sc.Subscribe(s.Subject, s.MessageHandler, stan.SetManualAckMode(), stan.DurableName(s.DurableName))
	if err != nil {
		log.Fatalf("Can't subscribe to the channel %s. Reason: %s", s.Subject, err)
	}
	defer sub.Close()
	log.Printf("Listening on the channel %s...\n", s.Subject)

	<-exitCh
}
