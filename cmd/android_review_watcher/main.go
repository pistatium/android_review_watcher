package main

import (
	"os"
	"fmt"
	"sync"
	"log"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"google.golang.org/api/androidpublisher/v2"
	"github.com/codegangsta/cli"
	"github.com/operando/golack"
)

var waitGroup sync.WaitGroup

type AppReview struct {
	App     TargetApp
	Reviews []*androidpublisher.Review
}

func getReview(service *androidpublisher.Service, app TargetApp) AppReview {
	res, err := service.Reviews.List(app.PackageName).Do()
	if err != nil {
		log.Fatalf("Unable to access review API: ", err)
	}
	return AppReview{
		App:         app,
		Reviews:     res.Reviews,
	}
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

	appConfig, err := LoadConfig(config_file)
	if err != nil {
		log.Fatal("Unable to parse config file: ", err)
	}
	fmt.Printf("%v", appConfig)

	ctx := context.Background()
	b, err := ioutil.ReadFile(oauth_key)
	if err != nil {
		log.Fatal("Unable to read client secret file: ", err)
	}
	config, err := google.JWTConfigFromJSON(b, androidpublisher.AndroidpublisherScope)
	if err != nil {
		log.Fatal("Unable to parse client secret file to config: ", err)
	}
	client := config.Client(ctx)
	service, err := androidpublisher.New(client)
	if err != nil {
		log.Fatal("Unable to get service: ", err)
	}

	results := make(chan AppReview, 2)
	for _, app := range appConfig.TargetApps {
		log.Print(app.PackageName)
		waitGroup.Add(1)
		go func(app TargetApp) {
			defer waitGroup.Done()
			results <- getReview(service, app)
		}(app)
	}
	go func() {
		waitGroup.Wait()
		close(results)
	}()

	for result := range results {
		for _, review := range result.Reviews {
			fmt.Println(review.Comments[0].UserComment.StarRating)
			fmt.Println(review.Comments[0].UserComment.Text)
			payload := golack.Payload{
				Slack: result.App.SlackConf,
			}
			payload.Slack.Text = review.Comments[0].UserComment.Text
			// golack.Post(payload, appConfig.SlackWebHook)
		}
	}
	return nil
}
