package http_web

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/paragor/parashort/pkg/domain/parashort"
	"github.com/paragor/parashort/pkg/domain/storage"
	"github.com/paragor/parashort/pkg/domain/user"
)

func (webServer WebServer) generateFrontendLink(key string) string {
	return "/z/" + key
}
func (webServer WebServer) registerFrontend(engine *gin.Engine) {
	engine.LoadHTMLGlob(path.Join(webServer.templateDir, "/*"))
	web := engine.Group("/")
	web.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	web.GET("/z/:short_url", func(c *gin.Context) {
		var (
			key        string
			text       string
			errMsg     string
			isNotFound bool
			url        string
		)
		key = c.Param("short_url")
		text, err := webServer.shortApp.LoadText(key)

		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				isNotFound = true
			} else {
				errMsg = err.Error()
			}
		}
		if len(key) > 0 {
			url = "http://" + c.Request.Host + webServer.generateFrontendLink(key)
		}

		c.HTML(200, "show.html", gin.H{
			"key":        key,
			"text":       text,
			"errMsg":     errMsg,
			"isNotFound": isNotFound,
			"url":        url,
		})
	})
	web.POST("/save", func(c *gin.Context) {
		var request = struct {
			Text        string `form:"text" binding:"required"`
			RequiredKey string `form:"required_key"`
		}{}
		err := c.Bind(&request)
		if err != nil {
			c.HTML(200, "show.html", gin.H{
				"errMsg": fmt.Sprintf("cant parse form: %s", err.Error()),
			})
			return
		}

		saveItem := parashort.SaveItem{
			RequiredKey: request.RequiredKey,
			Text:        request.Text,
			User:        user.NewUserInfo(net.ParseIP(c.ClientIP())),
			TTL:         0,
		}
		result, err := webServer.shortApp.SaveText(saveItem)
		if err != nil {
			c.HTML(200, "show.html", gin.H{
				"errMsg": fmt.Sprintf("cant parse form: %s", err.Error()),
			})
			return
		}
		c.Redirect(http.StatusFound, webServer.generateFrontendLink(result.Key))
	})

	engine.GET("/list", func(c *gin.Context) {
		keys, err := webServer.shortApp.List()
		if err != nil {
			c.HTML(200, "list.html", gin.H{
				"errMsg": fmt.Sprintf("cant parse form: %s", err.Error()),
			})

			return
		}
		c.HTML(200, "list.html", gin.H{
			"keys": keys,
		})
	})

}
