package http_web

import (
	"github.com/gin-gonic/gin"
)

func (webServer WebServer) registerAssets(engine *gin.Engine) {
	//engine.StaticFS("/assets/", http.Dir(webServer.assetsDir))
}
