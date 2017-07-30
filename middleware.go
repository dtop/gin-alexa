package ginalexa

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-alexa/alexa/parser"
	"github.com/go-alexa/alexa/response"
	"github.com/go-alexa/alexa/validations"
	"github.com/nicksnyder/go-i18n/i18n"
)

var MiddlewareLogInput bool = false
var MiddlewareLogOutput bool = true

// EchoMiddlewareAutomatic Acts as middleware and endpoint for your router definitions
func EchoMiddlewareAutomatic(app *EchoApplication) gin.HandlerFunc {

	return func(c *gin.Context) {

		app.Context.GinContext = c
		ec := app.Context

		var r *http.Request = c.Request
		validations.AppID = app.AppID

		// Verify certificate is good
		cert, err := validations.ValidateCertificate(r)
		if err != nil {
			log.Println("validation of certificate failed", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Verify signature is good
		body, err := validations.ValidateSignature(r, cert)
		if err != nil {
			log.Println("validation of signature failed", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if MiddlewareLogInput {
			log.Println(string(body))
		}

		var data json.RawMessage

		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println("unmarshalling of json rawmessage failed", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		req, err := parser.Parse(data)
		if err != nil {
			log.Println("parsing of json failed", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Make sure the request is good
		if err = validations.ValidateRequest(req); err != nil {
			log.Println("validation of request failed", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		T, err := i18n.Tfunc(req.Request.Locale, "en-US")
		if err != nil {
			log.Println("loading of translate failed", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		app.Context.T = T

		res := response.New()
		c.Set("echoRequest", req)

		if app.OnAuthCheck != nil {
			if err := app.OnAuthCheck(ec, req, res); err != nil {

				res.AddLinkAccountCard()

				if req.Request.Locale == "de-DE" {
					res.AddSSMLSpeech("<speak>Um my mailbox nutzen zu können, musst Du die Kontoverknüpfung durchführen. Öffne dazu Deine Alexa App und folge den Anweisungen.</speak>")
				} else {
					res.AddSSMLSpeech("<speak>To use my mailbox you are required to do an account link. Please open you Alexa App on your phone and follow the instructions.</speak>")
				}

				c.Header("Content-Type", "application/json;charset=UTF-8")
				c.JSON(200, res)
				return
			}
		}

		if app.Session != nil {

			sess, err := app.Session.New(req.Session.ID)
			if err != nil {
				log.Panicln(err)
			}

			app.Context.Session = sess
		}

		switch req.Request.Type {
		case "LaunchRequest":

			if app.OnLaunch != nil {
				app.OnLaunch(ec, req, res)
			}

			app.Context.Session.Store()

		case "IntentRequest":

			name := req.Request.Intent.Name
			proc, ok := app.intents[name]
			if !ok {

				c.AbortWithStatus(http.StatusInternalServerError)
				panic("unknown event " + name)
			}

			proc(ec, req, res)
			app.Context.Session.Store()

		case "SessionEndedRequest":

			if app.OnSessionEnded != nil {
				app.OnSessionEnded(ec, req, res)
				app.Context.Session.DeleteSession()
			}

		default:

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if MiddlewareLogOutput {

			raw, err := json.Marshal(res)
			if err != nil {
				log.Println(err)
			}

			log.Println(string(raw))
		}

		c.Header("Content-Type", "application/json;charset=UTF-8")
		c.JSON(200, res)

		c.Next()
	}
}

// EchoMiddleware delivers all things neccessary to check if the request was legit
func EchoMiddleware(AppID string) gin.HandlerFunc {

	return func(c *gin.Context) {

		var r *http.Request = c.Request
		validations.AppID = AppID

		// Verify certificate is good
		cert, err := validations.ValidateCertificate(r)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Verify signature is good
		body, err := validations.ValidateSignature(r, cert)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var data json.RawMessage

		err = json.Unmarshal(body, &data)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ev, err := parser.Parse(data)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Make sure the request is good
		if err = validations.ValidateRequest(ev); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Set("echoRequest", ev)
		c.Next()
	}
}
