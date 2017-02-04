package main

import (
	"github.com/codegangsta/cli"
	. "github.com/pistatium/android_review_watcher"

	"log"
	"os"
	"sync"
)

var waitGroup sync.WaitGroup

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

	apps, err := LoadApps(configFile)
	if err != nil {
		log.Fatal("Unable to parse config file: ", err)
	}

	service, err := GetGoogleService(oauthKey)
	if err != nil {
		log.Fatal("Unable to get google service:", err)
	}

	for _, app := range apps {
		waitGroup.Add(1)
		go func(app App) {
			defer waitGroup.Done()
			review := GetReview(service, app)
			review = FilterDuplicated(app, review)
			formatted := FormatReviews(review)

			for _, r := range formatted {
				log.Printf("%s", r)
				if dry_run {
					continue
				}
				app.Writer.Write(r)
			}
		}(app)
	}
	waitGroup.Wait()
	return nil
}
