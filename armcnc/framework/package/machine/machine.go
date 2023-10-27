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
	"github.com/djherbis/times"
	"os"
	"sort"
	"strings"
	"time"
)

type Machine struct {
	Path string `json:"path"`
}

type Data struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Describe     string    `json:"describe"`
	Version      string    `json:"version"`
	Control      int       `json:"control"`
	Coordinate   string    `json:"coordinate"`
	Increments   string    `json:"increments"`
	LinearUnits  string    `json:"linear_units"`
	AngularUnits string    `json:"angular_units"`
	Time         time.Time `json:"-"`
}

type INI struct {
	Emc struct {
		Machine string `ini:"MACHINE"`
		Debug   string `ini:"DEBUG"`
		Version string `ini:"VERSION"`
	} `ini:"EMC"`
	Display struct {
		Increments string `ini:"INCREMENTS"`
	} `ini:"DISPLAY"`
	Traj struct {
		Coordinates  string `ini:"COORDINATES"`
		LinearUnits  string `ini:"LINEAR_UNITS"`
		AngularUnits string `ini:"ANGULAR_UNITS"`
	} `ini:"TRAJ"`
}

type USER struct {
	Base struct {
		Name     string `ini:"NAME"`
		Describe string `ini:"DESCRIBE"`
		Control  int    `ini:"CONTROL"`
	} `ini:"BASE"`
}

func Init() *Machine {
	return &Machine{
		Path: Config.Get.Basic.Workspace + "/configs/",
	}
}

func (machine *Machine) Select() []Data {
	data := make([]Data, 0)

	files, err := os.ReadDir(machine.Path)
	if err != nil {
		return data
	}

	for _, file := range files {
		item := Data{}
		if file.IsDir() {
			item.Path = file.Name()
			timeData, _ := times.Stat(machine.Path + file.Name())
			item.Time = timeData.BirthTime()
			if strings.Contains(file.Name(), "default_") {
				item.Time = item.Time.Add(-10 * time.Minute)
			}
			ini := machine.GetIni(file.Name())
			if ini.Emc.Version != "" {
				item.Version = ini.Emc.Version
				item.Coordinate = ini.Traj.Coordinates
				item.Increments = ini.Display.Increments
				item.LinearUnits = ini.Traj.LinearUnits
				item.AngularUnits = ini.Traj.AngularUnits
				data = append(data, item)
			}
			user := machine.GetUser(file.Name())
			item.Name = user.Base.Name
			item.Describe = user.Base.Describe
			item.Control = user.Base.Control
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Time.After(data[j].Time)
	})

	return data
}

func (machine *Machine) GetUser(path string) USER {
	data := USER{}
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.user")
	if exists {
		userFile, err := IniUtils.Load(machine.Path + path + "/machine.user")
		if err == nil {
			err = IniUtils.MapTo(userFile, &data)
		}
	}
	return data
}

func (machine *Machine) GetIni(path string) INI {
	data := INI{}
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.ini")
	if exists {
		iniFile, err := IniUtils.Load(machine.Path + path + "/machine.ini")
		if err == nil {
			err = IniUtils.MapTo(iniFile, &data)
		}
	}
	return data
}

func (machine *Machine) GetIniContent(path string) string {
	content := ""
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.ini")
	if exists {
		contentByte, err := FileUtils.ReadFile(machine.Path + path + "/machine.ini")
		if err == nil {
			content = string(contentByte)
		}
	}

	return content
}

func (machine *Machine) UpdateIniContent(path string, content string) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.ini")
	if exists {
		write := FileUtils.WriteFile(content, machine.Path+path+"/machine.ini")
		if write == nil {
			status = true
		}
	}
	return status
}
