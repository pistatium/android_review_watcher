package main

import (
	"log"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"google.golang.org/api/androidpublisher/v2"
	"fmt"
)

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, androidpublisher.AndroidpublisherScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := config.Client(ctx)
	appId := "com.appspot.pistatium.tomorrow"
	res, err := client.Get("https://www.googleapis.com/androidpublisher/v2/applications/" + appId + "/reviews")
	if err != nil {
		log.Fatalf("err %v", err)
	}
	fmt.Printf("success %v", res.Body)
}
