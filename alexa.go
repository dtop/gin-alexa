package gin_alexa

import (
    "github.com/gin-gonic/gin"
    echoRequest "github.com/go-alexa/alexa/parser"
    echoResponse "github.com/go-alexa/alexa/response"
)

type EchoApplication struct {
    AppID           string
    OnLaunch        func(*gin.Context, *echoRequest.Event, *echoResponse.Response)
    OnIntent        func(*gin.Context, *echoRequest.Event, *echoResponse.Response)
    OnSessionEnded  func(*gin.Context, *echoRequest.Event, *echoResponse.Response)
    OnAuthCheck     func(*gin.Context, *echoRequest.Event, *echoResponse.Response) error
}