/**
 ******************************************************************************
 * @file    router.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Service

import (
	"armcnc/framework/config"
	"armcnc/framework/service/backup"
	"armcnc/framework/service/config"
	"armcnc/framework/service/machine"
	"armcnc/framework/service/message"
	"armcnc/framework/service/plugin"
	"armcnc/framework/service/program"
	"armcnc/framework/service/upload"
	"armcnc/framework/service/version"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router() http.Handler {

	router := gin.Default()

	gin.SetMode(gin.DebugMode)

	router.Use(gin.Recovery())

	router.Use(cors.Default())

	router.Static("/touch", Config.Get.Basic.Workspace+"/touch/")

	router.Static("/assets", Config.Get.Basic.Workspace+"/touch/assets/")

	router.Static("/static", Config.Get.Basic.Workspace+"/touch/static/")

	router.Static("/monacoeditorwork", Config.Get.Basic.Workspace+"/touch/monacoeditorwork/")

	resources := router.Group("resources")
	{
		resources.Static("/programs", Config.Get.Basic.Workspace+"/programs/")

		resources.Static("/plugins", Config.Get.Basic.Workspace+"/plugins/")

		resources.Static("/backups", Config.Get.Basic.Workspace+"/backups/")

		resources.Static("/runtime", Config.Get.Basic.Runtime+"/")
	}

	message := router.Group("message")
	{
		message.GET("/service", MessageService.Service)
	}

	config := router.Group("config")
	{
		config.GET("/index", ConfigService.Index)
	}

	machine := router.Group("machine")
	{
		machine.GET("/select", MachineService.Select)

		machine.GET("/get", MachineService.Get)

		machine.POST("/update", MachineService.Update)

		machine.GET("/download", MachineService.Download)

		machine.GET("/delete", MachineService.Delete)

		machine.POST("/update/launch", MachineService.UpdateLaunch)

		machine.POST("/update/hal", MachineService.UpdateHal)

		machine.POST("/update/xml", MachineService.UpdateXml)

		machine.GET("/default", MachineService.Default)

		machine.POST("/upload", UploadService.UploadMachine)
	}

	program := router.Group("program")
	{
		program.GET("/select", ProgramService.Select)

		program.GET("/read/line", ProgramService.ReadLine)

		program.GET("/read/content", ProgramService.ReadContent)

		program.POST("/update/content", ProgramService.UpdateContent)

		program.GET("/delete", ProgramService.Delete)

		program.POST("/upload", UploadService.UploadProgram)
	}

	plugin := router.Group("plugin")
	{
		plugin.GET("/select", PluginService.Select)
	}

	settings := router.Group("settings")
	{
		settings.GET("/backup/select", BackupService.Select)

		settings.GET("/backup/pack", BackupService.Pack)

		settings.GET("/backup/restore", BackupService.Restore)

		settings.GET("/backup/delete", BackupService.Delete)

		settings.GET("/version/check", VersionService.Check)
	}

	return router
}
