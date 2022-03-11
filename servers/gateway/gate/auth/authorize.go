package auth

type GateAuth interface {
	// Check

	Auth() error
}
