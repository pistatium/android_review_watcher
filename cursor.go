package android_review_watcher

import (
	"io/ioutil"
	"os"
	"strconv"
)

const CURSOR_DIR = ".cursor"

type cursor struct {
	key string
}

func NewCursor(key string) cursor {
	return cursor{key: key}
}

func (c *cursor) Load() (int64, error) {
	path := c.filepath()

	// If file is not exist, return default value
	_, err := os.Stat(path)
	if err != nil {
		return 0, nil
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseInt(string(buf), 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (c *cursor) Save(value int64) error {
	path := c.filepath()

	return ioutil.WriteFile(path, []byte(strconv.FormatInt(value, 10)), os.ModePerm)
}

func (c *cursor) filepath() string {
	return CURSOR_DIR + "/" + c.key
}
