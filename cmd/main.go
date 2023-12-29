package main

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"fedimint-go-client/pkg/fedimint"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:5000"
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		password = "password"
	}

	fedimintClient := fedimint.NewFedimintClient(baseUrl, password)

	info, err := fedimintClient.Info()
	if err != nil {
		fmt.Println("Error getting info: ", err)
		return
	}
	fmt.Println("Current Total Msats Ecash: ", info.TotalAmountMsat)

	invoiceRequest := fedimint.LnInvoiceRequest{
		AmountMsat:  10000,
		Description: "test",
	}

	invoiceResponse, err := fedimintClient.Modules.Ln.CreateInvoice(invoiceRequest)
	if err != nil {
		fmt.Println("Error creating invoice: ", err)
		return
	}

	fmt.Println("Created 10 sat Invoice: ", invoiceResponse.Invoice)

	fmt.Println("Waiting for payment...")

	awaitInvoiceRequest := fedimint.AwaitInvoiceRequest{
		OperationID: invoiceResponse.OperationID,
	}

	_, err = fedimintClient.Modules.Ln.AwaitInvoice(awaitInvoiceRequest)
	if err != nil {
		fmt.Println("Error awaiting invoice: ", err)
		return
	}

	fmt.Println("Payment Received!")
	// fmt.Println("New Total Msats Ecash: ", awaitInvoiceResponse.TotalAmountMsat)
}
