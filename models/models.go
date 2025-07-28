package models

import (
	"context"
	"time"

	pb "github.com/calculi-corp/api/go/vsm/report"
	"github.com/opensearch-project/opensearch-go"
)

type UpdateProjectActivityResponse struct {
	Response       string
	ReportId       string
	Replacements   map[string]any
	Client         *opensearch.Client
	OutputResponse map[string]interface{}
	IsIdle         bool
	JobIds         []string
	IsFragile      bool
}

type GetSubReportWidget struct {
	FilterData       []string
	BaseData         []string
	Req              *pb.ReportServiceRequest
	Replacements     map[string]any
	Ctx              context.Context
	ReplacementsSpec map[string]any
	ParentWidget     *pb.Widget
}

type CalculateDateBydurationType struct {
	StartDateStr       string
	EndDateStr         string
	DurationType       pb.DurationType
	Replacements       map[string]any
	ReplacementsSpec   map[string]any
	NormalizeMonthFlag bool
	IsComputFlag       bool
	CurrentTime        time.Time
}
