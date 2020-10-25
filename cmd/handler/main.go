package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lushc/pack-calc-go/pkg/calculator"
)

// Request is of type APIGatewayProxyRequest since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
type Request events.APIGatewayProxyRequest

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
type Response events.APIGatewayProxyResponse

// Parameters are those sent in the Request body
type Parameters struct {
	Quantity  int
	PackSizes []int
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, req Request) (Response, error) {
	// decode request parameters
	var params Parameters
	err := json.Unmarshal([]byte(req.Body), &params)
	if err != nil {
		return Response{StatusCode: 400}, err
	}

	// decide which calculator to use based on available pack sizes
	var packs calculator.PackCalculator
	if len(params.PackSizes) == 1 {
		packs = calculator.SimplePackCalculator{PackSize: params.PackSizes[0]}
	} else {
		packs = calculator.GraphPackCalculator{PackSizes: params.PackSizes}
	}

	// TODO: validate parameters
	response := packs.Calculate(params.Quantity)
	body, err := json.Marshal(response)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
