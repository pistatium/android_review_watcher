package android_review_watcher

import (
	"bytes"
	"google.golang.org/api/androidpublisher/v2"
	"log"
	"text/template"
)

type Review []byte

func Int2Stars(rate int64) string {
	stars := []rune("★★★★★☆☆☆☆")
	return string(stars[5-rate : 10-rate])
}

func GetReview(service *androidpublisher.Service, app App) []*androidpublisher.Review {
	reviews, err := service.Reviews.List(app.PackageName).Do()
	if err != nil {
		log.Fatalf("Unable to access review API: ", err)
	}
	return reviews.Reviews
}

func FilterDuplicated(app App, reviews []*androidpublisher.Review) []*androidpublisher.Review {
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

func FormatReviews(reviews []*androidpublisher.Review) []Review {
	formatted := make([]Review, len(reviews))
	funcMap := template.FuncMap{
		"stars": Int2Stars,
	}
	tpl := template.Must(template.New("post.tpl").Funcs(funcMap).ParseFiles("templates/post.tpl"))

	for i, r := range reviews {
		buf := bytes.Buffer{}
		if err := tpl.Execute(&buf, r); err != nil {
			log.Print(err)
		}
		formatted[i] = Review(buf.Bytes())
	}
	return formatted
}
