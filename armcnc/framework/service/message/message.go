/**
 ******************************************************************************
 * @file    message.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MessageService

import (
	"armcnc/framework/package/launch"
	"armcnc/framework/utils/socket"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"os/exec"
)

func Service(c *gin.Context) {

	client, _ := SocketUtils.SocketGrader.Upgrade(c.Writer, c.Request, nil)

	if !SocketUtils.SocketStruct.Status {
		SocketUtils.SocketStruct.User = make(map[*websocket.Conn]bool)
		SocketUtils.SocketStruct.Status = true
	}

	SocketUtils.SocketStruct.User[client] = true

	for {
		_, data, err := client.ReadMessage()
		if err != nil {
			break
		}

		jsonFormat := SocketUtils.SocketMessageFormat{}
		err = json.Unmarshal(data, &jsonFormat)
		if err == nil {
			if jsonFormat.Command != "" {
				if jsonFormat.Command == "desktop:device:restart" {
					go func() {
						launch := LaunchPackage.Init()
						launch.OnRestart()
					}()
				}
				if jsonFormat.Command == "desktop:device:shutdown" {
					cmd := exec.Command("shutdown", "-h", "now")
					cmd.Run()
				}
				SocketUtils.SendMessage(jsonFormat.Command, jsonFormat.Message, jsonFormat.Data)
			}
		}
	}
}
