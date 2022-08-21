package transport

type SessionStat int32

const (
	Connected    SessionStat = 1
	Disconnected SessionStat = 2
)

type OnMessageFunc func(Session, interface{})
type OnConnStatFunc func(Session, SessionStat)
type NewSessionIDFunc func() string

func (s SessionStat) String() string {
	switch s {
	case Connected:
		return "connected"
	case Disconnected:
		return "disconnected"
	}
	return "unknown"
}

type SessionMeta interface {
	MetaLoad(key string) (interface{}, bool)
	MetaStore(key string, value interface{})
	MetaDelete(key string)
}

type Session interface {
	ID() string
	String() string

	RemoteAddr() string
	LocalAddr() string

	Send(interface{}) error
	Close()
	SessionMeta
}
