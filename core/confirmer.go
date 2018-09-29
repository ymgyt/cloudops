package core

//go:generate mockgen -destination=../testutil/confirmer_mock.go -package=testutil github.com/ymgyt/cloudops/core Confirmer

// Confirmer have responsibility to handle user operation confirmation.
type Confirmer interface {
	Confirm(operation string, resources Resources) (bool, error)
}
