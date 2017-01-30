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
	"time"
	"text/template"
	"bytes"
)

var waitGroup sync.WaitGroup

type Review string

func getReview(service *androidpublisher.Service, app TargetApp) []*androidpublisher.Review {
	reviews, err := service.Reviews.List(app.PackageName).Do()
	if err != nil {
		log.Fatalf("Unable to access review API: ", err)
	}
	return reviews.Reviews
}

func postSlack(reviews []Review, app TargetApp) {
	for _, r := range reviews {
		payload := golack.Payload{
			Slack: app.SlackConf,
		}
		payload.Slack.Text = string(r)
		// golack.Post(payload, appConfig.SlackWebHook)
	}
}

func Int2Stars(args ...interface{}) string {
	rate := args[0].(int)
	return "★★★★★☆☆☆☆"[5-rate:10-rate]
}

func formatReviews(reviews []*androidpublisher.Review, interval int) []Review {
	t := int64(time.Now().Add(time.Duration(-interval) * time.Minute).Second())
	formatted := make([]Review, len(reviews))
	for i, r := range reviews {
		if int64(r.Comments[0].UserComment.LastModified.Seconds) < t {
			continue
		}
		tpl := template.Must(template.ParseFiles("templates/post.tpl"))
		tpl.Funcs(template.FuncMap{
			"stars": Int2Stars,
		})
		buf := &bytes.Buffer{}
		tpl.Execute(buf, r)
		formatted[i] = Review(buf.String())
	}
	return formatted
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
		cli.IntFlag{
			Name: "duration",
			Usage: "Fetch duration of reviews. (minutes)",
			Value: 24 * 60,
		},
		cli.BoolFlag{
			Name: "dry_run",
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
	fmt.Printf("%v", appConfig)

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
		log.Print(app.PackageName)
		waitGroup.Add(1)
		go func(app TargetApp) {
			defer waitGroup.Done()
			review := getReview(service, app)
			formatted := formatReviews(review, 24)
			if dry_run {
				for _, r := range formatted {
					fmt.Print(r)
				}
			} else {
				postSlack(formatted, app)
			}
		}(app)
	}
	waitGroup.Wait()
	return nil
}
