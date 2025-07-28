package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"text/template"

	pb "github.com/calculi-corp/api/go/vsm/report"
	"github.com/calculi-corp/config"
	"github.com/calculi-corp/reports-service/constants"
	"github.com/calculi-corp/reports-service/exceptions"

	"github.com/calculi-corp/log"
)

type DbQuery struct {
	AliasName   string
	QueryString string
}

type VsmDashboardLayout struct {
	DashboardName   string             `json:"dashboard_name"`
	UserID          string             `json:"user_id"`
	OrgID           string             `json:"org_id"`
	DashboardLayout pb.DashboardLayout `json:"dashboard_layout"`
}

// GetWidgetEntity returns the widget entity, i.e. the configuration data stored in the widget definition JSON for a particular widget ID
func GetWidgetEntity(widgetId string, replacements map[string]any) (*pb.WidgetEntity, error) {

	fileName, ok := WidgetDefinitionMap[widgetId]
	if !ok {
		log.Errorf(errors.New("widget definition not found"), "widget definition not found for Id : ", widgetId)
		return nil, ErrInternalServer
	}

	// open the JSON file
	jsonFile, err := os.Open(config.Config.GetString(constants.DB_REPORT_DEF_FILEPATH_TEMPLATE) + fileName)
	if err != nil {
		log.CheckErrorf(err, exceptions.ErrOpeningWidgetDefJsonFile)
		return nil, ErrFileNotFound
	}
	defer jsonFile.Close()

	// read from the JSON file
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.CheckErrorf(err, exceptions.ErrReadingWidgetDefJsonFile)
		return nil, ErrFileNotFound
	}

	var we = pb.WidgetEntity{}
	found := false
	var replaceWidgets = []string{"cs3-sub", "e9-sub", "css3-sub", "es9-sub", "d1", "d2", "d3", "d4", "ds1", "ds2", "ds3", "ds4"}
	for _, replaceWidget := range replaceWidgets {
		if widgetId == replaceWidget {
			found = true
			break
		}
	}
	if found {
		updatedJSON, err := ReplaceJSONplaceholders(replacements, string(byteValue))
		if log.CheckErrorf(err, "could not replace json placeholders :", string(byteValue)) {
			return nil, err
		}
		err = json.Unmarshal([]byte(updatedJSON), &we)
	} else {
		err = json.Unmarshal(byteValue, &we)
	}

	if log.CheckErrorf(err, "error fetching data from DB : ") {
		return nil, ErrDBConnection
	}
	return &we, nil
}

// Returns the configuration for a particular widget in its corresponding proto structure
func GetComponentComparisonConfig(widgetId string, replacements map[string]any) (*pb.ComponentComparisonConfig, error) {

	configData, err := readWidgetConfigFile(widgetId)
	if log.CheckErrorf(err, "error fetching widget configuration data: ") {
		return nil, err
	}

	var widgetConfig = pb.ComponentComparisonConfig{}

	err = json.Unmarshal(configData, &widgetConfig)
	if log.CheckErrorf(err, "error unmarshalling widget configuration: ") {
		return nil, ErrDBConnection
	}
	return &widgetConfig, nil
}

// Reads widget configuration data from the configuration file for a particular widget
func readWidgetConfigFile(widgetId string) ([]byte, error) {

	fileName, ok := WidgetDefinitionMap[widgetId]
	if !ok {
		log.Errorf(errors.New("widget definition not found"), "widget definition not found for Id : ", widgetId)
		return nil, ErrInternalServer
	}

	// open the JSON file
	jsonFile, err := os.Open(config.Config.GetString(constants.DB_REPORT_DEF_FILEPATH_TEMPLATE) + fileName)
	if err != nil {
		log.CheckErrorf(err, exceptions.ErrOpeningWidgetDefJsonFile)
		return nil, ErrFileNotFound
	}
	defer jsonFile.Close()

	// read from the JSON file
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.CheckErrorf(err, exceptions.ErrReadingWidgetDefJsonFile)
		return nil, ErrFileNotFound
	}

	return byteValue, nil
}

func GetReportEntity(reportId int64) (*pb.ScaReportEntity, error) {

	fileName, ok := ReportDefinitionMap[reportId]
	if !ok {
		log.Errorf(errors.New("report definition not found"), "report definition not found for Id : ", reportId)
		return nil, ErrInternalServer
	}

	// open the JSON file
	jsonFile, err := os.Open(config.Config.GetString(constants.DB_REPORT_DEF_FILEPATH_TEMPLATE) + "sca/" + fileName)
	if err != nil {
		log.CheckErrorf(err, "error opening report definition json file ")
		return nil, ErrFileNotFound
	}
	defer jsonFile.Close()

	// read from the JSON file
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.CheckErrorf(err, "error reading report definition json file: ")
		return nil, ErrFileNotFound
	}

	var re = pb.ScaReportEntity{}
	err = json.Unmarshal(byteValue, &re)
	if err != nil {
		log.CheckErrorf(err, "error fetching data from DB ")
		return nil, ErrDBConnection
	}

	return &re, nil
}

/*
ReplaceJSONplaceholders replaces placeholders in a string.
Placeholders to be replaced by a string are to be specified as {{.placeholder}}
and placeholders to be replaced by JSON data are to be specified as {{json .placeholder}} in the input string.
Everything between and including the double curly brackets are replaced.
*/
func ReplaceJSONplaceholders(data interface{}, query string) (string, error) {
	//adding a function called "json" which can be used to marshal json data into a string during placeholder replacement, hence enabling us to replace more than just strings
	var t = template.Must(template.New("").Funcs(template.FuncMap{
		"json": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(query))
	buf := bytes.Buffer{}
	err := t.Execute(&buf, data)
	if log.CheckErrorf(err, "failed to replace placeholders") {
		return "", err
	}
	return buf.String(), nil
}

type WidgetDefinition struct {
	Id string  `json:"id"`
	W  float32 `json:"w"`
	H  float32 `json:"h"`
}

func GetWidgetDefinitionList(dashboardName string) ([]*pb.ReportWidget, error) {
	fileName, ok := DashboardWidgetsDefinitionMap[dashboardName]
	if !ok {
		log.Errorf(errors.New("widgets definition not found"), "widgets definition not found for dashboard : ", dashboardName)
		return nil, ErrInternalServer
	}
	jsonFile, err := os.Open(config.Config.GetString(constants.DB_REPORT_DEF_FILEPATH_TEMPLATE) + fileName)
	if err != nil {
		log.CheckErrorf(err, "error opening widgets definition json file ")
		return nil, ErrFileNotFound
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.CheckErrorf(err, exceptions.ErrReadingWidgetDefJsonFile)
		return nil, ErrInvalidRequest
	}
	widgets := make([]WidgetDefinition, 0)
	err = json.Unmarshal(byteValue, &widgets)
	if err != nil {
		log.CheckErrorf(err, "error while unmarshalling JSON data ")
		return nil, ErrInternalServer
	}
	pbRes := make([]*pb.ReportWidget, len(widgets))

	for i, widget := range widgets {
		pbRes[i] = &pb.ReportWidget{
			Id: widget.Id,
			W:  float32(widget.W),
			H:  float32(widget.H),
		}
	}
	return pbRes, nil
}
