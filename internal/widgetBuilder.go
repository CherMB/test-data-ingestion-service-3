package internal

import (
	"encoding/json"
	"sync"

	pb "github.com/calculi-corp/api/go/vsm/report"
	"github.com/calculi-corp/log"
	"github.com/calculi-corp/reports-service/exceptions"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type WidgetBuilder struct {
	widget pb.Widget
}

func newWidgetBuilder() *WidgetBuilder {
	return &WidgetBuilder{}
}

func (wb *WidgetBuilder) setWidgetInfo(w *pb.Widget) {
	wb.widget.Id = w.Id
	wb.widget.Title = w.Title
	wb.widget.Description = w.Description
	wb.widget.Content = w.Content
	wb.widget.Data = w.Data
	wb.widget.Pagination = w.Pagination
	wb.widget.EnableComponentsCompare = w.EnableComponentsCompare
	wb.widget.ComponentsCompareId = w.ComponentsCompareId
}

// Build the Headers.
func (wb *WidgetBuilder) setHeaders(qr map[string]json.RawMessage, fr map[string]json.RawMessage) error {
	for i, ele := range wb.widget.Content {
		for j, v := range ele.Header {
			var structValue structpb.Struct
			if len(v.FunctionName) > 0 {
				if len(v.SpecKey) > 0 { //Expect fr to contain map[string]interface{}
					fnData := make(map[string]json.RawMessage)
					err := json.Unmarshal(fr[v.FunctionName], &fnData)
					if err != nil {
						return err
					}
					if _, ok := fnData[v.SpecKey]; ok {
						err := protojson.Unmarshal([]byte(fnData[v.SpecKey]), &structValue)
						if err != nil {
							return err
						}
					}
				} else {
					err := protojson.Unmarshal(fr[v.FunctionName], &structValue)
					if err != nil {
						return err
					}
				}
			} else if len(v.PostProcessFunctionName) > 0 {
				data, err := ExecutePostProcessFunction(v.PostProcessFunctionName, v.SpecKey, qr, nil)
				if log.CheckErrorf(err, exceptions.ErrExecutePostProcess, v.PostProcessFunctionName) {
					return err
				}
				err = protojson.Unmarshal(data, &structValue)
				if err != nil {
					return err
				}
			}
			wb.widget.Content[i].Header[j] = &pb.MetricInfo{
				Title:                   v.Title,
				Description:             v.Description,
				Data:                    &structValue,
				Type:                    v.Type,
				DrillDown:               v.DrillDown,
				EnableComponentsCompare: v.EnableComponentsCompare,
				ComponentsCompareId:     v.ComponentsCompareId,
			}
		}
	}
	return nil
}

// Build the Footer.
func (wb *WidgetBuilder) setFooter(qr map[string]json.RawMessage, fr map[string]json.RawMessage) error {

	for i, ele := range wb.widget.Content {
		for j, v := range ele.Footer {
			var structValue structpb.Struct

			if len(v.FunctionName) > 0 {
				//TBD past duration here
				err := protojson.Unmarshal(fr[v.FunctionName], &structValue)
				if err != nil {
					return err
				}

			} else if len(v.PostProcessFunctionName) > 0 {
				data, err := ExecutePostProcessFunction(v.PostProcessFunctionName, v.SpecKey, qr, nil)
				if log.CheckErrorf(err, exceptions.ErrExecutePostProcess, v.PostProcessFunctionName) {
					return err
				}
				err = protojson.Unmarshal(data, &structValue)
				if err != nil {
					return err
				}
			}
			wb.widget.Content[i].Footer[j] = &pb.MetricInfo{
				Title:       v.Title,
				Description: v.Description,
				Data:        &structValue,
			}
		}
	}
	return nil
}

type DataInfo struct {
	Data json.RawMessage `json:"data"`
	Info json.RawMessage `json:"info"`
}

// Build the Chart section.
func (wb *WidgetBuilder) setSection(qr map[string]json.RawMessage, fr map[string]json.RawMessage, replacements map[string]any) error {

	for i, ele := range wb.widget.Content {
		for j, v := range ele.Section {

			wb.widget.Content[i].Section[j] = &pb.ChartInfo{
				Title:            v.Title,
				Type:             v.Type,
				CategoryType:     v.CategoryType,
				ShowLegends:      v.ShowLegends,
				ColorScheme:      v.ColorScheme,
				LightColorScheme: v.LightColorScheme,
				DataType:         v.DataType,
				ColumnType:       v.ColumnType,
				ShowPagination:   v.ShowPagination,
				Orientation:      v.Orientation,
				DrillDown:        v.DrillDown,
				LegendsData:      v.LegendsData,
			}

			var dataValue structpb.ListValue
			if len(v.FunctionName) > 0 {
				if len(fr[v.FunctionName]) == 0 {
					log.Warn("Empty function data")
				} else {
					shouldReturn, returnValue := dataInfoSeperation(v, fr[v.FunctionName], &dataValue, wb, i, j)
					if shouldReturn {
						return returnValue
					}
				}
			} else if len(v.PostProcessFunctionName) > 0 {
				data, err := ExecutePostProcessFunction(v.PostProcessFunctionName, v.SpecKey, qr, replacements)
				if log.CheckErrorf(err, exceptions.ErrExecutePostProcess, v.PostProcessFunctionName) {
					return err
				}
				shouldReturn, returnValue := dataInfoSeperation(v, data, &dataValue, wb, i, j)
				if shouldReturn {
					return returnValue
				}
			}
			wb.widget.Content[i].Section[j].Data = &dataValue
		}
	}
	return nil
}

func unmarshalBatchProtoJSON(data []json.RawMessage, structValue *structpb.Struct, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	for _, record := range data {
		msg := structpb.Struct{}
		err := protojson.Unmarshal(record, &msg)
		if err != nil {
			log.Errorf(err, "Error unmarshaling record:")
			continue
		}
		mutex.Lock()
		*structValue = msg
		mutex.Unlock()
	}
}

func unmarshalProtoJSON(data []json.RawMessage) (structpb.Struct, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	numRecords := len(data)
	batchSize := 20000
	numBatches := (numRecords + batchSize - 1) / batchSize
	var structValue structpb.Struct

	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > numRecords {
			end = numRecords
		}
		wg.Add(1)
		go unmarshalBatchProtoJSON(data[start:end], &structValue, &wg, &mutex)
	}

	wg.Wait()

	return structValue, nil
}

// Build the Data.
func (wb *WidgetBuilder) setData(qr map[string]json.RawMessage, fr map[string]json.RawMessage) error {
	data := wb.widget.Data
	if data != nil {
		fields := data.GetFields()
		functionName, ok := fields["function_name"]
		var structValue structpb.Struct
		if ok {
			if functionName.GetStringValue() == "Insight Completed Runs Widget" {
				batchData := []json.RawMessage{fr[functionName.GetStringValue()]}
				batchStruct, err := unmarshalProtoJSON(batchData)
				if err != nil {
					return err
				}
				structValue = batchStruct
			} else {
				err := protojson.Unmarshal(fr[functionName.GetStringValue()], &structValue)
				if err != nil {
					return err
				}
			}
		}
		for key, value := range structValue.Fields {
			fields[key] = value
		}
		data.Fields = fields
		delete(data.Fields, "function_name")
		delete(data.Fields, "query_key")
		delete(data.Fields, "spec_key")
		wb.widget.Data = data
	}
	return nil
}

func dataInfoSeperation(v *pb.ChartInfo, data json.RawMessage, dataValue *structpb.ListValue, wb *WidgetBuilder, i int, j int) (bool, error) {
	if v.DataType == 1 {
		var infoValue structpb.ListValue
		var result DataInfo
		err := json.Unmarshal(data, &result)
		if log.CheckErrorf(err, "Exception while transforming function : ") {
			return true, err
		} else if result.Data != nil && result.Info != nil {
			err := protojson.Unmarshal(result.Data, dataValue)
			if log.CheckErrorf(err, "Exception while transforming function data : ") {
				return true, err
			}
			err = protojson.Unmarshal(result.Info, &infoValue)
			if log.CheckErrorf(err, "Exception while transforming function info : ") {
				return true, err
			}
		}
		wb.widget.Content[i].Section[j].Info = &infoValue
	} else {
		err := protojson.Unmarshal(data, dataValue)
		if err != nil {
			return true, err
		}
	}
	return false, nil
}

// Get all widget contents
func (wb *WidgetBuilder) getWidget() pb.Widget {
	return pb.Widget{
		Id:                      wb.widget.Id,
		Title:                   wb.widget.Title,
		Description:             wb.widget.Description,
		Content:                 wb.widget.Content,
		Data:                    wb.widget.Data,
		Pagination:              wb.widget.Pagination,
		EnableComponentsCompare: wb.widget.EnableComponentsCompare,
		ComponentsCompareId:     wb.widget.ComponentsCompareId,
	}
}
