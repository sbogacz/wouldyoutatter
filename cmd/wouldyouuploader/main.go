package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/sbogacz/wouldyoutatter/contender"
	"github.com/sbogacz/wouldyoutatter/service"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "this is the contender uploader for wouldyoutatter"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "svgpath",
			Usage: "Path to SVG file set",
			Value: "unset",
		},
		cli.StringFlag{
			Name:  "admintoken",
			Usage: "Admin Token",
			Value: service.DefaultMasterKey,
		},
		cli.StringFlag{
			Name:  "endpoint",
			Usage: "HTTP endpoint to call",
			Value: "http://localhost:8080/contenders",
		},
	}
	app.Action = upload

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func upload(c *cli.Context) error {
	if c.String("svgpath") == "unset" {
		// there's proably a better way to do required args
		return errors.New("svgpath is required")
	}
	contenders, err := loadContenders(c.String("svgpath"))
	if err != nil {
		return err
	}
	for _, contender := range *contenders {
		b, err := json.Marshal(&contender)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", c.String("endpoint"), bytes.NewBuffer(b))
		if err != nil {
			return err
		}
		req.Header.Set("X-Tatter-Master", c.String("admintoken"))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusCreated {
			log.Printf("Got back status of %d\n", resp.StatusCode)
			return errors.New("No status created code")
		}
	}
	return nil
}

func loadContenders(svgpath string) (*[]contender.Contender, error) {
	contenders := []contender.Contender{}
	f, err := os.Open(svgpath)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".svg") {
			continue
		}
		contenderName := strings.TrimSuffix(file.Name(), ".svg")
		fmt.Println(contenderName)
		content, err := ioutil.ReadFile(path.Join(svgpath, file.Name()))
		if err != nil {
			return nil, err
		}
		contenders = append(contenders, contender.Contender{
			Name:        contenderName,
			Description: contenderName,
			SVG:         content,
		})
	}
	return &contenders, nil
}
