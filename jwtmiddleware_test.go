package jwtmiddleware_test

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noobj/jwtmiddleware"
	"github.com/noobj/jwtmiddleware/internal/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JwtAuth", func() {
	var fakeRequest events.APIGatewayV2HTTPRequest

	var fakeGetUserFromPayload = func(payload interface{}) (interface{}, error) {
		return payload, nil
	}

	BeforeEach(func() {
		fakeRequest = events.APIGatewayV2HTTPRequest{}

		os.Setenv("ACCESS_TOKEN_SECRET", "codeeatsleep")
	})

	Context("when use jwt auth as middleware before handler", func() {

		It("should contains user in context when passing eligible cookie", func() {
			fakeHandler := func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
				user := int(ctx.Value(helper.ContextKeyUser).(float64))
				Expect(user).To(Equal(1234))

				return events.APIGatewayProxyResponse{}, nil
			}
			fakeRequest.Cookies = []string{
				"access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjoxMjM0fQ.WXQUnBQL7ACuJX7oPnK5PVVhzY8hCEE8xn8VBP9W_Og",
			}
			jwtmiddleware.Handle(fakeHandler, fakeGetUserFromPayload)(context.Background(), fakeRequest)
		})

		It("should return 404 response when no cookie passed", func() {
			fakeHandler := func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
				return events.APIGatewayProxyResponse{}, nil
			}
			res, err := jwtmiddleware.Handle(fakeHandler, fakeGetUserFromPayload)(context.Background(), fakeRequest)

			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(401))
			Expect(res.Body).To(Equal("please login in"))
		})

	})

})
