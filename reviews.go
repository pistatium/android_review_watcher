package android_review_watcher

import (
	"bytes"
	"google.golang.org/api/androidpublisher/v2"
	"log"
	"text/template"
)

type Review string

func Int2Stars(rate int64) string {
	stars := []rune("★★★★★☆☆☆☆")
	return string(stars[5-rate : 10-rate])
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
		formatted[i] = Review(buf.String())
	}
	return formatted
}
