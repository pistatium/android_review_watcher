package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/operando/golack"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"text/template"
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
	rate := args[0].(int64)
	return "★★★★★☆☆☆☆"[5-rate : 10-rate]
}

func formatReviews(reviews []*androidpublisher.Review) []Review {
	formatted := make([]Review, len(reviews))
	funcMap := template.FuncMap{
		"stars": Int2Stars,
	}
	for i, r := range reviews {
		buf := &bytes.Buffer{}
		tpl := template.Must(template.New("post.tpl").Funcs(funcMap).ParseFiles("templates/post.tpl"))
		if err := tpl.Execute(buf, r); err != nil {
			log.Print(err)
		}
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
		log.Print(app.PackageName)
		waitGroup.Add(1)
		go func(app TargetApp) {
			defer waitGroup.Done()
			review := getReview(service, app)
			formatted := formatReviews(review)
			if dry_run {
				for _, r := range formatted {
					fmt.Println(r)
				}
			} else {
				postSlack(formatted, app)
			}
		}(app)
	}
	waitGroup.Wait()
	return nil
}
