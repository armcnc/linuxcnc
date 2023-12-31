/**
 ******************************************************************************
 * @file    backup.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package BackupService

import (
	"armcnc/framework/config"
	"armcnc/framework/package/backup"
	"armcnc/framework/utils"
	"armcnc/framework/utils/file"
	"github.com/gin-gonic/gin"
)

type responseSelect struct {
	Backup []BackupPackage.Data `json:"backup"`
}

func Select(c *gin.Context) {

	returnData := responseSelect{}
	returnData.Backup = make([]BackupPackage.Data, 0)

	backup := BackupPackage.Init()
	returnData.Backup = backup.Select()

	Utils.Success(c, 0, "", returnData)
	return
}

func Pack(c *gin.Context) {

	Type := c.DefaultQuery("type", "")
	if Type == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	backup := BackupPackage.Init()
	status := backup.Pack(Type)
	if !status {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}

func Restore(c *gin.Context) {

	fileName := c.DefaultQuery("file_name", "")
	if fileName == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	backup := BackupPackage.Init()

	zip := FileUtils.Unzip(backup.Path+fileName, Config.Get.Basic.Workspace+"/")
	if !zip {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}

func Delete(c *gin.Context) {

	fileName := c.DefaultQuery("file_name", "")
	if fileName == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	backup := BackupPackage.Init()
	status := backup.Delete(fileName)
	if !status {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}
