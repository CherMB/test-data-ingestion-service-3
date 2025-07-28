package db

import (
	"errors"
)

var (
	//ErrDBConnection for database connection errors
	ErrDBConnection = errors.New("DB connection error")
	//ErrProtoMarshal for proto marshalling erros
	ErrProtoMarshal = errors.New("error in marshaling proto")
	//ErrInvalidRequest for missing grpc request attributes
	ErrInvalidRequest = errors.New("missing request attributes")
	//ErrEmptyDbData for no db data for the select query and condition(s)
	ErrEmptyDbData = errors.New("no db data for the select request")
	// ErrProcessingReport Failed to process report
	ErrProcessingReport = errors.New("failed to process report")
	// ErrProcessingDate Failed to process report
	ErrParsingDate = errors.New("failed to parse date")
	// ErrFileNotFound File not found
	ErrFileNotFound = errors.New("file not found")
	// ErrInternalServer for internal server error
	ErrInternalServer = errors.New("internal server error")
	// ErrTransformationFailed when transformation returns null
	ErrTransformationFailed = errors.New("transformation operation returned null")
	// ErrInvalidDuration when there's an invalid duration type
	ErrInvalidDuration = errors.New("invalid duration type")
	// ErrDataIngestionFailed when ingesting data
	ErrDataIngestionFailed = errors.New("data ingestion failed")
	// ErrMultiSearchResponseError when mismatch in multi-search response
	ErrMultiSearchResponseError = errors.New("multi search respone mismatch")
	// ErrNoDataFound when OpenSearch returns no data for a query
	ErrNoDataFound              = errors.New("No data found")
)
