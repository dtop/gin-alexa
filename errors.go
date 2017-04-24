package gin_alexa

type HttpError struct {
    error string
    code int
    expl string
}

func NewError(msg string, code int, explination ...string) HttpError {

    expl := ""
    if len(explination) > 0 {
        expl = explination[0]
    }

    return HttpError{error:msg, code:code, expl: expl}
}

func (e HttpError) Error() string {
    return e.error
}

func (e HttpError) Code() int {
    return e.code
}

func (e HttpError) Explination() string {
    return e.expl
}