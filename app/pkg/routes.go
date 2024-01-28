package pkg

import (
	"embed"
	"fmt"
	"net/http"

	instachart "github.com/kevincobain2000/instachart/pkg"
	"github.com/labstack/echo/v4"
)

const (
	DIST_DIR     = "frontend/dist"
	DOCS_URL     = "https://github.com/kevincobain2000/action-coveritup"
	FAVICON_FILE = "favicon.ico"
	ROBOTS_FILE  = "robots.txt"
	ROBOTS_TXT   = `User-agent: *
Allow: /
Disallow: /upload
Disallow: /chart
Disallow: /badge
Disallow: /destroy
Disallow: /pr
Disallow: /readme`
)

func SetupRoutes(e *echo.Echo, baseURL string, publicDir embed.FS, favicon embed.FS) {

	e.GET(baseURL+"", func(c echo.Context) error {
		filename := fmt.Sprintf("%s/%s", DIST_DIR, "index.html")
		content, err := publicDir.ReadFile(filename)
		if err != nil {
			return c.String(http.StatusNotFound, "404")
		}
		return ResponseHTML(c, content)
	})

	// /robots.txt
	e.GET(baseURL+ROBOTS_FILE, NewRobotsHandler().Get)

	// /favicon.ico
	e.GET(baseURL+FAVICON_FILE, NewFaviconHandler(&favicon).Get)

	// /upload to post new coverage for a branch
	e.POST(baseURL+"upload", NewUploadHandler().Post, HasAuthorizationHeader())

	// /destroy
	e.POST(baseURL+"destroy", NewDestroyHandler().Post, HasAuthorizationHeader())

	// /README.md to return markdown for embedings of badge and charts
	e.GET(baseURL+"readme", NewReadmeHandler().Get)

	// /badge to return badges
	e.GET(baseURL+"badge", NewBadgeHandler().Get)

	// /bar
	e.GET(baseURL+"bar", func(c echo.Context) error {
		img, err := instachart.NewBarChartHandler().Get(c)
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "", img)
	})

	// /chart to return charts
	e.GET(baseURL+"chart", NewChartHandler().Get)

	// /pr to return pr comment for a branch and base branch
	e.GET(baseURL+"pr", NewPRHandler().Get)
}
