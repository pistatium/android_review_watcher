package main

import (
	"github.com/codegangsta/cli"
	"github.com/operando/golack"
	. "github.com/pistatium/android_review_watcher"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var waitGroup sync.WaitGroup

func getReview(service *androidpublisher.Service, app TargetApp) []*androidpublisher.Review {
	reviews, err := service.Reviews.List(app.PackageName).Do()
	if err != nil {
		log.Fatalf("Unable to access review API: ", err)
	}
	return reviews.Reviews
}

func postSlack(review Review, app TargetApp, webhook golack.Webhook) {
	payload := golack.Payload{
		Slack: app.SlackConf,
	}
	payload.Slack.Text = string(review)
	golack.Post(payload, webhook)
}

func filterDuplicated(app TargetApp, reviews []*androidpublisher.Review) []*androidpublisher.Review {
	cursor := NewCursor(app.PackageName)
	c, err := cursor.Load()
	if err != nil {
		log.Fatal("Load cursor error: ", err)
	}
	var index int
	for i, r := range reviews {
		rts := r.Comments[0].UserComment.LastModified.Seconds
		if rts <= c {
			break
		}
		index = i + 1
	}
	if len(reviews) == 0 {
		return reviews
	}
	cursor.Save(reviews[0].Comments[0].UserComment.LastModified.Seconds)
	return reviews[:index]
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
			Value: "client_secret.json",
		},
		cli.StringFlag{
			Name:  "config_file, c",
			Usage: "Application setting file (TOML file)",
			Value: "config.toml",
		},
		cli.BoolFlag{
			Name:  "dry_run",
			Usage: "Get reviews only. (without posting to slack)",
		},
	}

	app.Action = watchReview

	app.Run(os.Args)

}

func watchReview(c *cli.Context) error {

	oauthKey := c.GlobalString("oauth_key")
	configFile := c.GlobalString("config_file")
	dry_run := c.GlobalBool("dry_run")

	appConfig, err := LoadConfig(configFile)
	if err != nil {
		log.Fatal("Unable to parse config file: ", err)
	}
	log.Printf("%v", appConfig)

	ctx := context.Background()
	b, err := ioutil.ReadFile(oauthKey)
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

	for _, app := range appConfig.TargetApps {
		waitGroup.Add(1)
		go func(app TargetApp) {
			defer waitGroup.Done()
			review := getReview(service, app)
			review = filterDuplicated(app, review)
			formatted := FormatReviews(review)

			for _, r := range formatted {
				log.Println(r)
				if dry_run {
					continue
				}
				postSlack(r, app, appConfig.SlackWebHook)
			}
		}(app)
	}
	waitGroup.Wait()
	return nil
}
