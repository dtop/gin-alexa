# gin-alexa
Amazon Alexa module for gin

Changes:

changed the actual alexa stuff from [skillserver](https://github.com/mikeflynn/go-alexa) to
[go-alexa](https://github.com/go-alexa/alex) which is also a dependency now.

Added glide as depsys

Import:
```

go get github.com/dtop/gin-alexa

```

Example:

```go

import (
    "github.com/gin-gonic/gin"
    echoRequest "github.com/go-alexa/alexa/parser"
    echoResponse "github.com/go-alexa/alexa/response"
)

func YourOnIntentFunc(c *gin.Context, req *echoRequest.Event, res *echoResponse.Response) {

    // ... handle OnIntent
    // ... respond
}

func Routes(r *gin.Engine) {

    app1 := &EchoApplication{
        AppID: "<YOUR APP ID>",
        OnIntent: YourOnIntentFunc,
        OnLaunch: YourOnLaunchFunc,
        OnSessionEnded: YourOnSessionEndedFunc,
    }

    alexa := r.Group("/echo")
    {
        alexa.GET("/App1", EchoMiddlewareAutomatic(app1))
    }
}

```


```go

func MyUsualGinHandlerFunc(c *gin.Context) {

    alexaRequest, ok := c.Get("alexaRequest")
    if ok {
    
        // ... handle
        
        echoResponse := response.New()
        // ... respond
        
        c.Header("Content-Type", "application/json;charset=UTF-8")
        c.JSON(200, echoResponse)
    }
}

func Routes(r *gin.Engine) {

    alexa := r.Group("/echo")
    {
        alexa.GET("/App1", EchoMiddleware("<YOUR APP ID>"), MyUsualGinHandlerFunc)
    }
}

```