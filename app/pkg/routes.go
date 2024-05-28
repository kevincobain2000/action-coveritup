package pkg

import (
	"embed"
	"net/http"

	instachart "github.com/kevincobain2000/instachart/pkg"
	"github.com/labstack/echo/v4"
)

const (
	DIST_DIR     = "frontend/dist"
	DOCS_URL     = "https://github.com/kevincobain2000/action-coveritup"
	FAVICON_FILE = "favicon.ico"
)

func SetupRoutes(e *echo.Echo, baseURL string, publicDir embed.FS, favicon embed.FS) {
	e.GET(baseURL+"", NewAssetsHandler(publicDir, "index.html").GetHTML)
	e.GET(baseURL+"readme", NewAssetsHandler(publicDir, "readme/index.html").GetHTML)

	e.GET(baseURL+"robots.txt", NewAssetsHandler(publicDir, "robots.txt").GetPlain)
	e.GET(baseURL+"ads.txt", NewAssetsHandler(publicDir, "ads.txt").GetPlain)

	// /favicon.ico
	e.GET(baseURL+FAVICON_FILE, NewFaviconHandler(&favicon).Get)

	// /upload to post new coverage for a branch
	e.POST(baseURL+"upload", NewUploadHandler().Post, HasAuthorizationHeader())

	// /destroy
	e.POST(baseURL+"destroy", NewDestroyHandler().Post, HasAuthorizationHeader())

	// /README.md to return markdown for embedings of badge and charts
	e.GET(baseURL+"api/readme", NewReadmeHandler().Get)

	// /badge to return badges
	e.GET(baseURL+"badge", NewBadgeHandler().Get)

	// /badge to return badges
	e.GET(baseURL+"progress", NewProgressHandler().Get)

	// /bar
	e.GET(baseURL+"bar", func(c echo.Context) error {
		img, err := instachart.NewBarChartHandler().Get(c)
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "", img)
	})

	// /line
	e.GET(baseURL+"line", func(c echo.Context) error {
		img, err := instachart.NewLineChartHandler().Get(c)
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "", img)
	})
	// /line
	e.GET(baseURL+"table", NewTableHandler().Get)

	// /chart to return charts
	e.GET(baseURL+"chart", NewChartHandler().Get)

	// /pr to return pr comment for a branch and base branch
	e.GET(baseURL+"pr", NewPRHandler().Get)
}
