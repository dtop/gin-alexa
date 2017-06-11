package ginalexa

import (
	"github.com/gin-gonic/gin"
	echoRequest "github.com/go-alexa/alexa/parser"
	echoResponse "github.com/go-alexa/alexa/response"
	"github.com/nicksnyder/go-i18n/i18n"
)

const (
	EventOnLaunch       = "OnLaunch"
	EventOnIntent       = "OnIntent"
	EventOnSessionEnded = "OnSessionEnded"
	EventOnAuthCheck    = "OnAuthCheck"
)

type (

	// EchoMethod is a shortform for the actual callback for all methods
	EchoMethod func(*EchoContext, *echoRequest.Event, *echoResponse.Response)

	// EchoAction represents the actions which can be registered to the app
	EchoAction interface {
		GetType() string
		GetName() string
		GetCallback() EchoMethod
	}

	// EchoContext provides a unified set of information to each callup
	EchoContext struct {
		AppID      string
		GinContext *gin.Context
		EchoConfig Configurable
		T          i18n.TranslateFunc
	}

	// EchoApplication is the actual application you define per endpoint
	EchoApplication struct {
		AppID          string
		config         Configurable
		OnLaunch       EchoMethod
		OnSessionEnded EchoMethod
		OnAuthCheck    func(*EchoContext, *echoRequest.Event, *echoResponse.Response) error
		Context        *EchoContext

		intents map[string]EchoMethod
	}

	echoAct struct {
		eaName     string
		eaType     string
		eaCallback EchoMethod
	}

	echoAuthAct struct {
		echoAct

		eaName         string
		eaType         string
		eaAuthCallback func(*EchoContext, *echoRequest.Event, *echoResponse.Response) error
	}
)

// New creates a new representation of EchoApplication
func New(appID string, config Configurable, c *gin.Context) *EchoApplication {

	_appid := appID
	if _appid == "" {
		_appid = config.GetString("AppID")
	}

	if _appid == "" {
		panic("no appID given")
		return nil
	}

	return &EchoApplication{
		AppID:   _appid,
		config:  config,
		intents: make(map[string]EchoMethod),
		Context: &EchoContext{
			AppID:      appID,
			GinContext: c,
			EchoConfig: config,
		},
	}
}

// Set sets echo actions
func (ea *EchoApplication) Set(actions ...EchoAction) {

	for _, v := range actions {

		switch v.GetType() {

		case EventOnLaunch:
			ea.OnLaunch = v.GetCallback()
			break

		case EventOnSessionEnded:
			ea.OnSessionEnded = v.GetCallback()
			break

		case EventOnIntent:
			ea.intents[v.GetName()] = v.GetCallback()
			break

		case EventOnAuthCheck:
			eac := v.(*echoAuthAct)
			ea.OnAuthCheck = eac.GetAuthCallback()
		}
	}
}

// ######################### EchoMethod

// ######################### EchoAction

func (ac echoAct) GetType() string         { return ac.eaType }
func (ac echoAct) GetName() string         { return ac.eaName }
func (ac echoAct) GetCallback() EchoMethod { return ac.eaCallback }

func (ac echoAuthAct) GetAuthCallback() func(*EchoContext, *echoRequest.Event, *echoResponse.Response) error {
	return ac.eaAuthCallback
}

// ######################### HELPERS

func MkEchoAction(theName, theType string, theCallback EchoMethod) EchoAction {

	return echoAct{
		eaName:     theName,
		eaType:     theType,
		eaCallback: theCallback,
	}
}

func MkEchoAuthAction(theCallback func(*EchoContext, *echoRequest.Event, *echoResponse.Response) error) EchoAction {

	return echoAuthAct{
		eaName:         "",
		eaType:         EventOnAuthCheck,
		eaAuthCallback: theCallback,
	}
}
