# Coffee Company Shipping Label Generator

This repository contains a set of Go programs that facilitate the generation of shipping labels for coffee orders. The programs connect to both Shipstation Orders and Google Drive PhotoDump via APIs to download and print bulk labels for coffee orders by searching picture names in format itemID:SKU which pulled off of shipstation pending orders. This allows all the pictures to be stored in a folder for bulk printing/shipping.

## Prerequisites

Before using the programs, make sure you have the following files in the same directory as the programs:

- **ShipstationAPI-KeyAccess.txt**: This file should contain your Shipstation API credentials, including the public key and private key portions. If the file is not found, it will be created, and you should fill in the appropriate credentials.

- **GoogleDriveCredentials.json**: This file should contain your OAuth JSON data for connecting to Google Drive. If the file is not found, it will be created, and you should provide the necessary OAuth credentials.

### Setting Up Google API Console

Go to Google Cloud Console
Create a new project.
Navigate to Library, search for Google Drive API, and enable it.
Go to Credentials and create an OAuth 2.0 Client ID.
Choose Desktop App as the application type.
Make sure you add the drive.DriveReadonlyScope and drive.DriveScope
Download the credentials in .json format.

Note: Must use a google account with access to that drive or one associated with shared drive organization.


### Setting Up Shipstation API

Open Shipstation
Click Settings
Account > API Settings
Copy the API Key and API Secret into the ShipstationAPI-KeyAccess.txt

## Program Descriptions

### Main Program (main.go)

The `main.go` program is the entry point of the application. It performs the following steps:

1. Reads Shipstation API credentials from the `ShipstationAPI-KeyAccess.txt` file.
2. Validates and extracts the API keys.
3. Reads Google Drive OAuth credentials from the `GoogleDriveCredentials.json` file.
4. Initializes the Google Drive client using the OAuth credentials.
5. Fetches coffee orders from Shipstation using the provided API keys.
6. Extracts order names and quantities.
7. Searches Google Drive for images related to the orders using order names.
8. Downloads the images and stores them in a `Pictures` directory.

### GoogleDrive Package (GoogleDrive/google_drive.go)

The `GoogleDrive` package contains functions related to interacting with Google Drive:

- **GetClient**: Retrieves a Google Drive API client using OAuth credentials.
- **GetTokenFromWeb**: Requests and retrieves an OAuth token from the web.
- **TokenFromFile**: Retrieves an OAuth token from a local file.
- **SaveToken**: Saves an OAuth token to a file.
- **SearchFilesOnDrive**: Searches Google Drive for image files based on specified search strings.
- **DownloadFile**: Downloads a file from Google Drive and saves it to a local path.

### CoffeeLabel Package (CoffeeLabel/coffee_label.go)

The `CoffeeLabel` package contains functions to connect to Shipstation and fetch coffee orders:

- **Connect**: Connects to the Shipstation API using provided credentials and fetches coffee orders.

## Go Installation

1. **Install Go**: If you haven't already, install the Go programming language by following the installation instructions on the [official Go website](https://golang.org/doc/install).

## Usage
- **Run the program from main for it to generate the necessary files**

- **Optionally can run the executable file ShippingLabelsCoffeeCompany.exe from any folder or path**

1. **Fill in Credentials**: Open the `ShipstationAPI-KeyAccess.txt` and `GoogleDriveCredentials.json` files in a text editor and provide the required credentials as follows:

   - **ShipstationAPI-KeyAccess.txt**: Replace the placeholders with your actual Shipstation API credentials. The file should look like this:

       ```
       Public Key: YOUR_PUBLIC_KEY
       Private Key: YOUR_PRIVATE_KEY
       ```

   - **GoogleDriveCredentials.json**: Replace the content of this file with your actual OAuth JSON data for connecting to Google Drive.
   
      ```
       {
          "installed": 
       {
          "client_id": "YOUR_CLIENT_ID",
          "project_id": "YOUR_PROJECT_ID",
          "auth_uri": "https://accounts.google.com/o/oauth2/auth",
          "token_uri": "https://oauth2.googleapis.com/token",
          "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
          "client_secret": "YOUR_CLIENT_SECRET",
          "redirect_uris": ["http://localhost"]
          }}
       ```

2. **Run the Program**: After filling in the credentials, run the `main.go` or the executable `ShippingLabelsCoffeeCompany.exe` program using a Go-compatible environment. This program will perform the following steps:

   - Connect to the Shipstation API using the provided credentials.
   - Connect to Google Drive using OAuth credentials from the JSON file.
   - Fetch coffee orders from Shipstation.
   - Search for image names related to the orders from shipstation on Google Drive by **storeID:SKU**.
     
**NOTE** : Required naming convention for the images. They must be in .png format and have a name that is a concatenation of the SKU and StoreId from the shipstation. The SKU can be found in Orders.Items.Sku and the StoreId can be found in Orders.AdvancedOptionsStoreId
   - Download the images and store them in a `Pictures` directory.

Please note that this is a basic overview of the programs' functionality. You might need to adapt and extend the code according to your specific requirements and the evolving APIs of Shipstation and Google Drive.

## Author

These programs were authored by Guy Kogan.

Feel free to modify, extend, and distribute these programs as needed. If you have any questions or need assistance, don't hesitate to reach out.
