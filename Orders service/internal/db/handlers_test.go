package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"orders_service/configs"
	"orders_service/internal/model"
	"testing"
)

func TestInsertDelivery(t *testing.T) {
	db, _ := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configs.TestDBUser, configs.TestDBPassword, configs.TestDBName))
	defer db.Close()
	delivery := model.Delivery{
		Name:    "testname",
		Phone:   "123456789",
		Zip:     "123",
		City:    "testcity",
		Address: "testaddress",
		Region:  "testregion",
		Email:   "testemail",
	}
	id, _ := InsertDelivery(db, delivery)
	var deliveryRes model.Delivery
	db.Get(&deliveryRes, "SELECT * FROM deliveries WHERE name = $1", delivery.Name)

	assert.NotEmpty(t, deliveryRes)
	assert.Equal(t, id, deliveryRes.Id)
}

func TestInsertPayment(t *testing.T) {
	db, _ := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configs.TestDBUser, configs.TestDBPassword, configs.TestDBName))
	defer db.Close()
	payment := model.Payment{Transaction: "test", RequestId: "test_requestid", Currency: "ru", Provider: "test_provider",
		Amount: 1, PaymentDt: 123, Bank: "test_banck", DeliveryCost: 1, GoodsTotal: 1, CustomFee: 1}
	tx := db.MustBegin()
	InsertPayment(tx, payment)
	tx.Commit()
	var paymentRes model.Payment
	db.Get(&paymentRes, "SELECT * FROM payments WHERE transaction = $1", payment.Transaction)

	assert.NotEmpty(t, paymentRes)
	assert.Equal(t, paymentRes.Bank, payment.Bank)
}

func TestInsertOrder(t *testing.T) {
	db, _ := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configs.TestDBUser, configs.TestDBPassword, configs.TestDBName))
	defer db.Close()
	order := model.Order{
		OrderUid:          "test",
		TrackNumber:       "test",
		Entry:             "test",
		DeliveryId:        1,
		PaymentId:         "test",
		Locale:            "ru",
		InternalSignature: "test",
		CustomerId:        "test",
		DeliveryService:   "test",
		Shardkey:          "test",
		SmId:              1,
		DateCreated:       "test",
		OofShard:          "test",
	}
	tx := db.MustBegin()
	InsertOrder(tx, order)
	tx.Commit()
	var orderRes model.Order
	db.Get(&orderRes, "SELECT * FROM orders WHERE order_uid = $1", order.OrderUid)

	assert.NotEmpty(t, orderRes)
	assert.Equal(t, orderRes.Locale, order.Locale)
}

func TestInsertItem(t *testing.T) {
	db, _ := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configs.TestDBUser, configs.TestDBPassword, configs.TestDBName))
	defer db.Close()
	item := make([]model.Item, 1)
	item[0] = model.Item{
		ChrtId:      1110,
		TrackNumber: "test",
		Price:       333,
		Rid:         "test_rid",
		Name:        "test",
		Sale:        0,
		Size:        "0",
		TotalPrice:  333,
		NmId:        22,
		Brand:       "test",
		Status:      202,
		OrderUid:    "testtest",
	}

	assert.Error(t, InsertItem(db.MustBegin(), item))

	item[0].OrderUid = "test"
	tx := db.MustBegin()
	InsertItem(tx, item)
	tx.Commit()
	var itemRes model.Item
	db.Get(&itemRes, "SELECT * FROM items WHERE rid = $1", item[0].Rid)

	assert.NotEmpty(t, itemRes)
	assert.Equal(t, itemRes.OrderUid, item[0].OrderUid)
}

func TestSelectAllOrders(t *testing.T) {
	db, _ := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configs.TestDBUser, configs.TestDBPassword, configs.TestDBName))
	defer db.Close()
	orders := make([]model.Order, 0)

	SelectAllOrders(db, &orders)

	assert.NotEmpty(t, orders)
	assert.Equal(t, len(orders), 1)
}

func TestGetDeliveryById(t *testing.T) {
	db, _ := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configs.TestDBUser, configs.TestDBPassword, configs.TestDBName))
	defer db.Close()
	delivery, _ := GetDeliveryById(db, 1)

	assert.NotEmpty(t, delivery)

	_, err := GetDeliveryById(db, 11)

	assert.Error(t, err)
}

func TestGetPaymentById(t *testing.T) {
	db, _ := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configs.TestDBUser, configs.TestDBPassword, configs.TestDBName))
	defer db.Close()
	payment, _ := GetPaymentById(db, "test")

	assert.NotEmpty(t, payment)

	_, err := GetPaymentById(db, "nottest")

	assert.Error(t, err)
}

func TestSelectItemsByOrderId(t *testing.T) {
	db, _ := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configs.TestDBUser, configs.TestDBPassword, configs.TestDBName))
	defer db.Close()
	items, _ := SelectItemsByOrderId(db, "test")

	assert.NotEmpty(t, items)
	assert.Equal(t, len(items), 1)

	items, _ = SelectItemsByOrderId(db, "nottest")

	assert.Empty(t, items)
}
