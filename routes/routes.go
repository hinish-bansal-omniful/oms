package routes

import (
	"context"
	"oms-service/controllers"

	// "github.com/newrelic/go-agent/v3/integrations/nrgin"

	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	// "github.com/omniful/go_commons/newrelic"
)

func Initialize(ctx context.Context, s *http.Server) error {
	// Health Check Route
	s.GET("/health", controllers.GetHealth)

	// API v1 Routes Group
	v1 := s.Engine.Group("/api/v1")
	{
		// Hubs Routes
		hubs := v1.Group("/orders")
		{
			hubs.GET("", controllers.GetAllOrders)
			hubs.GET("/:id", controllers.GetOrderByID)
			hubs.POST("/bulkorder", controllers.CreateBulkOrders)
		}
	}

	log.Infof("Routes initialized successfully")
	return nil
}
