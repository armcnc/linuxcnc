/**
 ******************************************************************************
 * @file    machine.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MachinePackage

import (
	"armcnc/framework/config"
	"armcnc/framework/utils/file"
	"armcnc/framework/utils/ini"
)

type Machine struct {
	Path string `json:"path"`
}

type Data struct {
	EMC struct {
		MACHINE      string `ini:"MACHINE"`
		DESCRIBE     string `ini:"DESCRIBE"`
		CONTROL_TYPE int    `ini:"CONTROL_TYPE"`
		DEBUG        string `ini:"DEBUG"`
		VERSION      string `ini:"VERSION"`
	} `ini:"EMC"`
}

func Init() *Machine {
	return &Machine{
		Path: Config.Get.Basic.Workspace + "/configs/",
	}
}

func (machine *Machine) Select() []Data {
	data := make([]Data, 0)
	return data
}

func (machine *Machine) Get(name string) Data {
	data := Data{}
	exists, _ := FileUtils.PathExists(machine.Path + name + "/machine.ini")
	if exists {
		iniFile, err := IniUtils.Load(machine.Path + name + "/machine.ini")
		if err == nil {
			err = IniUtils.MapTo(iniFile, &data)
		}
	}
	return data
}
