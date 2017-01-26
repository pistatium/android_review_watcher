package main

import (
	"log"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"google.golang.org/api/androidpublisher/v2"
	"fmt"
	"sync"
)

var waitGroup sync.WaitGroup

func getReview(service *androidpublisher.Service, appId string) []*androidpublisher.Review {
	res, err := service.Reviews.List(appId).Do()
	if err != nil {
		log.Fatalf("Unable to access review API: %v", err)
		return nil
	}
	return res.Reviews
}

func main() {
	ctx := context.Background()

	appIds := []string{
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

	service, err := androidpublisher.New(client)
	if err != nil {
		log.Fatal("Unable to get service: %v", err)
	}

	results := make(chan []*androidpublisher.Review, 2)

	for _, appId := range appIds {
		log.Print(appId)
		waitGroup.Add(1)
		go func(appId string) {
			defer waitGroup.Done()
			results <- getReview(service, appId)
		}(appId)
	}
	waitGroup.Wait()
	for range appIds {
		reviews := <-results
		for _, review := range reviews {
			fmt.Println(review.Comments[0].UserComment.StarRating)
			fmt.Println(review.Comments[0].UserComment.Text)
		}
	}
}
