package android_review_watcher

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v2"
	"io/ioutil"
)

func GetGoogleService(oauthKey string) (*androidpublisher.Service, error) {
	ctx := context.Background()
	b, err := ioutil.ReadFile(oauthKey)
	if err != nil {
		return nil, err
	}
	config, err := google.JWTConfigFromJSON(b, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}
	client := config.Client(ctx)
	service, err := androidpublisher.New(client)
	if err != nil {
		return nil, err
	}
	return service, nil
}
