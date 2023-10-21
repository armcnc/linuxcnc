/**
 ******************************************************************************
 * @file    launch.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package LaunchPackage

import (
	"armcnc/framework/config"
	"armcnc/framework/utils/file"
	"github.com/gookit/color"
	"log"
	"os/exec"
)

type Launch struct {
	Path string `json:"path"`
}

type Data struct {
}

func Init() *Launch {
	return &Launch{
		Path: Config.Get.Basic.Workspace,
	}
}

func (launch *Launch) Start(machine string) {
	log.Println("[launch]: " + color.Info.Sprintf(machine))
	write := FileUtils.WriteFile("MACHINE_PATH="+machine, Config.Get.Basic.Workspace+"/.armcnc/environment")
	if write == nil {
		if machine != "" {
			cmd := exec.Command("systemctl", "restart", "armcnc_linuxcnc.service")
			cmd.Output()
			cmd = exec.Command("systemctl", "restart", "armcnc_launch.service")
			cmd.Output()
		} else {
			cmd := exec.Command("systemctl", "stop", "armcnc_launch.service")
			cmd.Output()
			cmd = exec.Command("systemctl", "stop", "armcnc_linuxcnc.service")
			cmd.Output()
		}
	} else {
		cmd := exec.Command("systemctl", "stop", "armcnc_launch.service")
		cmd.Output()
		cmd = exec.Command("systemctl", "stop", "armcnc_linuxcnc.service")
		cmd.Output()
	}
}

func (launch *Launch) Stop() {
	cmd := exec.Command("systemctl", "stop", "armcnc_launch.service")
	cmd.Output()
	cmd = exec.Command("systemctl", "stop", "armcnc_linuxcnc.service")
	cmd.Output()
}
