package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sbogacz/wouldyoutatter/service"
	"github.com/urfave/cli"
)

var (
	config = &service.Config{}
	s      *service.Service
)

// Handler takes an APIGW request, converts it to an HTTP Request, and sends it to the
// http server
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	httpReq, err := APIGWReqToHTTP(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusBadRequest}, err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: http.StatusBadRequest}, err
	}

	return HTTPRespToAPIGW(resp)
}

func main() {
	app := cli.NewApp()
	app.Usage = "this is the CLI app version of wouldyoutatter"
	app.Flags = config.Flags()
	app.Action = serve

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("failed to start service with err: %v", err)
	}
}

func serve(c *cli.Context) error {

	var err error
	s, err = service.New(*config)
	if err != nil {
		return err
	}
	go s.Start()
	lambda.Start(Handler)

	s.Stop()
	return nil
}

// APIGWReqToHTTP converts APIGatewayProxyRequests and translates them to stdlib
// http.Requests
func APIGWReqToHTTP(req events.APIGatewayProxyRequest) (*http.Request, error) {
	addr := fmt.Sprintf("http://127.0.0.1:%d%s", config.Port, req.Path)
	httpReq, err := http.NewRequest(req.HTTPMethod, addr, strings.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	return httpReq, nil
}

// HTTPRespToAPIGW converts http.Responses and translates them to APIGatewayProxyResponses
func HTTPRespToAPIGW(resp *http.Response) (events.APIGatewayProxyResponse, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if err = resp.Body.Close(); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	headers := make(map[string]string, len(resp.Header))
	for k := range resp.Header {
		headers[k] = resp.Header.Get(k)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       string(body),
	}, nil
}
