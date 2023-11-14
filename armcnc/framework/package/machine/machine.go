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
	"strconv"
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

				user := machine.GetUser(file.Name())
				item.Name = user.Base.Name
				item.Describe = user.Base.Describe
				item.Control = user.Base.Control
				data = append(data, item)
			}
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

func (machine *Machine) UpdateUser(path string, data JsonUSER) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.user")
	if exists {
		iniFile, err := IniUtils.Load(machine.Path + path + "/machine.user")
		if err == nil {
			iniFile.Section("BASE").Key("NAME").SetValue(data.Base.Name)
			iniFile.Section("BASE").Key("DESCRIBE").SetValue(data.Base.Describe)
			iniFile.Section("BASE").Key("CONTROL").SetValue(strconv.Itoa(data.Base.Control))
			iniFile.Section("HANDWHEEL").Key("X_VELOCITY").SetValue(data.HandWheel.XVelocity)
			iniFile.Section("HANDWHEEL").Key("Y_VELOCITY").SetValue(data.HandWheel.YVelocity)
			iniFile.Section("HANDWHEEL").Key("Z_VELOCITY").SetValue(data.HandWheel.ZVelocity)
			iniFile.Section("HANDWHEEL").Key("A_VELOCITY").SetValue(data.HandWheel.AVelocity)
			iniFile.Section("HANDWHEEL").Key("B_VELOCITY").SetValue(data.HandWheel.BVelocity)
			iniFile.Section("HANDWHEEL").Key("C_VELOCITY").SetValue(data.HandWheel.CVelocity)
			err = IniUtils.SaveTo(iniFile, machine.Path+path+"/machine.user")
			if err == nil {
				status = true
			}
		}
	}
	return status
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

func (machine *Machine) UpdateIni(path string, data INI) bool {
	status := false
	exists, _ := FileUtils.PathExists(machine.Path + path + "/machine.ini")
	if exists {
		iniFile, err := IniUtils.Load(machine.Path + path + "/machine.ini")
		if err == nil {
			err := IniUtils.ReflectFrom(iniFile, data)
			if err == nil {
				err = IniUtils.SaveTo(iniFile, machine.Path+path+"/machine.ini")
				if err == nil {
					status = true
				}
			}
		}
	}
	return status
}
