package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"orders_service/configs"
)

import (
	"orders_service/internal/model"
)

// Connect - устанавливает соединение с базой данных
func Connect() *sqlx.DB {
	log.Println("Connecting to database...")
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configs.DBUser, configs.DBPassword, configs.DBName))
	if err != nil {
		log.Fatalf("Can't connect to database. Error: %s\n", err)
	}
	log.Println("Connection to database is established")

	return db
}

// InsertDelivery - добавляет данные в таблицу deliveries,
// а также возвращает автоматически генерируемый id добавленной записи
func InsertDelivery(db *sqlx.DB, delivery model.Delivery) (int, error) {
	var id int

	stmt, err := db.PrepareNamed("INSERT INTO deliveries (name, phone, zip, city, address, region, email) " +
		"VALUES (:name, :phone, :zip, :city, :address, :region, :email) RETURNING id")
	if err != nil {
		log.Printf("Error when adding data to the table deliveries: %s", err)
		return 0, err
	}
	err = stmt.Get(&id, &delivery)
	if err != nil {
		log.Printf("Error when commit the transaction to the table deliveries: %s", err)
		return 0, err
	}

	return id, nil
}

// InsertPayment - добавляет данные в таблицу payments
func InsertPayment(tx *sqlx.Tx, payment model.Payment) error {
	_, err := tx.NamedExec("INSERT INTO payments (transaction, request_id, currency, provider, "+
		"amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES (:transaction, :request_id, "+
		":currency, :provider, :amount, :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee)", &payment)
	if err != nil {
		log.Printf("Error when adding data to the table payments: %s", err)
		return err
	}

	return nil
}

// InsertItem - добавляет данные в таблицу items
func InsertItem(tx *sqlx.Tx, items []model.Item) error {
	for _, item := range items {
		_, err := tx.NamedExec("INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, "+
			"total_price, nm_id, brand, status, order_uid) VALUES (:chrt_id, :track_number, :price, :rid, :name, :sale, :size, "+
			":total_price, :nm_id, :brand, :status, :order_uid)", &item)
		if err != nil {
			log.Printf("Error when adding data to the table items: %s", err)
			return err
		}
	}

	return nil
}

// InsertOrder - добавляет данные в таблицу orders
func InsertOrder(tx *sqlx.Tx, order model.Order) error {
	_, err := tx.NamedExec("INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, "+
		"locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) "+
		"VALUES (:order_uid, :track_number, :entry, :delivery_id, :payment_id, :locale, :internal_signature, "+
		":customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard)", &order)
	if err != nil {
		log.Printf("Error when adding data to the table orders: %s", err)
		return err
	}

	return nil
}

// SelectAllOrders - выгружает все данные из таблицы orders
func SelectAllOrders(database *sqlx.DB, orders *[]model.Order) error {
	err := database.Select(orders, "SELECT * FROM orders")

	return err
}

// GetDeliveryById - возвращает данные из таблицы deliveries с конкретным id
func GetDeliveryById(database *sqlx.DB, id int) (model.Delivery, error) {
	var delivery model.Delivery

	err := database.Get(&delivery, "SELECT * FROM deliveries WHERE id = $1", id)

	return delivery, err
}

// GetPaymentById - возвращает данные из таблицы payments с конкретным id
func GetPaymentById(database *sqlx.DB, id string) (model.Payment, error) {
	var payment model.Payment

	err := database.Get(&payment, "SELECT * FROM payments WHERE transaction = $1", id)

	return payment, err
}

// SelectItemsByOrderId - выбирает продукты (items) по id заказа (order_uid)
func SelectItemsByOrderId(database *sqlx.DB, orderId string) ([]model.Item, error) {
	var items []model.Item

	err := database.Select(&items, "SELECT * FROM items WHERE order_uid = $1", orderId)

	return items, err
}

// UploadOrdersData - добавляет данные из базы данных в слайс структур model.Order
func UploadOrdersData(database *sqlx.DB, orders *[]model.Order) {
	log.Println("Uploading orders data from database...")

	if err := SelectAllOrders(database, orders); err != nil {
		log.Fatalf("Error when select data from orders table: %s\n", err)
	}

	for i, order := range *orders {
		var err error

		(*orders)[i].Delivery, err = GetDeliveryById(database, order.DeliveryId)
		if err != nil {
			log.Fatalf("Error when get data from deliveries table: %s\n", err)
		}

		(*orders)[i].Payment, err = GetPaymentById(database, order.PaymentId)
		if err != nil {
			log.Fatalf("Error when get data from payments table: %s\n", err)
		}

		(*orders)[i].Items, err = SelectItemsByOrderId(database, order.OrderUid)
		if err != nil {
			log.Fatalf("Error when select data from items table: %s\n", err)
		}
	}

	log.Println("Data was successfully uploaded to the memory")
}
