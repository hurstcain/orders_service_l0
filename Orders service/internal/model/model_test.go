package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelivery_CheckData(t *testing.T) {
	validDelivery := Delivery{
		Name:    "testname",
		Phone:   "123456789",
		Zip:     "123",
		City:    "testcity",
		Address: "testaddress",
		Region:  "testregion",
		Email:   "testemail",
	}
	invalidDelivery := Delivery{
		Phone:   "123456789",
		Zip:     "",
		City:    "testcity",
		Address: "testaddress",
		Region:  "testregion",
		Email:   "",
	}

	assert.NoError(t, validDelivery.CheckData())
	assert.Error(t, invalidDelivery.CheckData())
}

func TestPayment_CheckData(t *testing.T) {
	validPayment := Payment{
		Transaction:  "test",
		RequestId:    "test_requestid",
		Currency:     "ru",
		Provider:     "test_provider",
		Amount:       1,
		PaymentDt:    123,
		Bank:         "test_banck",
		DeliveryCost: 1,
		GoodsTotal:   1,
		CustomFee:    1,
	}
	invalidPayment := Payment{
		Transaction:  "test",
		RequestId:    "",
		Currency:     "ru",
		Provider:     "",
		PaymentDt:    123,
		Bank:         "test_banck",
		DeliveryCost: 1,
		CustomFee:    1,
	}

	assert.NoError(t, validPayment.CheckData())
	assert.Error(t, invalidPayment.CheckData())
}

func TestItem_CheckData(t *testing.T) {
	validItem := Item{
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
	invalidItem := Item{
		ChrtId:      0,
		TrackNumber: "test",
		Rid:         "test_rid",
		Name:        "",
		Size:        "0",
		TotalPrice:  333,
		NmId:        22,
		Brand:       "test",
		Status:      0,
	}

	assert.NoError(t, validItem.CheckData())
	assert.Error(t, invalidItem.CheckData())
}

func TestOrder_CheckData(t *testing.T) {
	validOrder := Order{
		OrderUid:    "test",
		TrackNumber: "test",
		Entry:       "test",
		Delivery: Delivery{
			Name:    "testname",
			Phone:   "123456789",
			Zip:     "123",
			City:    "testcity",
			Address: "testaddress",
			Region:  "testregion",
			Email:   "testemail",
		},
		Payment: Payment{
			Transaction:  "test",
			RequestId:    "test_requestid",
			Currency:     "ru",
			Provider:     "test_provider",
			Amount:       1,
			PaymentDt:    123,
			Bank:         "test_banck",
			DeliveryCost: 1,
			GoodsTotal:   1,
			CustomFee:    1,
		},
		Items: []Item{
			{
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
			},
		},
		Locale:            "ru",
		InternalSignature: "test",
		CustomerId:        "test",
		DeliveryService:   "test",
		Shardkey:          "test",
		SmId:              10,
		DateCreated:       "test",
		OofShard:          "test",
	}
	invalidOrder := Order{
		OrderUid:    "test",
		TrackNumber: "",
		Delivery: Delivery{
			Name:    "testname",
			Phone:   "123456789",
			Zip:     "123",
			City:    "testcity",
			Address: "testaddress",
			Region:  "testregion",
			Email:   "testemail",
		},
		Payment: Payment{
			Transaction:  "test",
			RequestId:    "test_requestid",
			Currency:     "ru",
			Provider:     "test_provider",
			Amount:       1,
			PaymentDt:    123,
			Bank:         "test_banck",
			DeliveryCost: 1,
			GoodsTotal:   1,
			CustomFee:    1,
		},
		Items: []Item{
			{
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
			},
		},
		Locale:            "ru",
		InternalSignature: "test",
		DateCreated:       "test",
		OofShard:          "",
	}

	assert.NoError(t, validOrder.CheckData())
	assert.Error(t, invalidOrder.CheckData())
}
