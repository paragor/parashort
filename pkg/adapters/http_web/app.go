package http_web

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paragor/parashort/pkg/domain/parashort"
)

type WebServer struct {
	shortApp    *parashort.ParashortApp
	templateDir string
	assetsDir   string
}

func NewWebServer(shortApp *parashort.ParashortApp, templateDir string, assetsDir string) *WebServer {
	return &WebServer{shortApp: shortApp, templateDir: templateDir, assetsDir: assetsDir}
}

func (webServer WebServer) Run(ctx context.Context) error {
	//gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	webServer.registerAssets(engine)
	webServer.registerFrontend(engine)
	webServer.registerApi(engine)

	server := http.Server{
		Addr:    ":8000",
		Handler: engine,
	}

	errChan := make(chan error)
	go func() {
		errChan <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		server.Shutdown(ctx)
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}
