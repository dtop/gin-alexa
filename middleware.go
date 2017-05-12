package gin_alexa

import (
    "github.com/gin-gonic/gin"
    "encoding/json"
    "github.com/go-alexa/alexa/validations"
    "net/http"
    "github.com/go-alexa/alexa/parser"
    "github.com/go-alexa/alexa/response"
    "log"
)

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
                c.AbortWithStatus(http.StatusUnauthorized)
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



//func AlexaMiddlewareAutomatic(app *AlexaApplication) gin.HandlerFunc {
//
//    return func (c *gin.Context) {
//
//        if err := verifyJSON(c, app.AppID); err != nil {
//
//            //msg := err.Error()
//            code := err.(*HttpError).Code()
//            expl := err.(*HttpError).Explination()
//
//            c.AbortWithError(code, err);
//            log.Println(expl)
//            return
//        }
//
//        if err := validateRequest(c); err != nil {
//
//            //msg := err.Error()
//            code := err.(*HttpError).Code()
//            expl := err.(*HttpError).Explination()
//
//            c.AbortWithError(code, err);
//            log.Println(expl)
//            return
//        }
//
//        // use app
//        req, ok := c.Get("alexaRequest")
//        if ok {
//            echoReq  := req.(*AlexaRequest)
//            echoResp := NewAlexaResponse()
//
//            if app.OnAuthCheck != nil {
//                if err := app.OnAuthCheck(c, echoReq, echoResp); err != nil {
//                    c.AbortWithStatus(403)
//                    return
//                }
//            }
//
//            switch echoReq.GetRequestType() {
//
//            case "LaunchRequest":
//                if app.OnLaunch != nil {
//                    app.OnLaunch(c, echoReq, echoResp)
//                }
//            case "IntentRequest":
//                if app.OnIntent != nil {
//                    app.OnIntent(c, echoReq, echoResp)
//                }
//            case "SessionEndedRequest":
//                if app.OnSessionEnded != nil {
//                    app.OnSessionEnded(c, echoReq, echoResp)
//                }
//            default:
//                c.AbortWithStatus(500)
//            }
//
//            c.Header("Content-Type", "application/json;charset=UTF-8")
//            c.JSON(200, echoResp)
//            return
//        }
//
//        c.Next()
//    }
//}
//
//func AlexaMiddleware(AppID string) gin.HandlerFunc {
//
//
//    return func (c *gin.Context) {
//
//        if err := verifyJSON(c, AppID); err != nil {
//
//            //msg := err.Error()
//            code := err.(*HttpError).Code()
//            expl := err.(*HttpError).Explination()
//
//            c.AbortWithError(code, err);
//            log.Println(expl)
//            return
//        }
//
//        if err := validateRequest(c); err != nil {
//
//            //msg := err.Error()
//            code := err.(*HttpError).Code()
//            expl := err.(*HttpError).Explination()
//
//            c.AbortWithError(code, err);
//            log.Println(expl)
//            return
//        }
//
//        c.Next()
//    }
//}
//
//
//// Decode the JSON request and verify it.
//func verifyJSON(c *gin.Context, AppID string) error  {
//
//    r := c.Request
//
//    var echoReq *AlexaRequest
//
//    raw, err := ioutil.ReadAll(r.Body)
//    if err != nil {
//        return NewError("Bad Request", 400)
//    }
//
//    log.Println(string(raw))
//    err = json.Unmarshal(raw, &echoReq)
//
//    //
//    //err := json.NewDecoder(r.Body).Decode(&echoReq)
//    if err != nil {
//        return NewError("Bad Request", 400)
//    }
//
//    // Check the timestamp
//    if !echoReq.VerifyTimestamp() && r.URL.Query().Get("_dev") == "" {
//        return NewError("Bad Request", 400, "Request too old to continue (>150s).")
//    }
//
//    // Check the app id
//    if !echoReq.VerifyAppID(AppID) {
//        return NewError("Bad Request", 400, "Alexa AppID mismatch!")
//    }
//
//    c.Set("alexaRequest", echoReq)
//    return nil
//}
//
//// Run all mandatory Amazon security checks on the request.
//func validateRequest(c *gin.Context) error {
//
//    r := c.Request
//
//    // Check for debug bypass flag
//    devFlag := r.URL.Query().Get("_dev")
//
//    isDev := devFlag != ""
//
//    if !isDev {
//
//        certURL := r.Header.Get("SignatureCertChainUrl")
//        log.Println(certURL)
//
//        // Verify certificate URL
//        if !verifyCertURL(certURL) && devFlag == "" {
//            return NewError("Not Authorized", 401, "Invalid cert URL: "+certURL)
//        }
//
//        // Fetch certificate data
//        certContents, err := readCert(certURL)
//        if err != nil {
//            return NewError("Not Authorized", 401, err.Error())
//        }
//        log.Println("cert:", string(certContents))
//
//        // Decode certificate data
//        block, _ := pem.Decode(certContents)
//        if block == nil {
//            return NewError("Not Authorized", 401, "Failed to parse certificate PEM.")
//        }
//
//        cert, err := x509.ParseCertificate(block.Bytes)
//        if err != nil {
//            return NewError("Not Authorized", 401)
//        }
//
//        // Check the certificate date
//        if time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix() {
//            return NewError("Not Authorized", 401, "Amazon certificate expired.")
//        }
//
//        // Check the certificate alternate names
//        foundName := false
//        for _, altName := range cert.Subject.Names {
//            if altName.Value == "echo-api.amazon.com" {
//                foundName = true
//            }
//        }
//
//        if !foundName && devFlag == "" {
//            return NewError("Not Authorized", 401, "Amazon certificate invalid.")
//        }
//
//        // Verify the key
//        publicKey := cert.PublicKey
//        encryptedSig, _ := base64.StdEncoding.DecodeString(r.Header.Get("Signature"))
//        //log.Println("encsig:", string(encryptedSig))
//
//        // Make the request body SHA1 and verify the request with the public key
//        var bodyBuf bytes.Buffer
//        hash := sha1.New()
//        _, err = io.Copy(hash, io.TeeReader(r.Body, &bodyBuf))
//        if err != nil {
//            return NewError("Internal Error", 500, err.Error())
//        }
//        //log.Println(bodyBuf.String())
//        r.Body = ioutil.NopCloser(&bodyBuf)
//
//        sum := hash.Sum(nil)
//        log.Println("sum:", string(sum))
//
//        err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, sum, encryptedSig)
//        if err != nil {
//            return NewError("Not Authorized", 401, "Signature match failed.")
//        }
//    }
//
//    return nil
//}
//
//func readCert(certURL string) ([]byte, error) {
//    cert, err := http.Get(certURL)
//    if err != nil {
//        return nil, errors.New("Could not download Amazon cert file.")
//    }
//    defer cert.Body.Close()
//    certContents, err := ioutil.ReadAll(cert.Body)
//    if err != nil {
//        return nil, errors.New("Could not read Amazon cert file.")
//    }
//
//    return certContents, nil
//}
//
//func verifyCertURL(path string) bool {
//
//    link, _ := url.Parse(path)
//
//    if link.Scheme != "https" {
//        return false
//    }
//
//    if link.Host != "s3.amazonaws.com" && link.Host != "s3.amazonaws.com:443" {
//        return false
//    }
//
//    if !strings.HasPrefix(link.Path, "/echo.api/") {
//        return false
//    }
//
//    return true
//}