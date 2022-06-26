package websocket

type Option func(*Options)

type Options struct {
	Address string
	// Adapter gate.MethodAdapter
}
