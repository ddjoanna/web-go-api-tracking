package handler

import (
	"strconv"
	datastructure "tracking-service/internal/datastructures"
	service "tracking-service/internal/services"
	util "tracking-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	BaseHandler
	tenant_service   *service.TenantService
	platform_service *service.PlatformService
	app_service      *service.ApplicationService
	event_service    *service.EventService
}

func NewAdminHandler(
	tenant_service *service.TenantService,
	platform_service *service.PlatformService,
	app_service *service.ApplicationService,
	event_service *service.EventService,
) *AdminHandler {
	return &AdminHandler{
		tenant_service:   tenant_service,
		platform_service: platform_service,
		app_service:      app_service,
		event_service:    event_service,
	}
}

// CreateTenant godoc
// @Summary      建立新租戶
// @Description  建立新租戶
// @Tags         Admin/Tenant
// @Accept       json
// @Produce      json
// @Param        request  body  datastructure.Tenant  true  "新增租戶資料"
// @Success      200      {object}  datastructure.BaseResponse{data=datastructure.Tenant}  "成功回應，包含新租戶資料"
// @Failure      400      {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401      {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403      {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404      {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409      {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500      {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/tenants [post]
func (h *AdminHandler) CreateTenant(c *gin.Context) {
	ctx := c.Request.Context()
	var req datastructure.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	reqTenant := &datastructure.Tenant{
		Name:        req.Name,
		Description: req.Description,
	}

	tenant, err := h.tenant_service.CreateTenant(ctx, reqTenant)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respTenant := datastructure.Tenant{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Description: tenant.Description,
		CreatedAt:   util.ConvertTimeToTimeStamp(&tenant.CreatedAt),
		UpdatedAt:   util.ConvertTimeToTimeStamp(&tenant.UpdatedAt),
		DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(tenant.DeletedAt),
	}

	h.Success(c, respTenant)
}

// GetTenant godoc
// @Summary      取得指定租戶詳細資料
// @Description  取得指定租戶詳細資料
// @Tags         Admin/Tenant
// @Produce      json
// @Param        tenant_id  path      string  true  "租戶 ID"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.Tenant}  "成功回應，包含應用程式詳細資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/tenants/{tenant_id} [get]
func (h *AdminHandler) GetTenant(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	tenant, err := h.tenant_service.GetTenant(c.Request.Context(), tenantID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respTenant := datastructure.Tenant{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Description: tenant.Description,
		CreatedAt:   util.ConvertTimeToTimeStamp(&tenant.CreatedAt),
		UpdatedAt:   util.ConvertTimeToTimeStamp(&tenant.UpdatedAt),
		DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(tenant.DeletedAt),
	}

	h.Success(c, respTenant)
}

// UpdateTenant godoc
// @Summary      更新指定租戶
// @Description  更新指定租戶
// @Tags         Admin/Tenant
// @Produce      json
// @Param        tenant_id  path      string  true  "租戶 ID"
// @Param        request  body  datastructure.Tenant  true  "更新租戶資料"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/tenants/{tenant_id} [put]
func (h *AdminHandler) UpdateTenant(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	var req datastructure.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	reqTenant := &datastructure.Tenant{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.tenant_service.UpdateTenant(c.Request.Context(), tenantID, reqTenant); err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// GetTenants godoc
// @Summary      取得所有租戶
// @Description  取得所有租戶
// @Tags         Admin/Tenant
// @Produce      json
// @Success      200     {object}  datastructure.BaseResponse{data=[]datastructure.Tenant}  "成功回應，包含所有租戶陣列"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/tenants [get]
func (h *AdminHandler) GetTenants(c *gin.Context) {
	tenants, err := h.tenant_service.GetTenants(c.Request.Context())
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respTenants := make([]*datastructure.Tenant, 0)
	for _, tenant := range tenants {
		respTenants = append(respTenants, &datastructure.Tenant{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Description: tenant.Description,
			CreatedAt:   util.ConvertTimeToTimeStamp(&tenant.CreatedAt),
			UpdatedAt:   util.ConvertTimeToTimeStamp(&tenant.UpdatedAt),
			DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(tenant.DeletedAt),
		})
	}

	h.Success(c, respTenants)
}

// CreatePlatform godoc
// @Summary      建立新平台
// @Description  建立新平台
// @Tags         Admin/Platform
// @Accept       json
// @Produce      json
// @Param        request  body  datastructure.Platform  true  "新增平台資料"
// @Success      200      {object}  datastructure.BaseResponse{data=datastructure.Platform}  "成功回應，包含新平台資料"
// @Failure      400      {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401      {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403      {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404      {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409      {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500      {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/platforms [post]
func (h *AdminHandler) CreatePlatform(c *gin.Context) {
	var req datastructure.CreatePlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	reqPlatform := &datastructure.Platform{
		Name: req.Name,
	}

	platform, err := h.platform_service.CreatePlatform(c.Request.Context(), reqPlatform)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respPlatform := datastructure.Platform{
		ID:        platform.ID,
		Name:      platform.Name,
		CreatedAt: util.ConvertTimeToTimeStamp(&platform.CreatedAt),
		UpdatedAt: util.ConvertTimeToTimeStamp(&platform.UpdatedAt),
		DeletedAt: util.ConvertGormDeletedAtToTimeStamp(platform.DeletedAt),
	}

	h.Success(c, respPlatform)
}

// GetPlatform godoc
// @Summary      取得指定平台詳細資料
// @Description  取得指定平台詳細資料
// @Tags         Admin/Platform
// @Produce      json
// @Param        platform_id  path      string  true  "平台 ID"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.Platform}  "成功回應，包含應用程式詳細資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/platforms/{platform_id} [get]
func (h *AdminHandler) GetPlatform(c *gin.Context) {
	// 字串轉整數
	platformID, err := strconv.Atoi(c.Param("platform_id"))
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	platform, err := h.platform_service.GetPlatformByID(c.Request.Context(), platformID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respPlatform := datastructure.Platform{
		ID:        platform.ID,
		Name:      platform.Name,
		CreatedAt: util.ConvertTimeToTimeStamp(&platform.CreatedAt),
		UpdatedAt: util.ConvertTimeToTimeStamp(&platform.UpdatedAt),
		DeletedAt: util.ConvertGormDeletedAtToTimeStamp(platform.DeletedAt),
	}

	h.Success(c, respPlatform)
}

// GetPlatforms godoc
// @Summary      取得所有平台
// @Description  取得所有平台
// @Tags         Admin/Platform
// @Produce      json
// @Success      200     {object}  datastructure.BaseResponse{data=[]datastructure.Platform}  "成功回應，包含所有平台陣列"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/platforms [get]
func (h *AdminHandler) GetPlatforms(c *gin.Context) {
	platforms, err := h.platform_service.GetPlatforms(c.Request.Context())
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respPlatforms := make([]*datastructure.Platform, 0)
	for _, platform := range platforms {
		respPlatforms = append(respPlatforms, &datastructure.Platform{
			ID:        platform.ID,
			Name:      platform.Name,
			CreatedAt: util.ConvertTimeToTimeStamp(&platform.CreatedAt),
			UpdatedAt: util.ConvertTimeToTimeStamp(&platform.UpdatedAt),
			DeletedAt: util.ConvertGormDeletedAtToTimeStamp(platform.DeletedAt),
		})
	}

	h.Success(c, respPlatforms)
}

// CreateApp godoc
// @Summary      建立新應用程式
// @Description  建立新應用程式
// @Tags         Admin/Application
// @Accept       json
// @Produce      json
// @Param        request  body  datastructure.Application  true  "新增應用程式資料"
// @Success      200      {object}  datastructure.BaseResponse{data=datastructure.Application}  "成功回應，包含新應用程式資料"
// @Failure      400      {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401      {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403      {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404      {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409      {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500      {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/apps [post]
func (h *AdminHandler) CreateApp(c *gin.Context) {
	var req datastructure.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	reqApp := &datastructure.Application{
		Name:        req.Name,
		Description: req.Description,
	}

	app, err := h.app_service.CreateApplication(c.Request.Context(), reqApp)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respApp := datastructure.Application{
		ID:          app.ID,
		TenantID:    app.TenantID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   util.ConvertTimeToTimeStamp(&app.CreatedAt),
		UpdatedAt:   util.ConvertTimeToTimeStamp(&app.UpdatedAt),
		DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(app.DeletedAt),
	}

	h.Success(c, respApp)
}

// GetApp godoc
// @Summary      取得指定應用程式詳細資料
// @Description  取得指定應用程式詳細資料
// @Tags         Admin/Application
// @Produce      json
// @Param        app_id  path      string  true  "應用程式 ID"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.Application}  "成功回應，包含應用程式詳細資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/apps/{app_id} [get]
func (h *AdminHandler) GetApp(c *gin.Context) {
	appID := c.Param("app_id")

	app, err := h.app_service.GetApplicationByID(c.Request.Context(), appID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respApp := datastructure.Application{
		ID:          app.ID,
		TenantID:    app.TenantID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   util.ConvertTimeToTimeStamp(&app.CreatedAt),
		UpdatedAt:   util.ConvertTimeToTimeStamp(&app.UpdatedAt),
		DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(app.DeletedAt),
	}

	h.Success(c, respApp)
}

// UpdateApp godoc
// @Summary      更新指定應用程式
// @Description  更新指定應用程式
// @Tags         Admin/Application
// @Produce      json
// @Param        app_id  path      string  true  "應用程式 ID"
// @Param        request  body  datastructure.Application  true  "更新應用程式資料"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/apps/{app_id} [put]
func (h *AdminHandler) UpdateApp(c *gin.Context) {
	var req datastructure.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	appID := c.Param("app_id")
	reqApp := &datastructure.Application{
		ID:          appID,
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.app_service.UpdateApplicationByID(c.Request.Context(), appID, reqApp)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// DeleteApp godoc
// @Summary      刪除指定應用程式
// @Description  刪除指定應用程式
// @Tags         Admin/Application
// @Produce      json
// @Param        app_id  path      string  true  "應用程式 ID"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/apps/{app_id} [delete]
func (h *AdminHandler) DeleteApp(c *gin.Context) {
	appID := c.Param("app_id")

	err := h.app_service.DeleteApplicationByID(c.Request.Context(), appID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// GetApps godoc
// @Summary      取得所有應用程式
// @Description  取得所有應用程式
// @Tags         Admin/Application
// @Produce      json
// @Success      200     {object}  datastructure.BaseResponse{data=[]datastructure.Application}  "成功回應，包含所有應用程式陣列"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/apps [get]
func (h *AdminHandler) GetApps(c *gin.Context) {
	apps, err := h.app_service.GetApplications(c.Request.Context())
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respApps := make([]*datastructure.Application, 0)
	for _, app := range apps {
		respApps = append(respApps, &datastructure.Application{
			ID:          app.ID,
			TenantID:    app.TenantID,
			Name:        app.Name,
			Description: app.Description,
			CreatedAt:   util.ConvertTimeToTimeStamp(&app.CreatedAt),
			UpdatedAt:   util.ConvertTimeToTimeStamp(&app.UpdatedAt),
			DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(app.DeletedAt),
		})
	}

	h.Success(c, respApps)
}

// CreateAppAPIKey godoc
// @Summary      建立應用程式 API 密鑰
// @Description  建立應用程式 API 密鑰
// @Tags         Admin/Application
// @Produce      json
// @Param        app_id  path      string  true  "應用程式 ID"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.ApplicationAPIKey}  "成功回應，包含新應用程式 API 密鑰"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/apps/{app_id}/api_keys [post]
func (h *AdminHandler) CreateAppAPIKey(c *gin.Context) {
	appID := c.Param("app_id")

	apiKey, err := h.app_service.CreateApplicationAPIKey(c.Request.Context(), appID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respAppAPIKey := datastructure.ApplicationAPIKey{
		ID:            apiKey.ID,
		ApplicationID: appID,
		APIKey:        apiKey.APIKey,
		CreatedAt:     util.ConvertTimeToTimeStamp(&apiKey.CreatedAt),
		UpdatedAt:     util.ConvertTimeToTimeStamp(&apiKey.UpdatedAt),
		DeletedAt:     util.ConvertGormDeletedAtToTimeStamp(apiKey.DeletedAt),
	}

	h.Success(c, respAppAPIKey)
}

// DeleteAppAPIKey godoc
// @Summary      刪除應用程式 API 密鑰
// @Description  刪除應用程式 API 密鑰
// @Tags         Admin/Application
// @Produce      json
// @Param        app_id  path      string  true  "應用程式 ID"
// @Param        api_key_id  path      string  true  "應用程式 API 密鑰 ID"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/apps/{app_id}/api-keys/{api_key_id} [delete]
func (h *AdminHandler) DeleteAppAPIKey(c *gin.Context) {
	appID := c.Param("app_id")
	apiKeyID := c.Param("api_key_id")

	err := h.app_service.DeleteApplicationAPIKey(c.Request.Context(), appID, apiKeyID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// CreateEvent godoc
// @Summary      建立新事件
// @Description  建立新事件
// @Tags         Admin/Event
// @Accept       json
// @Produce      json
// @Param        request  body  datastructure.Event  true  "新增事件資料"
// @Success      200      {object}  datastructure.BaseResponse{data=datastructure.Event}  "成功回應，包含新事件資料"
// @Failure      400      {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401      {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403      {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404      {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409      {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500      {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/events [post]
func (h *AdminHandler) CreateEvent(c *gin.Context) {
	var req datastructure.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	reqEvent := &datastructure.Event{
		ApplicationID: req.ApplicationID,
		PlatformID:    req.PlatformID,
		Name:          req.Name,
		Description:   req.Description,
		IsActive:      req.IsActive,
	}

	event, err := h.event_service.CreateEvent(c.Request.Context(), reqEvent)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respEvent := datastructure.Event{
		ID:            event.ID,
		ApplicationID: event.ApplicationID,
		PlatformID:    event.PlatformID,
		Name:          event.Name,
		Description:   event.Description,
		IsActive:      event.IsActive,
		CreatedAt:     util.ConvertTimeToTimeStamp(&event.CreatedAt),
		UpdatedAt:     util.ConvertTimeToTimeStamp(&event.UpdatedAt),
		DeletedAt:     util.ConvertGormDeletedAtToTimeStamp(event.DeletedAt),
	}

	h.Success(c, respEvent)
}

// GetEvent godoc
// @Summary      取得指定事件詳細資料
// @Description  取得指定事件詳細資料
// @Tags         Admin/Event
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.Event}  "成功回應，包含應用程式詳細資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/events/{event_id} [get]
func (h *AdminHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("event_id")

	event, err := h.event_service.GetEventByID(c.Request.Context(), eventID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respEvent := datastructure.Event{
		ID:            event.ID,
		ApplicationID: event.ApplicationID,
		PlatformID:    event.PlatformID,
		Name:          event.Name,
		Description:   event.Description,
		IsActive:      event.IsActive,
		CreatedAt:     util.ConvertTimeToTimeStamp(&event.CreatedAt),
		UpdatedAt:     util.ConvertTimeToTimeStamp(&event.UpdatedAt),
		DeletedAt:     util.ConvertGormDeletedAtToTimeStamp(event.DeletedAt),
	}

	h.Success(c, respEvent)
}

// UpdateEvent godoc
// @Summary      更新指定事件
// @Description  更新指定事件
// @Tags         Admin/Event
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Param        request  body  datastructure.Event  true  "更新事件資料"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/events/{event_id} [put]
func (h *AdminHandler) UpdateEvent(c *gin.Context) {
	var req datastructure.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	eventID := c.Param("event_id")
	appID := c.Param("app_id")
	reqEvent := &datastructure.Event{
		ID:            eventID,
		ApplicationID: appID,
		PlatformID:    req.PlatformID,
		Name:          req.Name,
		Description:   req.Description,
		IsActive:      req.IsActive,
	}

	err := h.event_service.UpdateEventByID(c.Request.Context(), eventID, reqEvent)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// DeleteEvent godoc
// @Summary      刪除指定事件
// @Description  刪除指定事件
// @Tags         Admin/Event
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/events/{event_id} [delete]
func (h *AdminHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("event_id")

	err := h.event_service.DeleteEventByID(c.Request.Context(), eventID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// GetEvents godoc
// @Summary      取得所有事件
// @Description  取得所有事件
// @Tags         Admin/Event
// @Produce      json
// @Success      200     {object}  datastructure.BaseResponse{data=[]datastructure.Event}  "成功回應，包含所有事件陣列"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/events [get]
func (h *AdminHandler) GetEvents(c *gin.Context) {
	events, err := h.event_service.GetEvents(c.Request.Context())
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.Success(c, events)
}

// CreateEventFields godoc
// @Summary      建立新事件欄位
// @Description  建立新事件欄位
// @Tags         Admin/Event
// @Accept       json
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Param        request  body  datastructure.EventField  true  "新增事件欄位資料"
// @Success      200      {object}  datastructure.BaseResponse{data=datastructure.EventField}  "成功回應，包含新事件欄位資料"
// @Failure      400      {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401      {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403      {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404      {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409      {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500      {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/events/{event_id}/fields [post]
func (h *AdminHandler) CreateEventFields(c *gin.Context) {
	var req datastructure.CreateEventFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	reqField := &datastructure.EventField{
		Name:        req.Name,
		EventID:     req.EventID,
		DataType:    req.DataType,
		Description: req.Description,
		IsRequired:  req.IsRequired,
	}

	eventField, err := h.event_service.CreateEventField(c.Request.Context(), reqField)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respEventField := datastructure.EventField{
		ID:          eventField.ID,
		EventID:     eventField.EventID,
		Name:        eventField.Name,
		DataType:    eventField.DataType,
		IsRequired:  eventField.IsRequired,
		Description: eventField.Description,
		CreatedAt:   util.ConvertTimeToTimeStamp(&eventField.CreatedAt),
		UpdatedAt:   util.ConvertTimeToTimeStamp(&eventField.UpdatedAt),
		DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(eventField.DeletedAt),
	}

	h.Success(c, respEventField)
}

// GetEventField godoc
// @Summary      取得指定事件欄位
// @Description  取得指定事件欄位
// @Tags         Admin/Event
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Param        field_id  path      string  true  "欄位 ID"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.EventField}  "成功回應，包含應用程式詳細資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/apps/{app_id}/events/{event_id}/fields/{field_id} [get]
func (h *AdminHandler) GetEventField(c *gin.Context) {
	eventID := c.Param("event_id")
	fieldID := c.Param("field_id")

	field, err := h.event_service.GetEventFieldByEventIDAndID(c.Request.Context(), eventID, fieldID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respEventField := datastructure.EventField{
		ID:          field.ID,
		EventID:     field.EventID,
		Name:        field.Name,
		DataType:    field.DataType,
		IsRequired:  field.IsRequired,
		Description: field.Description,
		CreatedAt:   util.ConvertTimeToTimeStamp(&field.CreatedAt),
		UpdatedAt:   util.ConvertTimeToTimeStamp(&field.UpdatedAt),
		DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(field.DeletedAt),
	}

	h.Success(c, respEventField)
}

// UpdateEventField godoc
// @Summary      更新指定事件欄位
// @Description  更新指定事件欄位
// @Tags         Admin/Event
// @Produce      json
// @Param        event_id  path  string  true  "事件 ID"
// @Param        field_id  path  string  true  "欄位 ID"
// @Param        request  body  datastructure.EventField  true  "更新事件欄位資料"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/apps/{app_id}/events/{event_id}/fields/{field_id} [put]
func (h *AdminHandler) UpdateEventField(c *gin.Context) {
	var req datastructure.UpdateEventFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	eventID := c.Param("event_id")
	fieldID := c.Param("field_id")
	reqField := &datastructure.EventField{
		EventID:     eventID,
		Name:        req.Name,
		DataType:    req.DataType,
		Description: req.Description,
		IsRequired:  req.IsRequired,
	}

	err := h.event_service.UpdateEventFieldByEventIDAndID(c.Request.Context(), eventID, fieldID, reqField)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// DeleteEventField godoc
// @Summary      刪除指定事件欄位
// @Description  刪除指定事件欄位
// @Tags         Admin/Event
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Param        field_id  path      string  true  "欄位 ID"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/events/{event_id}/fields/{field_id} [delete]
func (h *AdminHandler) DeleteEventField(c *gin.Context) {
	eventID := c.Param("event_id")
	fieldID := c.Param("field_id")

	err := h.event_service.DeleteEventFieldByEventIDAndID(c.Request.Context(), eventID, fieldID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// GetEventFields godoc
// @Summary      取得指定事件欄位
// @Description  取得指定事件欄位
// @Tags         Admin/Event
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Success      200     {object}  datastructure.BaseResponse{data=[]datastructure.EventField}  "成功回應，包含應用程式詳細資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /admin/events/{event_id}/fields [get]
func (h *AdminHandler) GetEventFields(c *gin.Context) {
	eventID := c.Param("event_id")

	eventFields, err := h.event_service.GetEventFields(c.Request.Context(), eventID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respEventFields := make([]*datastructure.EventField, 0, len(eventFields))
	for _, eventField := range eventFields {
		respEventFields = append(respEventFields, &datastructure.EventField{
			ID:          eventField.ID,
			EventID:     eventField.EventID,
			Name:        eventField.Name,
			DataType:    eventField.DataType,
			IsRequired:  eventField.IsRequired,
			Description: eventField.Description,
			CreatedAt:   util.ConvertTimeToTimeStamp(&eventField.CreatedAt),
			UpdatedAt:   util.ConvertTimeToTimeStamp(&eventField.UpdatedAt),
			DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(eventField.DeletedAt),
		})
	}

	h.Success(c, respEventFields)
}
