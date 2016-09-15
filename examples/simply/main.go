package main

import (
	"log"

	payments "github.com/xboston/go-flexidengi"
)

func main() {

	var (
		serviceID int
		secretKey string
	)

	flexi := payments.NewFlexi(serviceID, secretKey)

	flexi.
		SetOrderID(123).
		SetCustomerID("user#123").
		SetProductID(123)

	log.Println(flexi.Sign())

	log.Println(flexi.MakeForm())
}
