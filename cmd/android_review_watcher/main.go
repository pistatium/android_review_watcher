package main

import (
	"fmt"
	"sync"
	"log"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"google.golang.org/api/androidpublisher/v2"
	"github.com/codegangsta/cli"
	"os"
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
	app := cli.NewApp()
	app.Name = "android_review_watcher"
	app.Usage = "Android Review watcher"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "oauth_key, o",
			Usage: "Google OAuth key file (JSON file)",
		},
		cli.StringFlag{
			Name:  "config_file, c",
			Usage: "Application setting file (TOML file)",
		},
	}

	app.Action = watchReview

	app.Run(os.Args)

}

func watchReview(c *cli.Context) error {

	oauth_key := c.GlobalString("oauth_key")
	config_file := c.GlobalString("config_file")

	var appConfig Config

	LoadConfig(config_file, &appConfig)

	ctx := context.Background()

	b, err := ioutil.ReadFile(oauth_key)
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
	fmt.Printf("%v", appConfig)
	for _, app := range appConfig.TargetApps {
		log.Print(app.PackageName)
		waitGroup.Add(1)
		go func(appId string) {
			defer waitGroup.Done()
			results <- getReview(service, appId)
		}(app.PackageName)
	}
	waitGroup.Wait()
	for range appConfig.TargetApps {
		reviews := <-results
		for _, review := range reviews {
			fmt.Println(review.Comments[0].UserComment.StarRating)
			fmt.Println(review.Comments[0].UserComment.Text)
		}
	}
	return nil
}