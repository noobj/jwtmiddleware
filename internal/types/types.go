package types

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

type ApiResponse interface {
	events.APIGatewayProxyResponse | events.APIGatewayV2HTTPResponse
}

type ApiRequest interface {
	events.APIGatewayProxyRequest | events.APIGatewayV2HTTPRequest
}

type JwtToken struct {
	Token     string
	ExpiresIn int
}

type MyCustomClaims struct {
	Payload interface{} `json:"payload"`
	jwt.StandardClaims
}

type HandlerFunc[T ApiRequest, R ApiResponse] func(context.Context, T) (R, error)
