package ginalexa

type (
	Session interface {
		New(ssid string) (Session, error)
		Store() error
		Get(key string, ptr interface{}, force ...bool) error
		GetGeneric(key string) interface{}
		Set(key string, val interface{})
		Del(key string)
		DeleteSession()
	}
)
