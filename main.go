package main

import (
	"os"
	"github.com/calculi-corp/api/go/auth/permission"
	pb "github.com/calculi-corp/api/go/vsm/report"
	"github.com/calculi-corp/common/pkg/messaging"
	"github.com/calculi-corp/common/pkg/service"
	"github.com/calculi-corp/config"
	hostflags "github.com/calculi-corp/grpc-hostflags"
	"github.com/calculi-corp/log"
	"github.com/calculi-corp/nats"
	opensearchconfig "github.com/calculi-corp/opensearch-config"
	"github.com/calculi-corp/reports-service/cache"
	"github.com/calculi-corp/reports-service/handler"
)

const (
	handlerName = "ReportsHandler"
)

func init() {
	config.Config.DefineStringFlag("report.definition.filepath", "", "Report Definiftion filepath")
	opensearchconfig.DefineOpensearchFlags()
}

func main() {

	// read command line flags, config files, and environment variables to set config values
	config.Config.SetCliFlags()

	// Test Opensearch connection
	_, e := opensearchconfig.GetOpensearchConnection()
	if log.CheckErrorf(e, "could not connect to Opensearch") {
		panic(e)
	} else {
		log.Info("Successfullly connected to OpenSearch")
	}

	svc := service.NewService(hostflags.ReportsService)

	msgClt, err := nats.NewMessagingClient()
	if log.CheckErrorf(err, "unexpected error connecting NATS client") {
		panic(e)
	}
	initCaches(svc, msgClt)
	initAuthInterceptors(svc)

	svc.AddHandler(handlerName, service.NewHandlerWrapper(handler.NewDefaultReportsHandler), true)

	svc.Run()
}

func initAuthInterceptors(svc *service.Service) {

	svcName := pb.ReportServiceHandler_ServiceDesc.ServiceName
	err := svc.Server.GetAuthInterceptor().AddService(svcName, permission.ApiType_SERVICE)
	if err != nil {
		log.Fatal("Error adding reports service to interceptor", err)
	}

	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_GetReportData_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_BuildReportLayout_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_BuildDrilldownReport_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_BuildReport_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_GetEnvironments_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_BuildComputedDrilldownReport_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_BuildComputedReport_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_GetInsightsIntegration_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_GetDashboardLayout_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_UpdateDashboardLayout_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_GetWidgets_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_BuildComponentComparisonReport_FullMethodName, permission.PermissionAction_READ)
	svc.Server.GetAuthInterceptor().AddPermission(pb.ReportServiceHandler_StreamCIInsightsCompletedRun_FullMethodName, permission.PermissionAction_READ)

}

func initCaches(svc *service.Service, msgClt messaging.Messaging) {
	err := cache.InitializeCache(svc.Client, msgClt)
	if err != nil {
		log.Error("Could not initialize resource and endpoint cache",err)
		os.Exit(1)
	}
}
