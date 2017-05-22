package ginalexa

import (
	"github.com/gin-gonic/gin"
	echoRequest "github.com/go-alexa/alexa/parser"
	echoResponse "github.com/go-alexa/alexa/response"
)

type EchoContext struct {
	AppID      string
	GinContext *gin.Context
	EchoConfig Configurable
}

// EchoApplication is the actual application you define per endpoint
type EchoApplication struct {
	AppID          string
	Config		   Configurable
	OnLaunch       func(*EchoContext, *echoRequest.Event, *echoResponse.Response)
	OnIntent       func(*EchoContext, *echoRequest.Event, *echoResponse.Response)
	OnSessionEnded func(*EchoContext, *echoRequest.Event, *echoResponse.Response)
	OnAuthCheck    func(*EchoContext, *echoRequest.Event, *echoResponse.Response) error
	Context        *EchoContext
}

func (ea *EchoApplication) Inject() *EchoApplication {

	if ea.AppID == "" {

		appId := ea.Config.GetString("AppID")
		if appId == "" {
			panic("appid is missing")
		}

		ea.AppID = appId
	}

	return ea
}

func (ea *EchoApplication) Init(c *gin.Context) {

	ea.Context = &EchoContext{
		AppID: ea.AppID,
		GinContext: c,
		EchoConfig: ea.Config,
	}
}