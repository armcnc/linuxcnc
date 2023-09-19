/**
 ******************************************************************************
 * @file    message.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MessageService

import (
	"armcnc/framework/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
				for user := range SocketUtils.SocketStruct.User {
					err := user.WriteMessage(websocket.TextMessage, data)
					if err != nil {
						user.Close()
						delete(SocketUtils.SocketStruct.User, user)
					}
				}
			}
		}
	}
}