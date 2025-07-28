package helper

import (
	"reflect"
	"testing"

	api "github.com/calculi-corp/api/go"
	"github.com/calculi-corp/api/go/endpoint"
	"github.com/calculi-corp/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.Config.Set("logging.level", "INFO")
}

func TestGetDashboardLayoutForEndpoint(t *testing.T) {

	ep := &endpoint.GetUserPreferencesResponse{
		Properties: []*api.Property{
			{
				Name:        "dashboard1",
				Description: "dashboard1",
			},
			{
				Name:        "dashboard2",
				Description: "dashboard2",
			},
		},
	}
	dashboardName := "dashboard1"
	result := GetDashboardLayoutForEndpoint(ep, dashboardName)
	assert.Nil(t,result)
}

func TestAddTermFilter(t *testing.T) {
	filterName := "widget1"
	filterValue := "data1"

	expected := map[string]interface{}{
		"term": map[string]interface{}{
			"widget1": "data1",
		},
	}

	result := AddTermFilter(filterName, filterValue)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("AddTermFilter() returned unexpected result: got %v, want %v", result, expected)
	}
}

func TestAddTermsFilter(t *testing.T) {
	filterName := "widget1"
	filterValues := []string{"data1", "data2", "data3"}

	expected := map[string]interface{}{
		"terms": map[string]interface{}{
			"widget1": []string{"data1", "data2", "data3"},
		},
	}

	result := AddTermsFilter(filterName, filterValues)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("AddTermsFilter() returned unexpected result: got %v, want %v", result, expected)
	}
}

func TestCalculatePreviousDates(t *testing.T) {
	type args struct {
		inputDate string
		duration  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "Month validation",
			args: args{
				inputDate: "2023-10-20 15:04:05",
				duration:  "month",
			},
			want:  "2023-09-01",
			want1: "2023-09-30",
		},
		{
			name: "Week validation",
			args: args{
				inputDate: "2023-10-20 15:04:05",
				duration:  "week",
			},
			want:  "2023-10-09",
			want1: "2023-10-15",
		},
		{
			name: "Validate unknown type",
			args: args{
				inputDate: "2023-10-20 15:04:05",
				duration:  "test",
			},
			wantErr: true,
		},
		{
			name: "Validate date conversion error",
			args: args{
				inputDate: "2023/10/20 15:04:05",
				duration:  "week",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CalculatePreviousDates(tt.args.inputDate, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculatePreviousDates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculatePreviousDates() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CalculatePreviousDates() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
