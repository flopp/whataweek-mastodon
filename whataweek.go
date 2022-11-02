package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	mastodon "github.com/mattn/go-mastodon"
)

var (
	//go:embed whataweek.jpg
	image []byte
)

func readConfig(fileName string) (mastodon.Config, error) {
	confBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return mastodon.Config{}, fmt.Errorf("reading %s: %s", fileName, err)
	}

	var conf mastodon.Config
	if err = json.Unmarshal(confBytes, &conf); err != nil {
		return mastodon.Config{}, fmt.Errorf("unmarshalling %s: %s", fileName, err)
	}

	return conf, nil
}

func main() {
	if len(os.Args) != 2 {
		panic("config file missing")
	}

	conf, err := readConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	client := mastodon.NewClient(&conf)
	ctx := context.Background()

	reader := bytes.NewReader(image)
	attachment, err := client.UploadMediaFromReader(ctx, reader)
	if err != nil {
		panic(err)
	}

	toot := mastodon.Toot{
		Status: "",
	}
	toot.MediaIDs = append(toot.MediaIDs, attachment.ID)

	status, err := client.PostStatus(ctx, &toot)
	if err != nil {
		panic(err)
	}
	fmt.Println(status.URL)
}
