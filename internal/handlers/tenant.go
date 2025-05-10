package handler

import (
	shared "tracking-service/internal"
	datastructure "tracking-service/internal/datastructures"
	service "tracking-service/internal/services"
	util "tracking-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	BaseHandler
	tenant_service   *service.TenantService
	app_service      *service.ApplicationService
	platform_service *service.PlatformService
	event_service    *service.EventService
}

func NewTenantHandler(
	tenant_service *service.TenantService,
	app_service *service.ApplicationService,
	platform_service *service.PlatformService,
	event_service *service.EventService,
) *TenantHandler {
	return &TenantHandler{
		tenant_service:   tenant_service,
		app_service:      app_service,
		platform_service: platform_service,
		event_service:    event_service,
	}
}

// GetPlatforms godoc
// @Summary      取得平台列表
// @Description  取得所有平台的名稱列表
// @Tags         Tenant/Platform
// @Produce      json
// @Success      200  {object}  datastructure.BaseResponse{data=[]datastructure.Platform}  "成功回應，包含平台陣列"
// @Failure      400  {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401  {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403  {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404  {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409  {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500  {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /platforms [get]
func (h *TenantHandler) GetPlatforms(c *gin.Context) {
	platforms, err := h.platform_service.GetPlatforms(c.Request.Context())
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respPlatforms := make([]datastructure.Platform, 0, len(platforms))
	for _, platform := range platforms {
		respPlatforms = append(respPlatforms, datastructure.Platform{
			ID:        platform.ID,
			Name:      platform.Name,
			CreatedAt: util.ConvertTimeToTimeStamp(&platform.CreatedAt),
			UpdatedAt: util.ConvertTimeToTimeStamp(&platform.UpdatedAt),
			DeletedAt: util.ConvertGormDeletedAtToTimeStamp(platform.DeletedAt),
		})
	}

	h.Success(c, respPlatforms)
}

// GetApp godoc
// @Summary      取得應用程式詳細資料
// @Description  取得指定應用程式的詳細資料
// @Tags         Tenant/Application
// @Produce      json
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.Application}  "成功回應，包含應用程式詳細資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/profile [get]
func (h *TenantHandler) GetApp(c *gin.Context) {
	tenantID := c.GetString(string(shared.TenantIDKey))
	appID := c.GetString(string(shared.TenantApplicationIDKey))

	app, err := h.app_service.GetApplicationByTenantIDAndID(c.Request.Context(), tenantID, appID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respApp := datastructure.Application{
		TenantID:    app.TenantID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   util.ConvertTimeToTimeStamp(&app.CreatedAt),
		UpdatedAt:   util.ConvertTimeToTimeStamp(&app.UpdatedAt),
		DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(app.DeletedAt),
	}

	h.Success(c, respApp)
}

// CreateEvent godoc
// @Summary      建立事件
// @Description  建立新事件
// @Tags         Tenant/Event
// @Accept       json
// @Produce      json
// @Param        request  body      datastructure.CreateEventRequest true  "新增事件資料"
// @Success      200      {object}  datastructure.BaseResponse{data=datastructure.Event}  "成功回應，包含新事件資料"
// @Failure      400      {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401      {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403      {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404      {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409      {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500      {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/events [post]
func (h *TenantHandler) CreateEvent(c *gin.Context) {
	var req datastructure.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	req.ApplicationID = c.GetString(string(shared.TenantApplicationIDKey))
	appID := c.GetString(string(shared.TenantApplicationIDKey))
	eventReq := datastructure.Event{
		ApplicationID: appID,
		PlatformID:    req.PlatformID,
		Name:          req.Name,
		Description:   req.Description,
		IsActive:      req.IsActive,
	}

	event, err := h.event_service.CreateEvent(c.Request.Context(), &eventReq)
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
// @Summary      取得事件詳細資料
// @Description  取得指定事件的詳細資料
// @Tags         Tenant/Event
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.Event}  "成功回應，包含事件詳細資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/events/{event_id} [get]
func (h *TenantHandler) GetEvent(c *gin.Context) {
	applicationID := c.GetString(string(shared.TenantApplicationIDKey))
	eventID := c.Param("event_id")

	event, err := h.event_service.GetEventByTenant(c.Request.Context(), applicationID, eventID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	fields := make([]datastructure.EventField, 0, len(event.Fields))
	for _, field := range event.Fields {
		fields = append(fields, datastructure.EventField{
			ID:          field.ID,
			EventID:     field.EventID,
			Name:        field.Name,
			DataType:    field.DataType,
			IsRequired:  field.IsRequired,
			Description: field.Description,
			CreatedAt:   util.ConvertTimeToTimeStamp(&field.CreatedAt),
			UpdatedAt:   util.ConvertTimeToTimeStamp(&field.UpdatedAt),
			DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(field.DeletedAt),
		})
	}

	respEvent := datastructure.EventResponse{
		Event: datastructure.Event{
			ID:            event.ID,
			ApplicationID: event.ApplicationID,
			PlatformID:    event.PlatformID,
			Name:          event.Name,
			Description:   event.Description,
			IsActive:      event.IsActive,
			CreatedAt:     util.ConvertTimeToTimeStamp(&event.CreatedAt),
			UpdatedAt:     util.ConvertTimeToTimeStamp(&event.UpdatedAt),
			DeletedAt:     util.ConvertGormDeletedAtToTimeStamp(event.DeletedAt),
		},
		Fields: fields,
	}

	h.Success(c, respEvent)
}

// UpdateEvent godoc
// @Summary      更新事件
// @Description  更新指定事件
// @Tags         Tenant/Event
// @Produce      json
// @Param        event_id  path  string  true  "事件 ID"
// @Param        request  body  datastructure.UpdateEventRequest  true  "更新事件資料"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/events/{event_id} [put]
func (h *TenantHandler) UpdateEvent(c *gin.Context) {
	var req datastructure.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	eventID := c.Param("event_id")
	applicationID := c.GetString(string(shared.TenantApplicationIDKey))

	eventReq := datastructure.Event{
		ID:            eventID,
		ApplicationID: applicationID,
		PlatformID:    req.PlatformID,
		Name:          req.Name,
		Description:   req.Description,
		IsActive:      req.IsActive,
	}

	err := h.event_service.UpdateEventByTenant(c.Request.Context(), applicationID, &eventReq)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// DeleteEvent godoc
// @Summary      刪除事件
// @Description  刪除指定事件
// @Tags         Tenant/Event
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/events/{event_id} [delete]
func (h *TenantHandler) DeleteEvent(c *gin.Context) {
	applicationID := c.GetString(string(shared.TenantApplicationIDKey))
	eventID := c.Param("event_id")

	err := h.event_service.DeleteEventByTenant(c.Request.Context(), applicationID, eventID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// GetEvents godoc
// @Summary      取得事件列表
// @Description  取得指定應用程式的事件列表
// @Tags         Tenant/Event
// @Produce      json
// @Success      200     {object}  datastructure.BaseResponse{data=[]datastructure.Event}  "成功回應，包含事件陣列"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/events [get]
func (h *TenantHandler) GetEvents(c *gin.Context) {
	applicationID := c.GetString(string(shared.TenantApplicationIDKey))
	events, err := h.event_service.GetEventsByApplicationID(c.Request.Context(), applicationID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respEvents := make([]datastructure.EventResponse, 0, len(events))

	for _, event := range events {
		fields := make([]datastructure.EventField, 0, len(event.Fields))
		for _, field := range event.Fields {
			fields = append(fields, datastructure.EventField{
				ID:          field.ID,
				EventID:     field.EventID,
				Name:        field.Name,
				DataType:    field.DataType,
				IsRequired:  field.IsRequired,
				Description: field.Description,
				CreatedAt:   util.ConvertTimeToTimeStamp(&field.CreatedAt),
				UpdatedAt:   util.ConvertTimeToTimeStamp(&field.UpdatedAt),
				DeletedAt:   util.ConvertGormDeletedAtToTimeStamp(field.DeletedAt),
			})
		}
		respEvents = append(respEvents, datastructure.EventResponse{
			Event: datastructure.Event{
				ID:            event.ID,
				ApplicationID: event.ApplicationID,
				PlatformID:    event.PlatformID,
				Name:          event.Name,
				Description:   event.Description,
				IsActive:      event.IsActive,
				CreatedAt:     util.ConvertTimeToTimeStamp(&event.CreatedAt),
				UpdatedAt:     util.ConvertTimeToTimeStamp(&event.UpdatedAt),
				DeletedAt:     util.ConvertGormDeletedAtToTimeStamp(event.DeletedAt),
			},
			Fields: fields,
		})
	}

	h.Success(c, respEvents)
}

// CreateEventField godoc
// @Summary      建立事件欄位
// @Description  建立新事件欄位
// @Tags         Tenant/Event
// @Accept       json
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Param        request  body  datastructure.CreateEventFieldRequest  true  "新增事件欄位資料"
// @Success      200      {object}  datastructure.BaseResponse{data=datastructure.EventField}  "成功回應：新增事件欄位資料"
// @Failure      400      {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401      {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403      {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404      {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409      {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500      {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/events/{event_id}/fields [post]
func (h *TenantHandler) CreateEventField(c *gin.Context) {
	eventID := c.Param("event_id")
	var req datastructure.CreateEventFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	reqField := datastructure.EventField{
		Name:        req.Name,
		EventID:     eventID,
		DataType:    req.DataType,
		Description: req.Description,
		IsRequired:  req.IsRequired,
	}

	field, err := h.event_service.CreateEventField(c.Request.Context(), &reqField)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respEvent := datastructure.EventField{
		ID:          field.ID,
		EventID:     field.EventID,
		Name:        field.Name,
		DataType:    field.DataType,
		IsRequired:  field.IsRequired,
		Description: field.Description,
		CreatedAt:   util.ConvertTimeToTimeStamp(&field.CreatedAt),
		UpdatedAt:   util.ConvertTimeToTimeStamp(&field.UpdatedAt),
	}

	h.Success(c, respEvent)
}

// GetEventField godoc
// @Summary      取得事件欄位
// @Description  取得指定事件欄位
// @Tags         Tenant/Event
// @Produce      json
// @Param        event_id  path      string  true  "事件 ID"
// @Param        field_id  path      string  true  "欄位 ID"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.EventField}  "成功回應：應用程式詳細資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/events/{event_id}/fields/{field_id} [get]
func (h *TenantHandler) GetEventField(c *gin.Context) {
	applicationID := c.GetString(string(shared.TenantApplicationIDKey))
	eventID := c.Param("event_id")
	fieldID := c.Param("field_id")

	field, err := h.event_service.GetEventFieldByTenant(c.Request.Context(), applicationID, eventID, fieldID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respEvent := datastructure.EventField{
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

	h.Success(c, respEvent)
}

// UpdateEventField godoc
// @Summary      更新事件欄位
// @Description  更新指定事件欄位
// @Tags         Tenant/Event
// @Produce      json
// @Param        event_id  path  string  true  "事件 ID"
// @Param        field_id  path  string  true  "欄位 ID"
// @Param        request  body  datastructure.UpdateEventFieldRequest  true  "更新事件欄位資料"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/events/{event_id}/fields/{field_id} [put]
func (h *TenantHandler) UpdateEventField(c *gin.Context) {
	var req datastructure.UpdateEventFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	applicationID := c.GetString(string(shared.TenantApplicationIDKey))
	eventID := c.Param("event_id")
	fieldID := c.Param("field_id")

	reqField := datastructure.EventField{
		Name:        req.Name,
		EventID:     eventID,
		DataType:    req.DataType,
		Description: req.Description,
		IsRequired:  req.IsRequired,
	}

	err := h.event_service.UpdateEventFieldByTenant(c.Request.Context(), applicationID, eventID, fieldID, &reqField)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// DeleteEventField godoc
// @Summary      刪除事件欄位
// @Description  刪除指定事件欄位
// @Tags         Tenant/Event
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
// @Router       /tenant/events/{event_id}/fields/{field_id} [delete]
func (h *TenantHandler) DeleteEventField(c *gin.Context) {
	applicationID := c.GetString(string(shared.TenantApplicationIDKey))
	eventID := c.Param("event_id")
	fieldID := c.Param("field_id")

	err := h.event_service.DeleteEventFieldByTenant(c.Request.Context(), applicationID, eventID, fieldID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// CreateEventLog godoc
// @Summary      建立事件日誌
// @Description  建立新事件日誌
// @Tags         Tenant/Event
// @Produce      json
// @Param        event_id  path  string  true  "事件 ID"
// @Param        request  body  datastructure.CreateEventLogRequest  true  "新增事件日誌資料"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.EventLog}  "成功回應，包含新事件日誌資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/events/{event_id}/event_logs [post]
func (h *TenantHandler) CreateEventLog(c *gin.Context) {
	var req datastructure.CreateEventLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	appID := c.GetString(string(shared.TenantApplicationIDKey))
	eventID := c.Param("event_id")
	reqEventLog := datastructure.EventLog{
		ApplicationID: appID,
		SessionID:     req.SessionID,
		EventID:       eventID,
		PlatformID:    req.PlatformID,
		Properties:    req.Properties,
	}

	eventLog, err := h.event_service.CreateEventLog(c.Request.Context(), &reqEventLog)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	respEventLog := datastructure.EventLog{
		ApplicationID: eventLog.ApplicationID,
		SessionID:     eventLog.SessionID,
		EventID:       eventLog.EventID,
		PlatformID:    eventLog.PlatformID,
		Properties:    eventLog.Properties,
		CreatedAt:     util.ConvertTimeToTimeStamp(&eventLog.CreatedAt),
	}

	h.Success(c, respEventLog)
}

// CreateSession godoc
// @Summary      建立會話
// @Description  建立新會話
// @Tags         Tenant/Session
// @Produce      json
// @Param        event_id  path  string  true  "事件 ID"
// @Param        request  body  datastructure.CreateSessionRequest  true  "新增會話資料"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.Session}  "成功回應，包含新會話資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/events/{event_id}/sessions [post]
func (h *TenantHandler) CreateSession(c *gin.Context) {
	var req datastructure.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	appID := c.GetString(string(shared.TenantApplicationIDKey))
	reqSession := datastructure.Session{
		ApplicationID: appID,
		PlatformID:    req.PlatformID,
		SessionKey:    req.SessionKey,
		UserID:        req.UserID,
		UserAgent:     req.UserAgent,
		IPAddress:     req.IPAddress,
		StartedAt:     req.StartedAt,
		EndedAt:       req.EndedAt,
	}

	session, err := h.app_service.CreateSession(c.Request.Context(), &reqSession)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	endedAt := util.ConvertTimeToTimeStamp(session.EndedAt)
	respSession := datastructure.Session{
		ApplicationID: session.ApplicationID,
		PlatformID:    session.PlatformID,
		SessionKey:    session.SessionKey,
		UserID:        session.UserID,
		UserAgent:     session.UserAgent,
		IPAddress:     session.IPAddress,
		StartedAt:     util.ConvertTimeToTimeStamp(&session.StartedAt),
		EndedAt:       &endedAt,
		CreatedAt:     util.ConvertTimeToTimeStamp(&session.CreatedAt),
		UpdatedAt:     util.ConvertTimeToTimeStamp(&session.UpdatedAt),
		DeletedAt:     util.ConvertGormDeletedAtToTimeStamp(session.DeletedAt),
	}
	h.Success(c, respSession)
}

// GetSession godoc
// @Summary      取得會話詳細資料
// @Description  取得指定會話的詳細資料
// @Tags         Tenant/Session
// @Produce      json
// @Param        session_id  path      string  true  "會話 ID"
// @Success      200     {object}  datastructure.BaseResponse{data=datastructure.Session}  "成功回應，包含會話詳細資料"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/sessions/{session_id} [get]
func (h *TenantHandler) GetSession(c *gin.Context) {
	applicationID := c.GetString(string(shared.TenantApplicationIDKey))
	sessionID := c.Param("session_id")

	session, err := h.app_service.GetSessionByApplicationIDAndID(c.Request.Context(), applicationID, sessionID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	endedAt := util.ConvertTimeToTimeStamp(session.EndedAt)
	respSession := datastructure.Session{
		ApplicationID: session.ApplicationID,
		PlatformID:    session.PlatformID,
		SessionKey:    session.SessionKey,
		UserID:        session.UserID,
		UserAgent:     session.UserAgent,
		IPAddress:     session.IPAddress,
		StartedAt:     util.ConvertTimeToTimeStamp(&session.StartedAt),
		EndedAt:       &endedAt,
		CreatedAt:     util.ConvertTimeToTimeStamp(&session.CreatedAt),
		UpdatedAt:     util.ConvertTimeToTimeStamp(&session.UpdatedAt),
		DeletedAt:     util.ConvertGormDeletedAtToTimeStamp(session.DeletedAt),
	}

	h.Success(c, respSession)
}

// UpdateSession godoc
// @Summary      更新會話
// @Description  更新指定會話
// @Tags         Tenant/Session
// @Produce      json
// @Param        session_id  path  string  true  "會話 ID"
// @Param        request  body  datastructure.UpdateSessionRequest  true  "更新會話資料"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/sessions/{session_id} [put]
func (h *TenantHandler) UpdateSession(c *gin.Context) {
	var req datastructure.UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.InvalidInputErrorResponse(c, err)
		return
	}

	applicationID := c.GetString(string(shared.TenantApplicationIDKey))
	sessionID := c.Param("session_id")
	err := h.app_service.UpdateSessionByApplicationIDAndID(c.Request.Context(), applicationID, sessionID, &req)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}

// DeleteSession godoc
// @Summary      刪除會話
// @Description  刪除會話
// @Tags         Tenant/Session
// @Produce      json
// @Param        session_id  path      string  true  "會話 ID"
// @Success      204     "成功，無內容回應"
// @Failure      400     {object}  datastructure.ErrorResponseWithCode "錯誤回應：無效請求"
// @Failure      401     {object}  datastructure.ErrorResponseWithCode "錯誤回應：未授權"
// @Failure      403     {object}  datastructure.ErrorResponseWithCode "錯誤回應：禁止訪問"
// @Failure      404     {object}  datastructure.ErrorResponseWithCode "錯誤回應：找不到資源"
// @Failure      409     {object}  datastructure.ErrorResponseWithCode "錯誤回應：重複鍵"
// @Failure      500     {object}  datastructure.ErrorResponseWithCode "錯誤回應：伺服器錯誤"
// @Router       /tenant/sessions/{session_id} [delete]
func (h *TenantHandler) DeleteSession(c *gin.Context) {
	applicationID := c.GetString(string(shared.TenantApplicationIDKey))
	sessionID := c.Param("session_id")

	err := h.app_service.DeleteSessionByApplicationIDAndID(c.Request.Context(), applicationID, sessionID)
	if err != nil {
		h.ErrorResponse(c, err)
		return
	}

	h.SuccessWithoutContent(c)
}
