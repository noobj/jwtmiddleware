# lambda-jwt-middleware ![Test](https://github.com/noobj/jwtmiddleware/workflows/Go/badge.svg)


A middleware function used for verifying the jwt token from cookies which is compitable with [aws-lambda-go](https://github.com/aws/aws-lambda-go)

The basic structure of the function is based on [this article](https://www.zachjohnsondev.com/posts/lambda-go-middleware/).

>Since it uses Generics feature, Go version 1.18 or greater is required.

## Usage

It uses env variable to read the cookie secret if you've encrypted your cookie.

.env
```bash
ACCESS_TOKEN_SECRET=yoursecret
```

Define your own payloadHandler, as the example below, I extract the userId and search the entity in MongoDB, then store the entity into the context.

```go
func payloadHandler(ctx context.Context, payload interface{}) (context.Context, error) {
    userId, ok := payload.(string)
    userObjId, _ := primitive.ObjectIDFromHex(userId)
    if !ok {
        log.Printf("wrong payload format: %v", payload)
        return nil, fmt.Errorf("wrong payload format")
    }

    userRepo := UserRepository.New()
    defer userRepo.Disconnect()()
    var user UserRepository.User

    err := userRepo.FindOne(context.TODO(), bson.M{"_id": userObjId}).Decode(&user)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    ctx = context.WithValue(ctx, helper.ContextKeyUser, user)

    return ctx, nil
}
```

Then you use the middleware in the lambda.Start() function with the handler and payloadHandler as the arguments.

```go
func main() {
    lambda.Start(jwtmiddleware.Handle(Handler, payloadHandler))
}
```

Or you can combine the middleware with the payloadHandler to form another function which makes it shorter.

```go
import (
    "github.com/noobj/jwtmiddleware"
    "github.com/noobj/jwtmiddleware/types"
)

func Auth[T types.ApiRequest, R types.ApiResponse](f types.HandlerFunc[T, R]) types.HandlerFunc[T, R] {
    return jwtmiddleware.Handle(f, payloadHandler)
}

func main() {
    lambda.Start(Auth(Handler))
}
```