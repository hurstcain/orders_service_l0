package model

import "errors"

type Delivery struct {
	Id      int    `json:"-" db:"id"`
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Zip     string `json:"zip" db:"zip"`
	City    string `json:"city" db:"city"`
	Address string `json:"address" db:"address"`
	Region  string `json:"region" db:"region"`
	Email   string `json:"email" db:"email"`
}

// CheckData - проверяет, не пусты ли обязательные элементы структуты
func (d *Delivery) CheckData() error {
	if d.Name == "" {
		return errors.New("delivery data is invalid. Name field is empty")
	}
	if d.Phone == "" {
		return errors.New("delivery data is invalid. Phone field is empty")
	}
	if d.Zip == "" {
		return errors.New("delivery data is invalid. Zip field is empty")
	}
	if d.City == "" {
		return errors.New("delivery data is invalid. City field is empty")
	}
	if d.Address == "" {
		return errors.New("delivery data is invalid. Address field is empty")
	}
	if d.Region == "" {
		return errors.New("delivery data is invalid. Region field is empty")
	}
	if d.Email == "" {
		return errors.New("delivery data is invalid. Email field is empty")
	}

	return nil
}

type Payment struct {
	Transaction  string `json:"transaction" db:"transaction"`
	RequestId    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency"`
	Provider     string `json:"provider" db:"provider"`
	Amount       int    `json:"amount" db:"amount"`
	PaymentDt    int    `json:"payment_dt" db:"payment_dt"`
	Bank         string `json:"bank" db:"bank"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee"`
}

func (p *Payment) CheckData() error {
	if p.Transaction == "" {
		return errors.New("payment data is invalid. Transaction field is empty")
	}
	if p.Currency == "" {
		return errors.New("payment data is invalid. Currency field is empty")
	}
	if p.Provider == "" {
		return errors.New("payment data is invalid. Provider field is empty")
	}
	if p.PaymentDt == 0 {
		return errors.New("payment data is invalid. PaymentDt field is empty")
	}
	if p.Bank == "" {
		return errors.New("payment data is invalid. Bank field is empty")
	}

	return nil
}

type Item struct {
	ChrtId      int    `json:"chrt_id" db:"chrt_id"`
	TrackNumber string `json:"track_number" db:"track_number"`
	Price       int    `json:"price" db:"price"`
	Rid         string `json:"rid" db:"rid"`
	Name        string `json:"name" db:"name"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price"`
	NmId        int    `json:"nm_id" db:"nm_id"`
	Brand       string `json:"brand" db:"brand"`
	Status      int    `json:"status" db:"status"`
	OrderUid    string `json:"-" db:"order_uid"`
}

func (i *Item) CheckData() error {
	if i.ChrtId == 0 {
		return errors.New("item data is invalid. ChrtId field is empty")
	}
	if i.TrackNumber == "" {
		return errors.New("item data is invalid. TrackNumber field is empty")
	}
	if i.Rid == "" {
		return errors.New("item data is invalid. Rid field is empty")
	}
	if i.Name == "" {
		return errors.New("item data is invalid. Name field is empty")
	}
	if i.NmId == 0 {
		return errors.New("item data is invalid. NmId field is empty")
	}
	if i.Brand == "" {
		return errors.New("item data is invalid. Brand field is empty")
	}
	if i.Status == 0 {
		return errors.New("item data is invalid. Status field is empty")
	}

	return nil
}

type Order struct {
	OrderUid          string   `json:"order_uid" db:"order_uid"`
	TrackNumber       string   `json:"track_number" db:"track_number"`
	Entry             string   `json:"entry" db:"entry"`
	DeliveryId        int      `json:"-" db:"delivery_id"`
	Delivery          Delivery `json:"delivery" db:"-"`
	PaymentId         string   `json:"-" db:"payment_id"`
	Payment           Payment  `json:"payment" db:"-"`
	Items             []Item   `json:"items" db:"-"`
	Locale            string   `json:"locale" db:"locale"`
	InternalSignature string   `json:"internal_signature" db:"internal_signature"`
	CustomerId        string   `json:"customer_id" db:"customer_id"`
	DeliveryService   string   `json:"delivery_service" db:"delivery_service"`
	Shardkey          string   `json:"shardkey" db:"shardkey"`
	SmId              int      `json:"sm_id" db:"sm_id"`
	DateCreated       string   `json:"date_created" db:"date_created"`
	OofShard          string   `json:"oof_shard" db:"oof_shard"`
}

func (o *Order) CheckData() error {
	if err := o.Delivery.CheckData(); err != nil {
		return err
	}
	if err := o.Payment.CheckData(); err != nil {
		return err
	}
	for _, item := range o.Items {
		if err := item.CheckData(); err != nil {
			return err
		}
	}
	if o.OrderUid == "" {
		return errors.New("order data is invalid. OrderUid field is empty")
	}
	if o.TrackNumber == "" {
		return errors.New("order data is invalid. TrackNumber field is empty")
	}
	if o.Entry == "" {
		return errors.New("order data is invalid. Entry field is empty")
	}
	if o.Locale == "" {
		return errors.New("order data is invalid. Locale field is empty")
	}
	if o.CustomerId == "" {
		return errors.New("order data is invalid. CustomerId field is empty")
	}
	if o.DeliveryService == "" {
		return errors.New("order data is invalid. DeliveryService field is empty")
	}
	if o.Shardkey == "" {
		return errors.New("order data is invalid. Shardkey field is empty")
	}
	if o.SmId == 0 {
		return errors.New("order data is invalid. SmId field is empty")
	}
	if o.DateCreated == "" {
		return errors.New("order data is invalid. DateCreated field is empty")
	}
	if o.OofShard == "" {
		return errors.New("order data is invalid. OofShard field is empty")
	}

	return nil
}
