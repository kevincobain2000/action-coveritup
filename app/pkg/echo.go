package pkg

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-echarts/statsview"
	"github.com/go-echarts/statsview/viewer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEcho(baseURL string, publicDir embed.FS, favicon embed.FS, cors string) *echo.Echo {
	if os.Getenv("PPROF_HOST") != "" && os.Getenv("PPROF_PORT") != "" {
		Logger().Info("pprof enabled and listening on: ", os.Getenv("PPROF_HOST")+":"+os.Getenv("PPROF_PORT"))
		addr := os.Getenv("PPROF_HOST") + ":" + os.Getenv("PPROF_PORT")
		viewer.SetConfiguration(viewer.WithTheme(viewer.ThemeWesteros), viewer.WithAddr(addr))
		mgr := statsview.New()
		_ = mgr
		go mgr.Start()
		// mgr.Stop()
	}

	e := echo.New()

	e.HTTPErrorHandler = HTTPErrorHandler
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: ltsv(),
	}))
	SetupRoutes(e, baseURL, publicDir, favicon)

	if cors != "" {
		SetupCors(e, cors)
	}

	return e
}

func SetupCors(e *echo.Echo, cors string) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:" + cors},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
}

func StartEcho(e *echo.Echo, host string, port string) {
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", host, port)))
}

// HTTPErrorResponse is the response for HTTP errors
type HTTPErrorResponse struct {
	Error interface{} `json:"error"`
}

// HTTPErrorHandler handles HTTP errors for entire application
func HTTPErrorHandler(err error, c echo.Context) {
	SetHeadersResponseJSON(c.Response().Header())
	code := http.StatusInternalServerError
	var message interface{}
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
	} else {
		message = err.Error()
	}

	Logger().Error(message)
	if code == http.StatusInternalServerError {
		message = "Internal Server Error"
	}
	if err = c.JSON(code, &HTTPErrorResponse{Error: message}); err != nil {
		Logger().Error(err)
	}
}

func ltsv() string {
	timeCustom := time.Now().Format("2006-01-02 15:04:05")
	var format string
	format += fmt.Sprintf("time:%s\t", timeCustom)
	format += "host:${remote_ip}\t"
	format += "xff:${header:x-forwarded-for}\t"
	format += "req:-\t"
	format += "status:${status}\t"
	format += "method:${method}\t"
	format += "uri:${uri}\t"
	format += "size:${bytes_out}\t"
	format += "referer:${referer}\t"
	format += "ua:${user_agent}\t"
	format += "reqtime_ns:${latency}\t"
	format += "cache:-\t"
	format += "runtime:-\t"
	format += "apptime:-\t"
	format += "vhost:${host}\t"
	format += "reqtime_human:${latency_human}\t"
	format += "x-request-id:${id}\t"
	format += "host:${host}\n"
	return format
}
