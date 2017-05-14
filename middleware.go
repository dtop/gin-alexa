package ginalexa

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-alexa/alexa/parser"
	"github.com/go-alexa/alexa/response"
	"github.com/go-alexa/alexa/validations"
	"log"
	"net/http"
)

// EchoMiddlewareAutomatic Acts as middleware and endpoint for your router definitions
func EchoMiddlewareAutomatic(app *EchoApplication) gin.HandlerFunc {

	return func(c *gin.Context) {

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

		log.Println(string(body))
		var data json.RawMessage

		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println("unmarshalling of json rawmessage failed", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ev, err := parser.Parse(data)
		if err != nil {
			log.Println("parsing of json failed", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Make sure the request is good
		if err = validations.ValidateRequest(ev); err != nil {
			log.Println("validation of request failed", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Set("echoRequest", ev)
		resp := response.New()

		if app.OnAuthCheck != nil {
			if err := app.OnAuthCheck(c, ev, resp); err != nil {

				resp.AddLinkAccountCard()
				c.Header("Content-Type", "application/json;charset=UTF-8")
				c.JSON(http.StatusUnauthorized, resp)
				return
			}
		}

		switch ev.Request.Type {
		case "LaunchRequest":
			if app.OnLaunch != nil {
				app.OnLaunch(c, ev, resp)
			}
		case "IntentRequest":
			if app.OnIntent != nil {
				app.OnIntent(c, ev, resp)
			}
		case "SessionEndedRequest":
			if app.OnSessionEnded != nil {
				app.OnSessionEnded(c, ev, resp)
			}
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "application/json;charset=UTF-8")
		c.JSON(200, resp)

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
