package main

import (
	"log"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"google.golang.org/api/androidpublisher/v2"
	"fmt"
	"net/http"
	"sync"
)

var waitGroup sync.WaitGroup

func getReview(client *http.Client, appId string, result chan <- []*androidpublisher.Review ) {
	defer waitGroup.Done()
	res, err := androidpublisher.NewReviewsService(client).List(appId).Do()
	if err != nil {
		log.Fatalf("Unable to access review API: %v", err)
		result <- nil
		return
	}
	result <- res.Reviews
}

func main() {
	ctx := context.Background()

	appIds := []string {
		"com.appspot.pistatium.tomorrow",
		"com.appspot.pistatium.tenseconds",
	}

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, androidpublisher.AndroidpublisherScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)

	results := make(chan []*androidpublisher.Review, 2)

	for _, appId := range appIds {
		waitGroup.Add(1)
		go getReview(client, appId, results)
	}

	waitGroup.Wait()
	for _ := range appIds {
		reviews := <- results
		for review := range reviews {
			fmt.Print(review)
		}
	}
}
