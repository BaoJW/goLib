package baseApp

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HttpConfig struct {
	Host   string
	Port   string
	Cors   CorsConfig
	engine *gin.Engine
	ready  bool
}

type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           time.Duration
	AllowWildcards   bool
}

type HttpRoute struct {
	Service string
	Method  string
	Path    string
	Tags    string
	Remark  string
	Handle  gin.HandlerFunc
}

type HttpRouteGroup struct {
	BasePath    string
	Routes      []HttpRoute
	Middlewares []gin.HandlerFunc
	Sync        bool
}

// DefaultCorsConfig 默认跨域支持RESTFul风格
// REST表征状态转移，本质上就是一种Web的约束和规则，是一套方法论
// 有以下特点：
// 1.统一接口(GET、POST、PUT、DELETE)+资源(URI)
// 2.无状态
// 3.分层系统
// 4.可缓存性
func DefaultCorsConfig() CorsConfig {
	return CorsConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           time.Hour * 12,
		AllowWildcards:   false,
	}
}

func DefaultHttpConfig() HttpConfig {
	return HttpConfig{
		Host:   "0.0.0.0",
		Port:   "8080",
		Cors:   DefaultCorsConfig(),
		engine: nil,
		ready:  false,
	}
}

func (ba *BaseApp) InitHttp(config HttpConfig) {
	ba.http = config
	ba.http.ready = true
	ba.http.engine = gin.New()
	ba.http.engine.Use(gin.Recovery())
	ba.http.engine.Use(cors.New(cors.Config{
		AllowOrigins:     ba.http.Cors.AllowedOrigins,
		AllowMethods:     ba.http.Cors.AllowedMethods,
		AllowHeaders:     ba.http.Cors.AllowedHeaders,
		MaxAge:           ba.http.Cors.MaxAge,
		AllowCredentials: ba.http.Cors.AllowCredentials,
		AllowWildcard:    ba.http.Cors.AllowWildcards,
	}))
	// TODO 日志收集与上报
	// TODO 分布式链路追踪
	// TODO swagger文档自动生成
	ba.http.engine.GET("/health", func(ctx *gin.Context) { ctx.AbortWithStatus(200) })
}

func (ba *BaseApp) RegisterHttpRoute(config HttpRouteGroup) {
	group := ba.http.engine.Group(config.BasePath)
	group.Use(config.Middlewares...)
	for _, route := range config.Routes {
		switch route.Method {
		case "GET":
			group.GET(route.Path, route.Handle)
		case "POST":
			group.POST(route.Path, route.Handle)
		case "PUT":
			group.PUT(route.Path, route.Handle)
		case "DELETE":
			group.DELETE(route.Path, route.Handle)
		case "OPTIONS":
			group.OPTIONS(route.Path, route.Handle)
		case "HEAD":
			group.HEAD(route.Path, route.Handle)
		case "PATCH":
			group.PATCH(route.Path, route.Handle)
		case "ANY":
			group.Any(route.Path, route.Handle)
		}
	}

}

func (ba *BaseApp) startHttp() {
	err := ba.http.engine.Run(fmt.Sprintf("%s:%d", ba.http.Host, ba.http.Port))
	if err != nil {
		// TODO add log
		panic(err)
	}
}
