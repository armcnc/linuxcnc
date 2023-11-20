/**
 ******************************************************************************
 * @file    code.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package CodeService

import (
	"armcnc/framework/package/code"
	"armcnc/framework/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io/ioutil"
)

type responseSelect struct {
	Code []CodePackage.Data `json:"code"`
}

func Select(c *gin.Context) {

	returnData := responseSelect{}
	returnData.Code = make([]CodePackage.Data, 0)

	code := CodePackage.Init()
	returnData.Code = code.Select()

	Utils.Success(c, 0, "", returnData)
	return
}

type responseReadLine struct {
	Content string   `json:"content"`
	Line    []string `json:"lines"`
}

func ReadLine(c *gin.Context) {

	returnData := responseReadLine{}
	returnData.Line = make([]string, 0)

	fileName := c.DefaultQuery("file_name", "")
	if fileName == "" {
		Utils.Error(c, 10000, "", returnData)
		return
	}

	code := CodePackage.Init()
	returnData.Content = code.ReadContent(fileName)
	read := code.ReadLine(fileName)
	returnData.Line = read.Line

	Utils.Success(c, 0, "", returnData)
	return
}

type responseReadContent struct {
	Content string `json:"content"`
}

func ReadContent(c *gin.Context) {

	returnData := responseReadContent{}

	fileName := c.DefaultQuery("file_name", "")
	if fileName == "" {
		Utils.Error(c, 10000, "", returnData)
		return
	}

	code := CodePackage.Init()
	returnData.Content = code.ReadContent(fileName)

	Utils.Success(c, 0, "", returnData)
	return
}

type requestUpdateContent struct {
	FileName string `json:"file_name"`
	Content  string `json:"content"`
}

func UpdateContent(c *gin.Context) {

	requestJson := requestUpdateContent{}
	requestData, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(requestData, &requestJson)
	if err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	code := CodePackage.Init()
	update := code.UpdateContent(requestJson.FileName, requestJson.Content)
	if !update {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
	return
}
