package CoffeeLabel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Order struct {
	Orders []struct {
		OrderId          int         `json:"orderId"`
		OrderNumber      string      `json:"orderNumber"`
		OrderKey         string      `json:"orderKey"`
		OrderDate        string      `json:"orderDate"`
		CreateDate       string      `json:"createDate"`
		ModifyDate       string      `json:"modifyDate"`
		PaymentDate      string      `json:"paymentDate"`
		ShipByDate       interface{} `json:"shipByDate"`
		OrderStatus      string      `json:"orderStatus"`
		CustomerId       int         `json:"customerId"`
		CustomerUsername string      `json:"customerUsername"`
		CustomerEmail    string      `json:"customerEmail"`
		BillTo           struct {
			Name            string      `json:"name"`
			Company         interface{} `json:"company"`
			Street1         interface{} `json:"street1"`
			Street2         interface{} `json:"street2"`
			Street3         interface{} `json:"street3"`
			City            interface{} `json:"city"`
			State           interface{} `json:"state"`
			PostalCode      interface{} `json:"postalCode"`
			Country         interface{} `json:"country"`
			Phone           interface{} `json:"phone"`
			Residential     interface{} `json:"residential"`
			AddressVerified interface{} `json:"addressVerified"`
		} `json:"billTo"`
		ShipTo struct {
			Name            string      `json:"name"`
			Company         interface{} `json:"company"`
			Street1         string      `json:"street1"`
			Street2         string      `json:"street2"`
			Street3         interface{} `json:"street3"`
			City            string      `json:"city"`
			State           string      `json:"state"`
			PostalCode      string      `json:"postalCode"`
			Country         string      `json:"country"`
			Phone           interface{} `json:"phone"`
			Residential     bool        `json:"residential"`
			AddressVerified string      `json:"addressVerified"`
		} `json:"shipTo"`
		Items []struct {
			OrderItemId int    `json:"orderItemId"`
			LineItemKey string `json:"lineItemKey"`
			Sku         string `json:"sku"`
			Name        string `json:"name"`
			ImageUrl    string `json:"imageUrl"`
			Weight      struct {
				Value       float64 `json:"value"`
				Units       string  `json:"units"`
				WeightUnits int     `json:"WeightUnits"`
			} `json:"weight"`
			Quantity          int           `json:"quantity"`
			UnitPrice         float64       `json:"unitPrice"`
			TaxAmount         float64       `json:"taxAmount"`
			ShippingAmount    float64       `json:"shippingAmount"`
			WarehouseLocation interface{}   `json:"warehouseLocation"`
			Options           []interface{} `json:"options"`
			ProductId         int           `json:"productId"`
			FulfillmentSku    interface{}   `json:"fulfillmentSku"`
			Adjustment        bool          `json:"adjustment"`
			Upc               string        `json:"upc"`
			CreateDate        string        `json:"createDate"`
			ModifyDate        string        `json:"modifyDate"`
		} `json:"items"`
		OrderTotal               float64     `json:"orderTotal"`
		AmountPaid               float64     `json:"amountPaid"`
		TaxAmount                float64     `json:"taxAmount"`
		ShippingAmount           float64     `json:"shippingAmount"`
		CustomerNotes            interface{} `json:"customerNotes"`
		InternalNotes            interface{} `json:"internalNotes"`
		Gift                     bool        `json:"gift"`
		GiftMessage              interface{} `json:"giftMessage"`
		PaymentMethod            string      `json:"paymentMethod"`
		RequestedShippingService string      `json:"requestedShippingService"`
		CarrierCode              interface{} `json:"carrierCode"`
		ServiceCode              interface{} `json:"serviceCode"`
		PackageCode              interface{} `json:"packageCode"`
		Confirmation             string      `json:"confirmation"`
		ShipDate                 interface{} `json:"shipDate"`
		HoldUntilDate            interface{} `json:"holdUntilDate"`
		Weight                   struct {
			Value       float64 `json:"value"`
			Units       string  `json:"units"`
			WeightUnits int     `json:"WeightUnits"`
		} `json:"weight"`
		Dimensions       interface{} `json:"dimensions"`
		InsuranceOptions struct {
			Provider       interface{} `json:"provider"`
			InsureShipment bool        `json:"insureShipment"`
			InsuredValue   float64     `json:"insuredValue"`
		} `json:"insuranceOptions"`
		InternationalOptions struct {
			Contents     string      `json:"contents"`
			CustomsItems interface{} `json:"customsItems"`
			NonDelivery  string      `json:"nonDelivery"`
		} `json:"internationalOptions"`
		AdvancedOptions struct {
			WarehouseId          int           `json:"warehouseId"`
			NonMachinable        bool          `json:"nonMachinable"`
			SaturdayDelivery     bool          `json:"saturdayDelivery"`
			ContainsAlcohol      bool          `json:"containsAlcohol"`
			MergedOrSplit        bool          `json:"mergedOrSplit"`
			MergedIds            []interface{} `json:"mergedIds"`
			ParentId             interface{}   `json:"parentId"`
			StoreId              int           `json:"storeId"`
			CustomField1         interface{}   `json:"customField1"`
			CustomField2         interface{}   `json:"customField2"`
			CustomField3         interface{}   `json:"customField3"`
			Source               string        `json:"source"`
			BillToParty          interface{}   `json:"billToParty"`
			BillToAccount        interface{}   `json:"billToAccount"`
			BillToPostalCode     interface{}   `json:"billToPostalCode"`
			BillToCountryCode    interface{}   `json:"billToCountryCode"`
			BillToMyOtherAccount interface{}   `json:"billToMyOtherAccount"`
		} `json:"advancedOptions"`
		TagIds                    interface{} `json:"tagIds"`
		UserId                    interface{} `json:"userId"`
		ExternallyFulfilled       bool        `json:"externallyFulfilled"`
		ExternallyFulfilledBy     interface{} `json:"externallyFulfilledBy"`
		ExternallyFulfilledById   interface{} `json:"externallyFulfilledById"`
		ExternallyFulfilledByName interface{} `json:"externallyFulfilledByName"`
		LabelMessages             interface{} `json:"labelMessages"`
	} `json:"orders"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Pages int `json:"pages"`
}

func Connect(api string) (Order, error) {
	url := "https://ssapi.shipstation.com/orders"
	method := "GET"
	var orders Order

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Host", "ssapi.shipstation.com")
	req.Header.Add("Authorization", "Basic "+api)

	res, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Unable to close IO reader")
		}
	}(res.Body)
	body, _ := io.ReadAll(res.Body)
	// This code unpacks the json data from Shopify/orders into orders
	err = json.Unmarshal(body, &orders)
	if err != nil {
		log.Fatal("Unable to access ssapi.shipstation.com with api key: " + api)
	}
	fmt.Println("Successfully connected to Shopify API")
	if len(orders.Orders) == 0 {
		return orders, errors.New("no orders to process")
	}
	fmt.Println("Successfully pulled in Orders from Shopify")
	return orders, nil
}
