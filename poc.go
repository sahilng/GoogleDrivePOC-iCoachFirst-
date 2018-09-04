//This file is no longer being edited

package main

import (
        "encoding/json"
        "fmt"
        "io"
        "io/ioutil"
	"log"
        "net/http"
        "os"

        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
        tokenFile := "token.json"
        tok, err := tokenFromFile(tokenFile)
        if err != nil {
                tok = getTokenFromWeb(config)
                saveToken(tokenFile, tok)
        }
        return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code %v", err)
        }

        tok, err := config.Exchange(oauth2.NoContext, authCode)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web %v", err)
        }
        return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
        f, err := os.Open(file)
        defer f.Close()
        if err != nil {
                return nil, err
        }
        tok := &oauth2.Token{}
        err = json.NewDecoder(f).Decode(tok)
        return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
        fmt.Printf("Saving credential file to: %s\n", path)
        f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
        defer f.Close()
        if err != nil {
                log.Fatalf("Unable to cache oauth token: %v", err)
        }
        json.NewEncoder(f).Encode(token)
}



//only touch this function = above is OAuth

func main() {

//Authorization
        b, err := ioutil.ReadFile("credentials.json")
        if err != nil {
                log.Fatalf("Unable to read client secret file: %v", err)
        }

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, drive.DriveScope)
        if err != nil {
                log.Fatalf("Unable to parse client secret file to config: %v", err)
        }
        client := getClient(config)

        srv, err := drive.New(client)
        if err != nil {
                log.Fatalf("Unable to retrieve Drive client: %v", err)
        }



//begin iCoachFirst stuff
fmt.Println("Welcome to the Google Drive iCoachFirst POC!")
fmt.Println()
fmt.Println("Sample 1: Retrieving Google Files")
	//sample (get files and list all)
	r, err := srv.Files.List().PageSize(10).
                Fields("nextPageToken, files(id, name)").Do()
        if err != nil {
                log.Fatalf("Unable to retrieve files: %v", err)
        }

        if len(r.Files) == 0 {
                fmt.Println("No files found.")
        } else {
	fmt.Println("\n" + "List of Google Files: ") 
	       for _, i := range r.Files {
                        fmt.Printf("%s (%s)\n", i.Name, i.Id)
                }
	fmt.Println()
        }

//Sahil added these

fmt.Println("Sample 2: Downloading Google Files")
	//sample (download Google Doc)

	fileId := "1p3G8Cty3WXkx1I17t29vMGOPh-1Im9_YsWpoQ06INeo" //replace with appropriate fileId
	mimeType := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	res, err := srv.Files.Export(
		fileId,
		mimeType,
	).Download()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer res.Body.Close()
  	out,_ := os.Create(fileId + ".docx")
  	defer out.Close()
  	io.Copy(out, res.Body)
	fmt.Println("\nSuccessfully downloaded " + fileId + "\n")

fmt.Println("Sample 3: Uploading Google Files")
	//sample (save text file)
	var newFile drive.File
	newFile.Name = "new file test"
	newFile.MimeType = "application/vnd.google-apps.document"
	
	response, error := srv.Files.Create(
		&newFile).
		Fields("id").
	Do()
	
	if error != nil {
		log.Fatalf("Error: %v", error)
	} else {
		fmt.Println("\nSuccessfully uploaded file with id = " + response.Id + "\n")
	}

//fmt.Println("Sample 4: Setting Google Tags")
	//sample (set a tag)
	//can  be   done   behind the scenes   or here using https://developers.google.com/drive/api/v3/file -> Add Custom File Properties
	//Set key-value pair "tag" -> array of stringvalues

fmt.Println("Sample A: Sharing Google Files")
	//sample (share file externally)
	//	fileId := chosen file id 
	var addedPermissions drive.Permission
	addedPermissions.Type = "user"
	addedPermissions.Role = "writer"
	addedPermissions.EmailAddress = "sahilng1997@gmail.com"	
	_, perr := srv.Permissions.Create(
		fileId,
		&addedPermissions,).
	Do()	
	if perr != nil {
		log.Fatalf("Error: %v", perr)
	} else {
		fmt.Println("Successfully changed the permissions of " + fileId)
	}
}

