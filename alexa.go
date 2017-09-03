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
	AmazonCancelIntent  = "AMAZON.CancelIntent"
	AmazonHelpIntent    = "AMAZON.HelpIntent"
	AmazonNextIntent    = "AMAZON.NextIntent"
	AmazonStopIntent    = "AMAZON.StopIntent"
)

type (

	// EchoError can be used to transport dispatchable errors through your app
	EchoError struct {
		errorString       string
		errorCode         int
		responseString    string
		responseI18nToken string
	}

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
		Session    Session
	}

	// EchoApplication is the actual application you define per endpoint
	EchoApplication struct {
		AppID          string
		config         Configurable
		Session        Session
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
			Session:    nil,
		},
	}
}

// Set sets echo actions
func (ea *EchoApplication) Set(actions ...EchoAction) {

	for _, v := range actions {

		switch v.GetType() {

		case EventOnLaunch:
			ea.OnLaunch = v.GetCallback()

		case EventOnSessionEnded:
			ea.OnSessionEnded = v.GetCallback()

		case EventOnIntent:
			ea.intents[v.GetName()] = v.GetCallback()

		case EventOnAuthCheck:
			eac := v.(echoAuthAct)
			ea.OnAuthCheck = eac.GetAuthCallback()
		}
	}
}

// ######################### EchoMethod

// ######################### EchoAction

func (ac echoAct) GetType() string         { return ac.eaType }
func (ac echoAct) GetName() string         { return ac.eaName }
func (ac echoAct) GetCallback() EchoMethod { return ac.eaCallback }

func (ac echoAuthAct) GetType() string         { return ac.eaType }
func (ac echoAuthAct) GetName() string         { return "" }
func (ac echoAuthAct) GetCallback() EchoMethod { return nil }

func (ac echoAuthAct) GetAuthCallback() func(*EchoContext, *echoRequest.Event, *echoResponse.Response) error {
	return ac.eaAuthCallback
}

// ######################### HELPERS

func MkCustomIntent(theName, theType string, theCallback EchoMethod) EchoAction {

	return echoAct{
		eaName:     theName,
		eaType:     theType,
		eaCallback: theCallback,
	}
}

func MkEchoAuthAction(theCallback func(*EchoContext, *echoRequest.Event, *echoResponse.Response) error) EchoAction {

	return echoAuthAct{
		eaType:         EventOnAuthCheck,
		eaAuthCallback: theCallback,
	}
}

func MkCancelIntent(theCallback EchoMethod) EchoAction {

	return echoAct{
		eaName:     AmazonCancelIntent,
		eaType:     EventOnIntent,
		eaCallback: theCallback,
	}
}

func MkHelpIntent(theCallback EchoMethod) EchoAction {

	return echoAct{
		eaName:     AmazonHelpIntent,
		eaType:     EventOnIntent,
		eaCallback: theCallback,
	}
}

func MkNextIntent(theCallback EchoMethod) EchoAction {

	return echoAct{
		eaName:     AmazonNextIntent,
		eaType:     EventOnIntent,
		eaCallback: theCallback,
	}
}

func MkStopIntent(theCallback EchoMethod) EchoAction {

	return echoAct{
		eaName:     AmazonStopIntent,
		eaType:     EventOnIntent,
		eaCallback: theCallback,
	}
}

// ################### EchoError

func NewEchoError(msg string, code int, responseMessage, i18nToken string) error {

	return EchoError{
		errorString:       msg,
		errorCode:         code,
		responseString:    responseMessage,
		responseI18nToken: i18nToken,
	}
}

func EchoErrorFromError(err error, code int, responseMessage, i18nToken string) error {

	return NewEchoError(err.Error(), code, responseMessage, i18nToken)
}

func (ee EchoError) Error() string {

	return ee.errorString
}

func (ee EchoError) Code() int {
	return ee.errorCode
}

func (ee EchoError) ResponseMessage() string {
	return ee.responseString
}

func (ee EchoError) ResponseI18nToken() string {
	return ee.responseI18nToken
}
