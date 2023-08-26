// Author: Guy Kogan
//
//	Description:
//		This program connects to GoogleDrive PhotoDump
//		as well as Shipstation Orders via API to help
//		print bulk labels on orders for Coffee
//
// Required Files:
//
//	ShipstationApI-KeyAccess.txt : if file is not found it will be created
//		fill in private key and public key portion
//	GoogleDriveCredentials.json : if file is not found it will be created
//		Fill in with OAuth json data to connect google-drive
package main

import (
	"ShippingLabelsCoffeeCompany/CoffeeLabel"
	"ShippingLabelsCoffeeCompany/GoogleDrive"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed opening log file: %s", err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Could not write to log file")
		}
	}(file)
	// Set the log output to write to the file.
	log.SetOutput(file)

	credentialShipstation, err := os.ReadFile("ShipstationApI-KeyAccess.txt")
	if err != nil {
		create, err := os.Create("ShipstationApI-KeyAccess.txt")
		if err != nil {
			log.Fatalf("Could not create file")
		}
		_, err = create.Write([]byte("PUBLIC KEY: \nPRIVATE KEY: "))
		if err != nil {
			log.Fatalf("Could not write to file")
		}
		log.Println("ShipstationApI-KeyAccess.txt does not exist created file")
	}
	credentials := string(credentialShipstation)
	reg := regexp.MustCompile("[^ :A-Z\r\n]+")
	keys := reg.FindAllString(credentials, -1)
	b, err := os.ReadFile("GoogleDriveCredentials.json")
	if err != nil {
		_, err := os.Create("GoogleDriveCredentials.json")
		if err != nil {
			log.Fatalf("Could not create file")
		}
		log.Fatalf("GoogleDriveCredentials does not exist created file, Please store credentials" +
			" to googledrive OAuth photodump here as json")
	}
	if len(keys) != 2 {
		log.Fatalf("ShipstationApI-KeyAccess.txt keys improperly set")
	}
	data := []byte(keys[0] + ":" + keys[1])
	str := base64.StdEncoding.EncodeToString(data)
	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	//fmt.Println("---Please authenticate in browser with email that has access to the " +
	//"Google Drive folder of coffee labels---")
	client := GoogleDrive.GetClient(config)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	orders, err := CoffeeLabel.Connect(str)
	if err != nil {
		log.Fatalf("Could not pull in orders because orders are empty")
	}
	var orderSku []string
	quantity := 1
	for i := 0; i < len(orders.Orders); i++ {
		if orders.Orders[i].ShipDate != nil {
			continue
		}
		for j := 0; j < len(orders.Orders[i].Items); j++ {
			searchCoffee := orders.Orders[i].Items[j].Sku
			searchCoffee = strconv.Itoa(orders.Orders[i].AdvancedOptions.StoreId) + "#" + searchCoffee
			quantity = orders.Orders[i].Items[j].Quantity
			if quantity > 1 {
				for i := 0; i < quantity; i++ {
					orderSku = append(orderSku, searchCoffee)
				}
			} else {
				orderSku = append(orderSku, searchCoffee)
			}
		}
	}
	stringID := GoogleDrive.SearchFilesOnDrive(client, orderSku, "0ALwW2dhzNkP8Uk9PVA")
	_, err = os.Stat("Pictures")
	if err != nil {
		fmt.Println("Creating Directory: Pictures")
		err = os.Mkdir("Pictures", os.ModePerm)
		if err != nil {
			log.Fatalf("could not create folder Pictures")
		}
	}
	for i, id := range stringID {
		id = strings.Split(id, "#")[0]
		err = GoogleDrive.DownloadFile(client, id, ".\\Pictures\\"+strings.ReplaceAll(orderSku[i], ":", "")+"---"+strconv.Itoa(i+1)+".png")
		if err != nil {
			log.Println("\n-------\nCould not download file with name: " + orderSku[i] + "\n-------")
		}
	}
	fmt.Println("Successfully pulled in Pictures from Google Drive")
}
