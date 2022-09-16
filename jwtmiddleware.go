package jwtmiddleware

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
	"github.com/noobj/jwtmiddleware/internal/helper"
	"github.com/noobj/jwtmiddleware/internal/types"
)

func Handle[T types.ApiRequest, R types.ApiResponse](f types.HandlerFunc[T, R], getUserFromPayload func(interface{}) (interface{}, error)) types.HandlerFunc[T, R] {
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
		user, err := getUserFromPayload(payload)
		if err != nil {
			return helper.GenerateErrorResponse[R](401)
		}

		ctxWithUser := context.WithValue(ctx, helper.ContextKeyUser, user)

		return f(ctxWithUser, any(v2Request).(T))
	}
}
