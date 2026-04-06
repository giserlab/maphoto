package view

import (
	"embed"
	"fmt"
	"maphoto/internal/app/web/assets"
	"maphoto/internal/util"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
)

func AssetsFinder(c echo.Context) error {
	fileName := c.Param("filename")
	assetDir := fmt.Sprintf("dist/%s/%s", "static", fileName)
	var fs embed.FS = assets.FS
	if strings.HasSuffix(fileName, ".js") || strings.HasSuffix(fileName, ".map") {
		c.Response().Header().Set("Content-Type", "application/javascript")
	} else if strings.HasSuffix(fileName, ".png") {
		c.Response().Header().Set("Content-Type", "image/x-icon")

	} else if strings.HasSuffix(fileName, ".ico") {
		c.Response().Header().Set("Content-Type", "image/x-icon")

	} else {
		c.Response().Header().Set("Content-Type", "text/css")
	}
	content, err := fs.ReadFile(assetDir)
	if err != nil {
		util.Logger.Error(err.Error())
		return c.String(http.StatusNotFound, err.Error())
	}
	c.Response().Write(content)
	return nil
}

func PageIndex(c echo.Context) error {
	t, _ := template.ParseFS(assets.FS, "dist/index.html")
	t.Execute(c.Response(), nil)
	return nil
}

func ViewPhoto(c echo.Context) error {
	name := c.Param("name")
	picPath := filepath.Join(util.ExcutePath(), "maphoto_data", "photo", name)
	content, err := os.ReadFile(picPath)
	if err != nil {
		util.Logger.Error(err.Error())
		return c.String(http.StatusNotFound, err.Error())
	}
	c.Response().Write(content)
	return nil
}

func ViewThumbnail(c echo.Context) error {
	name := c.Param("name")
	picPath := filepath.Join(util.ExcutePath(), "maphoto_data", "thumbnail", name)
	content, err := os.ReadFile(picPath)
	if err != nil {
		util.Logger.Error(err.Error())
		return c.String(http.StatusNotFound, err.Error())
	}
	c.Response().Write(content)
	return nil
}
