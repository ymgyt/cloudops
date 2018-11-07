package core

import (
	"cloud.google.com/go/bigquery"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1beta1"
)

// GCPProjectService -
type GCPProjectService interface {
	List(*ListGCPProjectsInput) (*ListGCPProjectsOutput, error)
}

// ListGCPProjectsInput -
type ListGCPProjectsInput struct {
}

// ListGCPProjectsOutput -
type ListGCPProjectsOutput struct {
	Projects []*GCPProject
}

// GCPProject -
type GCPProject cloudresourcemanager.Project

// BQueryInput -
// legacySQLの切り替え必要..?
type BQueryInput struct {
	Query          string
	Dryrun         bool
	MaxBytesBilled int64 // query failed if read bytes exceed it
	DestDatasetID  string
	DestTableID    string

	Handler           BigqueryHandler
	CreateDisposition bigquery.TableCreateDisposition
	WriteDisposition  bigquery.TableWriteDisposition
}

// BigqueryHandler -
type BigqueryHandler func([]bigquery.Value)

// BQueryOutput -
type BQueryOutput struct {
	ProcessedBytes int64
}

// BigqueryService -
type BigqueryService interface {
	Query(*BQueryInput) (*BQueryOutput, error)
}
