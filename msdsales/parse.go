package msdsales

import (
	"github.com/amp-labs/connectors/common"
	"github.com/spyzhov/ajson"
)

func getTotalSize(node *ajson.Node) (int64, error) {
	return common.JsonManager.GetInteger(node, "@odata.count")
}

func getRecords(node *ajson.Node) ([]map[string]any, error) {
	arr, err := common.JsonManager.GetArr(node, "value")
	if err != nil {
		return nil, err
	}

	return common.JsonManager.ArrToMap(arr)
}

func getNextRecordsURL(node *ajson.Node) (string, error) {
	return common.JsonManager.GetString(node, "@odata.nextLink")
}

// TODO we must differentiate between GET and LIST
// sometimes the fields that user requests are either in singular record or list of records or hybrid in array and outside
func getMarshaledData(records []map[string]interface{}, fields []string) ([]common.ReadResultRow, error) {
	data := make([]common.ReadResultRow, len(records))
	for i, record := range records {
		data[i] = common.ReadResultRow{
			Fields: common.ExtractLowercaseFieldsFromRaw(fields, record),
			Raw:    record,
		}
	}

	return data, nil
}
