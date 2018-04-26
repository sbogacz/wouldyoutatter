package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sbogacz/wouldyoutatter/service"
	log "github.com/sirupsen/logrus"
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
	for _, f := range config.Flags() {
		f.Apply(flag.CommandLine)
	}

	var err error
	s, err = service.New(*config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on port: %d\n", config.Port)
	go s.Start()

	lambda.Start(Handler)
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
	q := httpReq.URL.Query()
	for k, v := range req.QueryStringParameters {
		q.Set(k, v)
	}
	httpReq.URL.RawQuery = q.Encode()

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
