package android_review_watcher

import "io"

type App struct {
	PackageName string
	Writer      io.Writer
}
