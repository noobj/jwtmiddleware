package jwtmiddleware

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
	"github.com/noobj/jwtmiddleware/internal/helper"
	"github.com/noobj/jwtmiddleware/types"
)

func Handle[T types.ApiRequest, R types.ApiResponse](f types.HandlerFunc[T, R], payloadHandler func(context.Context, interface{}) (context.Context, error)) types.HandlerFunc[T, R] {
	return func(ctx context.Context, r T) (R, error) {
		v2Request, ok := any(r).(events.APIGatewayV2HTTPRequest)
		if !ok {
			return helper.GenerateErrorResponse[R](401)
		}
		cookiesMap := helper.ParseCookie(v2Request.Cookies)
		if _, ok := cookiesMap["access_token"]; !ok {
			return helper.GenerateErrorResponse[R](401)
		}

		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found", err)
		}
		key := os.Getenv("ACCESS_TOKEN_SECRET")
		payload, err := helper.ExtractPayloadFromToken(key, cookiesMap["access_token"])
		if err != nil {
			return helper.GenerateErrorResponse[R](401)
		}
		ctx, err = payloadHandler(ctx, payload)
		if err != nil {
			return helper.GenerateErrorResponse[R](401)
		}

		return f(ctx, any(v2Request).(T))
	}
}
