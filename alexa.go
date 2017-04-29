package gin_alexa

import "github.com/gin-gonic/gin"

type AlexaApplication struct {
    AppID           string
    OnLaunch        func(*gin.Context, *AlexaRequest, *AlexaResponse)
    OnIntent        func(*gin.Context, *AlexaRequest, *AlexaResponse)
    OnSessionEnded  func(*gin.Context, *AlexaRequest, *AlexaResponse)
    OnAuthCheck     func(*gin.Context, *AlexaRequest, *AlexaResponse) error
}