package core

//go:generate mockgen -destination=../testutil/backend_mock.go -package=testutil github.com/ymgyt/cloudops/core Backend

// Backend -
type Backend interface {
	Put(*PutInput) (*PutOutput, error)
	Fetch(*FetchInput) (*FetchOutput, error)
	Remove(*RemoveInput) (*RemoveOutput, error)
}

// PutInput -
type PutInput struct {
	Dryrun    bool
	CreateDir bool
	Recursive bool
	Dest      string
	Resources Resources
}

// PutOutput -
type PutOutput struct {
	PutNum int
}

// FetchInput -
type FetchInput struct {
	Recursive bool
	Regexp    string
	Src       string
}

// FetchOutput -
type FetchOutput struct {
	Resources Resources
}

// RemoveInput -
type RemoveInput struct {
	Dryrun bool
	Resources
}

// RemoveOutput -
type RemoveOutput struct {
	RemoveNum int
}
