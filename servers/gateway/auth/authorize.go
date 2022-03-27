package auth

type GateAuth interface {
	Auth() error
}
