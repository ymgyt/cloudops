package core

// Backend -
type Backend interface {
	Put(*PutInput) (*PutOutput, error)
}

// PutInput -
type PutInput struct{}

// PutOutput -
type PutOutput struct{}
