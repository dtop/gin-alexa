# gin-alexa
Amazon Alexa module for gin

IMPORTANT:
most of the actual code is borrowed from [skillserver](https://github.com/mikeflynn/go-alexa) to use it with gin

Import:
```

go get github.com/dtop/gin-alexa

```

Example:

```go

func Routes(r *gin.Engine) {

    app1 := &AlexaApplication{
        AppID: "<YOUR APP ID>",
        OnIntent: YourOnIntentFunc,
        OnLaunch: YourOnLaunchFunc,
        OnSessionEnded: YourOnSessionEndedFunc,
    }

    alexa := r.Group("/alexa")
    {
        alexa.GET("/App1", AlexaMiddlewareAutomatic(app1))
    }
}

```


```go

func MyUsualGinHandlerFunc(c *gin.Context) {

    alexaRequest, ok := c.Get("alexaRequest")
    if ok {
    
        // ... handle
        
        alexaResponse := NewAlexaResponse()
        // ... respond
        
        c.Header("Content-Type", "application/json;charset=UTF-8")
        c.JSON(200, alexaResponse)
    }
}

func Routes(r *gin.Engine) {

    alexa := r.Group("/alexa")
    {
        alexa.GET("/App1", AlexaMiddleware("<YOUR APP ID>"), MyUsualGinHandlerFunc)
    }
}

```