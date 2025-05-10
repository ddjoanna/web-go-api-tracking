package route

import (
	shared "tracking-service/internal"
	handler "tracking-service/internal/handlers"
	middleware "tracking-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type AdminRoutes struct {
	config  *shared.Config
	handler *handler.AdminHandler
}

func NewAdminRoutes(
	config *shared.Config,
	handler *handler.AdminHandler,
) *AdminRoutes {
	return &AdminRoutes{
		config:  config,
		handler: handler,
	}
}

func (ar *AdminRoutes) RegisterRoutes(r *gin.Engine) {
	group := r.Group(
		"/admin",
		middleware.AdminAuthMiddleware(ar.config),
	)

	group.POST("/tenants", ar.handler.CreateTenant)
	group.GET("/tenants/:tenant_id", ar.handler.GetTenant)
	group.PUT("/tenants/:tenant_id", ar.handler.UpdateTenant)
	group.GET("/tenants", ar.handler.GetTenants)

	// 多租戶共用不提供修改和刪除，特殊需求請新增(或依需求擴充功能)
	group.POST("/platforms", ar.handler.CreatePlatform)
	group.GET("/platforms/:platform_id", ar.handler.GetPlatform)
	group.GET("/platforms", ar.handler.GetPlatforms)

	group.POST("/apps", ar.handler.CreateApp)
	group.GET("/apps/:app_id", ar.handler.GetApp)
	group.PUT("/apps/:app_id", ar.handler.UpdateApp)
	group.DELETE("/apps/:app_id", ar.handler.DeleteApp)
	group.GET("/apps", ar.handler.GetApps)
	group.POST("/apps/:app_id/api-keys", ar.handler.CreateAppAPIKey)
	group.DELETE("/apps/:app_id/api-keys/:api_key_id", ar.handler.DeleteAppAPIKey)

	group.POST("/apps/:app_id/events", ar.handler.CreateEvent)
	group.GET("/apps/:app_id/events/:event_id", ar.handler.GetEvent)
	group.PUT("/apps/:app_id/events/:event_id", ar.handler.UpdateEvent)
	group.DELETE("/apps/:app_id/events/:event_id", ar.handler.DeleteEvent)
	group.GET("/apps/:app_id/events", ar.handler.GetEvents)

	group.POST("/apps/:app_id/events/:event_id/fields", ar.handler.CreateEventFields)
	group.GET("/apps/:app_id/events/:event_id/fields/:field_id", ar.handler.GetEventField)
	group.PUT("/apps/:app_id/events/:event_id/fields/:field_id", ar.handler.UpdateEventField)
	group.DELETE("/apps/:app_id/events/:event_id/fields/:field_id", ar.handler.DeleteEventField)
	group.GET("/apps/:app_id/events/:event_id/fields", ar.handler.GetEventFields)
}
