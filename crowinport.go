package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	crowi "github.com/crowi/go-crowi"
)

var appConf Config
var importRecords ArayFormat
var ConfigPath = "./config.toml"

type Config struct {
	URL    string
	APIKey string
}

func (c *Config) ReadConfig(path string) error {
	if !exists(path) {
		return errors.New("no file")
	}
	_, err := toml.DecodeFile(path, &c)
	return err
}

type OutputFormat struct {
	Category string `toml:"category"`
	Body     string `toml:"body"`
}
type ArayFormat struct {
	Records []OutputFormat `toml:"records"`
}

func (c *ArayFormat) ReadConfig(path string) error {
	if !exists(path) {
		return errors.New("no file")
	}
	_, err := toml.DecodeFile(path, &c)
	return err
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func main() {
	importRecords.ReadConfig("./result.toml")
	appConf.ReadConfig("./config.toml")
	config := crowi.Config{
		URL:   appConf.URL,
		Token: appConf.APIKey,
	}
	client, err := crowi.NewClient(config)
	if err != nil {
		panic(err)
	}

	// var (
	// 	path = fmt.Sprintf("/user/%s/go-crowi-test-%d", "USER", time.Now().UnixNano())
	// 	body = "# this is a sample\n\ntest"
	// )
	for pos, v := range importRecords.Records {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := client.Pages.Create(ctx, "/"+v.Category, v.Body)
		if err != nil {
			panic(err)
		}

		if !res.OK {
			log.Printf("[ERROR] %s", res.Error)
		} else {
			fmt.Println(v.Category)
		}
		if pos > 10 {
			break
		}
		time.Sleep(1 * time.Second) // 3秒休む
	}
}
