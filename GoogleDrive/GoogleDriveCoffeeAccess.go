package GoogleDrive

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// GetClient Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := TokenFromFile(tokFile)
	if err != nil {
		fmt.Println("Please enter your token after logging in to google account in the URL: ")
		tok = GetTokenFromWeb(config)
		SaveToken(tokFile, tok)
	}
	fmt.Println("Successfully connected to Google Drive API")
	return config.Client(context.Background(), tok)
}

// GetTokenFromWeb Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// TokenFromFile Retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("could not open " + f.Name() + "correctly!")
		}
	}(f)
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// SaveToken Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("something went wrong with " + f.Name())
		}
	}(f)
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("could not encode json")
	}
}

func SearchFilesOnDrive(client *http.Client, searchStrings []string, sharedDriveId string) []string {
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	var results []string

	for _, searchString := range searchStrings {
		query := fmt.Sprintf("name contains '%s' and mimeType = 'image/png'", searchString)
		r, err := srv.Files.List().
			Q(query).
			IncludeItemsFromAllDrives(true). // Include shared drives
			SupportsAllDrives(true).         // Support for shared drives
			DriveId(sharedDriveId).          // ID of the shared drive
			Corpora("drive").                // Search specifically in the provided drive
			Do()
		if err != nil {
			log.Printf("Unable to retrieve files: %v", err)
			continue
		}
		if len(r.Files) > 0 {
			for i := 0; i < len(r.Files); i++ {
				results = append(results, r.Files[i].Id+"#"+r.Files[i].Name)
			}
		}
	}
	return results
}

func DownloadFile(client *http.Client, fileID string, savePath string) error {
	srv, err := drive.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	resp, err := srv.Files.Get(fileID).Download()
	if err != nil {
		return fmt.Errorf("unable to download file with ID %s: %v", fileID, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("could not close the body")
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read data for file with ID %s: %v", fileID, err)
	}

	err = os.WriteFile(savePath, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to write file %s: %v", savePath, err)
	}

	return nil
}
