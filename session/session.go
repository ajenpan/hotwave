package session

type Session interface {
	ID() string
	String() string
	Send(interface{}) error
	Close() error
}

type SessionStat int32

const (
	SessionStatConnected    SessionStat = 1
	SessionStatDisconnected SessionStat = 2
)

func (s SessionStat) String() string {
	switch s {
	case SessionStatConnected:
		return "connected"
	case SessionStatDisconnected:
		return "disconnected"
	}
	return "unknown"
}
