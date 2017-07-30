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

func HandleIntent1(c *ginalexa.EchoContext, req *echoRequest.Event, res *echoResponse.Response) {

    // ... handle OnIntent
    // ... respond
}

func HandleIntent2(c *ginalexa.EchoContext, req *echoRequest.Event, res *echoResponse.Response) {

    // ... handle OnIntent
    // ... respond
    session := c.Session
    ginContext := c.GinContext
}

func Routes(r *gin.Engine) {

	echoApp := ginalexa.New(
		<YOUR APP ID>,
		nil,
		nil,
	)

	echoApp.Session = <your implementation of the session interface>

	echoApp.Set(
		ginalexa.MkEchoAction("", ginalexa.EventOnLaunch, endpoints.HandleSessionStart),
		ginalexa.MkEchoAction("", ginalexa.EventOnSessionEnded, endpoints.HandleSessionEnded),
		ginalexa.MkEchoAuthAction(endpoints.Auth),
		ginalexa.MkEchoAction("Intent1", ginalexa.EventOnIntent, endpoints.HandleIntent1),
		ginalexa.MkEchoAction("Intent2", ginalexa.EventOnIntent, endpoints.HandleIntent2),
	)

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