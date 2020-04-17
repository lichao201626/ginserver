package main

import (
	// "ginserver/dao"
	"ginserver/middleware"
	"ginserver/resources"

	// "ginserver/services"
	"ginserver/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	// "fmt"
)

func main() {
	// Background services
	// services.CollectPriceData()

	// Utils
	util.LoadEnvVars()
	util.UseJSONLogFormat()

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(middleware.JSONLogMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID(middleware.RequestIDOptions{AllowSetting: false}))
	r.Use(middleware.ConnectMgo())
	r.Use(middleware.Auth())
	r.Use(middleware.CORS(middleware.CORSOptions{}))

	resources.NewAdminResource(r)
	resources.NewUserResource(r)
	resources.NewTokenResource(r)
	resources.NewPincodeResource(r)
	resources.NewPriceResource(r)
	resources.NewOrderResource(r)
	resources.NewPlanResource(r)
	resources.NewKeyResource(r)
	resources.NewWsResource(r)

	port := util.GetEnv("PORT", "8889")
	log.Info("Service starting on port " + port)

	r.Run(":" + port) // listen and serve
}
