package http_web

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/paragor/parashort/pkg/domain/parashort"
	"github.com/paragor/parashort/pkg/domain/storage"
	"github.com/paragor/parashort/pkg/domain/user"
)

func (webServer WebServer) registerApi(engine *gin.Engine) {
	v1 := engine.Group("/api/v1/")
	v1.GET("/url/:short_url", func(c *gin.Context) {
		text, err := webServer.shortApp.LoadText(c.Param("short_url"))
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				c.JSON(404, NewBadApiResponse(err))
			} else {
				c.JSON(400, NewBadApiResponse(err))
			}

			return
		}
		c.JSON(200, LoadApiResponse{
			ApiResponse: ApiResponse{Status: true},
			Text:        text,
		})
	})
	v1.GET("/list", func(c *gin.Context) {
		keys, err := webServer.shortApp.List()
		if err != nil {
			c.JSON(400, NewBadApiResponse(err))

			return
		}
		c.JSON(200, ListApiResponse{
			ApiResponse: ApiResponse{Status: true},
			Keys:        keys,
		})
	})

	v1.DELETE("/url/:short_url", func(c *gin.Context) {
		err := webServer.shortApp.Delete(c.Param("short_url"))
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.JSON(200, DeleteApiResponse{ApiResponse{Status: true}})
	})

	v1.POST("/save", func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, NewBadApiResponse(err))
			return
		}
		request := SaveRequest{}
		if err := json.Unmarshal(data, &request); err != nil {
			c.JSON(500, NewBadApiResponse(err))
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
			c.JSON(500, NewBadApiResponse(err))
			return
		}
		c.JSON(200, SaveApiResponse{
			ApiResponse: ApiResponse{Status: true},
			Key:         result.Key,
		})
	})
}
