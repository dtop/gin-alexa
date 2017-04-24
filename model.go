package gin_alexa

import (
    "time"
    "encoding/json"
    "errors"
)

// ****************************************************************
//
// Taken from SkillServer (https://github.com/mikeflynn/go-alexa)
//
// ****************************************************************


// Request Functions
func (a *AlexaRequest) VerifyTimestamp() bool {
    reqTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", a.Request.Timestamp)
    if time.Since(reqTimestamp) < time.Duration(150)*time.Second {
        return true
    }

    return false
}

func (a *AlexaRequest) VerifyAppID(myAppID string) bool {
    if a.Session.Application.ApplicationID == myAppID {
        return true
    }

    return false
}

func (a *AlexaRequest) GetSessionID() string {
    return a.Session.SessionID
}

func (a *AlexaRequest) GetUserID() string {
    return a.Session.User.UserID
}

func (a *AlexaRequest) GetRequestType() string {
    return a.Request.Type
}

func (a *AlexaRequest) GetIntentName() string {
    if a.GetRequestType() == "IntentRequest" {
        return a.Request.Intent.Name
    }

    return a.GetRequestType()
}

func (a *AlexaRequest) GetSlotValue(slotName string) (string, error) {
    if _, ok := a.Request.Intent.Slots[slotName]; ok {
        return a.Request.Intent.Slots[slotName].Value, nil
    }

    return "", errors.New("Slot name not found.")
}

func (a *AlexaRequest) AllSlots() map[string]AlexaSlot {
    return a.Request.Intent.Slots
}

// Response Functions
func NewAlexaResponse() *AlexaResponse {
    er := &AlexaResponse{
        Version: "1.0",
        Response: AlexaRespBody{
            ShouldEndSession: true,
        },
    }

    return er
}

func (a *AlexaResponse) OutputSpeech(text string) *AlexaResponse {
    a.Response.OutputSpeech = &AlexaRespPayload{
        Type: "PlainText",
        Text: text,
    }

    return a
}

func (a *AlexaResponse) OutputSpeechSSML(text string) *AlexaResponse {
    a.Response.OutputSpeech = &AlexaRespPayload{
        Type: "SSML",
        SSML: text,
    }

    return a
}

func (a *AlexaResponse) Card(title string, content string) *AlexaResponse {
    return a.SimpleCard(title, content)
}

func (a *AlexaResponse) SimpleCard(title string, content string) *AlexaResponse {
    a.Response.Card = &AlexaRespPayload{
        Type:    "Simple",
        Title:   title,
        Content: content,
    }

    return a
}

func (a *AlexaResponse) StandardCard(title string, content string, smallImg string, largeImg string) *AlexaResponse {
    a.Response.Card = &AlexaRespPayload{
        Type:    "Standard",
        Title:   title,
        Content: content,
    }

    if smallImg != "" {
        a.Response.Card.Image.SmallImageURL = smallImg
    }

    if largeImg != "" {
        a.Response.Card.Image.LargeImageURL = largeImg
    }

    return a
}

func (a *AlexaResponse) LinkAccountCard() *AlexaResponse {
    a.Response.Card = &AlexaRespPayload{
        Type: "LinkAccount",
    }

    return a
}

func (a *AlexaResponse) Reprompt(text string) *AlexaResponse {
    a.Response.Reprompt = &AlexaReprompt{
        OutputSpeech: AlexaRespPayload{
            Type: "PlainText",
            Text: text,
        },
    }

    return a
}

func (a *AlexaResponse) EndSession(flag bool) *AlexaResponse {
    a.Response.ShouldEndSession = flag

    return a
}

func (a *AlexaResponse) String() ([]byte, error) {
    jsonStr, err := json.Marshal(a)
    if err != nil {
        return nil, err
    }

    return jsonStr, nil
}

// Request Types

type AlexaRequest struct {
    Version string      `json:"version"`
    Session AlexaSession `json:"session"`
    Request AlexaReqBody `json:"request"`
}

type AlexaSession struct {
    New         bool   `json:"new"`
    SessionID   string `json:"sessionId"`
    Application struct {
                    ApplicationID string `json:"applicationId"`
                } `json:"application"`
    Attributes struct {
                    String map[string]interface{} `json:"string"`
                } `json:"attributes"`
    User struct {
                    UserID      string `json:"userId"`
                    AccessToken string `json:"accessToken,omitempty"`
                } `json:"user"`
}

type AlexaReqBody struct {
    Type      string     `json:"type"`
    RequestID string     `json:"requestId"`
    Timestamp string     `json:"timestamp"`
    Intent    AlexaIntent `json:"intent,omitempty"`
    Reason    string     `json:"reason,omitempty"`
}

type AlexaIntent struct {
    Name  string              `json:"name"`
    Slots map[string]AlexaSlot `json:"slots"`
}

type AlexaSlot struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}

// Response Types

type AlexaResponse struct {
    Version           string                 `json:"version"`
    SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
    Response          AlexaRespBody           `json:"response"`
}

type AlexaRespBody struct {
    OutputSpeech     *AlexaRespPayload `json:"outputSpeech,omitempty"`
    Card             *AlexaRespPayload `json:"card,omitempty"`
    Reprompt         *AlexaReprompt    `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
    ShouldEndSession bool             `json:"shouldEndSession"`
}

type AlexaReprompt struct {
    OutputSpeech AlexaRespPayload `json:"outputSpeech,omitempty"`
}

type AlexaRespImage struct {
    SmallImageURL string `json:"smallImageUrl,omitempty"`
    LargeImageURL string `json:"largeImageUrl,omitempty"`
}

type AlexaRespPayload struct {
    Type    string        `json:"type,omitempty"`
    Title   string        `json:"title,omitempty"`
    Text    string        `json:"text,omitempty"`
    SSML    string        `json:"ssml,omitempty"`
    Content string        `json:"content,omitempty"`
    Image   AlexaRespImage `json:"image,omitempty"`
}
