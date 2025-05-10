package route

import (
	handler "tracking-service/internal/handlers"
	middleware "tracking-service/internal/middlewares"
	service "tracking-service/internal/services"

	"github.com/gin-gonic/gin"
)

type TenantRoutes struct {
	service *service.ApplicationService
	handler *handler.TenantHandler
}

func NewTenantRoutes(
	service *service.ApplicationService,
	handler *handler.TenantHandler,
) *TenantRoutes {
	return &TenantRoutes{
		service: service,
		handler: handler,
	}
}

func (ur *TenantRoutes) RegisterRoutes(r *gin.Engine) {
	group := r.Group(
		"/tenant",
		middleware.TenantAuthMiddleware(ur.service),
	)

	group.GET("/platforms", ur.handler.GetPlatforms)

	group.GET("/profile", ur.handler.GetApp)

	group.POST("/events", ur.handler.CreateEvent)
	group.GET("/events/:event_id", ur.handler.GetEvent)
	group.PUT("/events/:event_id", ur.handler.UpdateEvent)
	group.DELETE("/events/:event_id", ur.handler.DeleteEvent)
	group.GET("/events", ur.handler.GetEvents)
	group.POST("/events/:event_id/fields", ur.handler.CreateEventField)
	group.GET("/events/:event_id/fields/:field_id", ur.handler.GetEventField)
	group.PUT("/events/:event_id/fields/:field_id", ur.handler.UpdateEventField)
	group.DELETE("/events/:event_id/fields/:field_id", ur.handler.DeleteEventField)

	group.POST("/events/:event_id/logs", ur.handler.CreateEventLog)

	group.POST("/sessions", ur.handler.CreateSession)
	group.GET("/sessions/:session_id", ur.handler.GetSession)
	group.PUT("/sessions/:session_id", ur.handler.UpdateSession)
	group.DELETE("/sessions/:session_id", ur.handler.DeleteSession)

}
