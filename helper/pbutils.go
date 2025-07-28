package helper

import (
	api "github.com/calculi-corp/api/go"
	"github.com/calculi-corp/api/go/endpoint"
	cutils "github.com/calculi-corp/common/utils"
)

func GetDashboardLayoutForEndpoint(ep *endpoint.GetUserPreferencesResponse, dashboardName string) *api.BinaryData {
	property := cutils.FindPropertyByName(ep.GetProperties(), dashboardName)
	dashboardLayout := property.GetData()
	return dashboardLayout
}
