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
	var params Parameters
	var pc calculator.PackCalculator
	var body []byte
	var buf bytes.Buffer

	// decode request parameters
	err := json.Unmarshal([]byte(req.Body), &params)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// decide which calculator to use based on available pack sizes
	if len(params.PackSizes) == 1 {
		pc = calculator.SimplePackCalculator{PackSize: params.PackSizes[0]}
	} else {
		pc = calculator.GraphPackCalculator{PackSizes: params.PackSizes}
	}

	packs, err := pc.Calculate(params.Quantity)

	// prepare the payload
	if err != nil {
		resp.StatusCode = 400
		body, err = json.Marshal(map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		body, err = json.Marshal(packs)
	}

	if err != nil {
		return Response{StatusCode: 500}, err
	}

	json.HTMLEscape(&buf, body)
	resp.Body = buf.String()

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
