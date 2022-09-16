package helper

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
	"github.com/noobj/lambda-jwt-middleware/internal/types"
)

func GenerateErrorResponse[T types.ApiResponse](statusCode int, messages ...string) (T, error) {
	messageResp := StatusCodeDefaultMsgMap[statusCode]
	if len(messages) != 0 {
		messageResp = strings.Join(messages, "")
	}

	var resType T
	var res any
	switch t := any(resType).(type) {
	case events.APIGatewayProxyResponse:
		t.Body = messageResp
		t.StatusCode = statusCode
		res = t
	case events.APIGatewayV2HTTPResponse:
		t.Body = messageResp
		t.StatusCode = statusCode
		res = t
	}

	return res.(T), nil
}

var StatusCodeDefaultMsgMap = map[int]string{
	401: "please login in",
	500: "internal error",
}

func ParseCookie(cookies []string) map[string]string {
	result := make(map[string]string)
	for _, cookie := range cookies {
		splitStrings := strings.SplitN(cookie, "=", 2)
		if len(splitStrings) != 2 {
			continue
		}

		result[splitStrings[0]] = splitStrings[1]
	}

	return result
}

func ExtractPayloadFromToken(key string, jwtToken string) (interface{}, error) {
	var claims types.MyCustomClaims
	token, err := jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		log.Printf("jwt parse error: %v", err)
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims.Payload, nil
}

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var ContextKeyUser = contextKey("user")
