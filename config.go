package ginalexa

type Configurable interface {
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    GetFloat(key string) float64
    GetVal(key string) interface{}
    Get(key string, obj interface{})
}


