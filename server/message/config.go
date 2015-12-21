package message

// The server config.
type Config interface {
	Export() []interface{}
	Import([]interface{})
}
