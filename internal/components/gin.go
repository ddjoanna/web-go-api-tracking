package component

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
	shared "tracking-service/internal"

	log "github.com/sirupsen/logrus"
	swaggoFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RouteRegistrar interface {
	RegisterRoutes(*gin.Engine)
}

func NewRouter(config *shared.Config, routeRegistrars []RouteRegistrar) *gin.Engine {
	// 設定為 config 中的模式（可設定為 "debug" 或 "release"）
	gin.SetMode(config.GinMode)
	router := gin.Default()
	// 自訂 Recovery
	router.Use(customRecovery())

	// 不追蹤的路由清單
	notTraced := map[string][]string{
		"GET": {"/healthcheck"},
	}
	router.Use(customOtelMiddleware(config.OtlpServiceName, notTraced))

	// 設定 X-TRACE-ID Header
	router.Use(traceIDHeaderMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggoFiles.Handler))

	// 註冊所有的路由
	for _, registrar := range routeRegistrars {
		registrar.RegisterRoutes(router)
	}

	return router
}

func customRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 取得請求資訊
		req := c.Request
		headers := formatHeaders(req.Header)
		method := req.Method
		path := req.URL.Path
		proto := req.Proto
		userAgent := req.UserAgent()
		clientIP := c.ClientIP()

		// 組合簡易日誌訊息
		msg := fmt.Sprintf("%s - \"%s %s %s %s\" Headers: {%s}",
			clientIP,
			method,
			path,
			proto,
			userAgent,
			headers,
		)

		// 取得 panic 內容與堆疊
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("Panic recovered: %v\n", recovered))
		buf.Write(debug.Stack())

		// 用 logrus 輸出，包含請求資訊與 panic 堆疊
		log.WithContext(c.Request.Context()).WithFields(log.Fields{
			"request": msg,
			"error":   buf.String(),
		}).Error("Panic recovered during request")

		// 回傳 500 給客戶端
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Internal Server Error",
		})
	})
}

func formatHeaders(h http.Header) string {
	var b bytes.Buffer
	for k, v := range h {
		b.WriteString(fmt.Sprintf("%s: %s; ", k, strings.Join(v, ",")))
	}
	return b.String()
}

func traceIDHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		traceID := span.SpanContext().TraceID().String()
		if traceID != "" {
			c.Writer.Header().Set("X-TRACE-ID", traceID)
		}
		c.Next()
	}
}

func customOtelMiddleware(serviceName string, notTracedEndpoints map[string][]string) gin.HandlerFunc {
	// 產生過濾函式，判斷是否需要追蹤
	filterTraces := func(req *http.Request) bool {
		if paths, ok := notTracedEndpoints[req.Method]; ok {
			for _, path := range paths {
				if req.URL.Path == path {
					return false // 不追蹤此路由
				}
			}
		}
		return true // 其他路由追蹤
	}

	return otelgin.Middleware(serviceName,
		otelgin.WithFilter(filterTraces),
		otelgin.WithPropagators(propagation.TraceContext{}),
	)
}

func NewHttpServer(
	lc fx.Lifecycle,
	config *shared.Config,
	router *gin.Engine,
) *http.Server {
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.HttpPort),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("HTTP server panic recovered: %v", r)
				}
			}()
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("HTTP server start failed: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			log.Info("Shutting down HTTP server...")
			if err := srv.Shutdown(shutdownCtx); err != nil {
				log.Errorf("Error during server shutdown: %v", err)
				return err
			}
			log.Info("HTTP server shut down gracefully.")
			return nil
		},
	})

	return srv
}
