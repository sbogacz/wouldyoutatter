package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sbogacz/wouldyoutatter/service"
	"github.com/urfave/cli"
)

var (
	config         = &service.Config{}
	s              *service.Service
	useLocalDynamo bool
)

func flags() []cli.Flag {
	return append(config.Flags(),
		cli.BoolFlag{
			Name:        "local-dynamo",
			Usage:       "you can configure the app to run against local Dynamo either setting this flag, or setting the AWS_REGION env var to true",
			Destination: &useLocalDynamo,
		})
}

func main() {
	app := cli.NewApp()
	app.Usage = "this is the CLI app version of wouldyoutatter"
	app.Flags = flags()
	app.Action = serve

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func serve(c *cli.Context) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var err error
	s, err = service.New(*config)
	if err != nil {
		return err
	}
	go s.Start()

	<-sigs
	s.Stop()
	return nil
}
