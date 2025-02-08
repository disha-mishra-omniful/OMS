package router

import (
	//"awesomeProject5/OMS/orders/services"
	"context"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	// "github.com/omniful/go_commons/jwt/private"
	"awesomeProject5/OMS/orders"
)

func Initialize(ctx context.Context, s *http.Server) (err error) {
	s.Engine.Use(log.RequestLogMiddleware(log.MiddlewareOptions{
		Format:      config.GetString(ctx, "log.format"),
		Level:       config.GetString(ctx, "log.level"),
		LogRequest:  config.GetBool(ctx, "log.requests"),
		LogResponse: config.GetBool(ctx, "log.responses"),
	}))

	oms_v1 := s.Engine.Group("/api/v1")

	var CSVUploadController *orders.CSVUploadController
	// oms_v1.POST("/create-order",omscontroller.CreateOrder)
	// oms_v1.POST("/create-bulk",omscontroller.CreateBulkOrder)
	// oms_v1.GET("/:order_id", private.AuthenticateJWT(), omscontroller.GetOrder)
	// oms_v1.GET("/orders",omscontroller.GetOrders)
	oms_v1.POST("/create", CSVUploadController.CreateBulkCsv)

	return nil

}
