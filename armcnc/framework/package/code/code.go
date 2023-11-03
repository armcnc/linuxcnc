/**
 ******************************************************************************
 * @file    code.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package CodePackage

import (
	"armcnc/framework/config"
	"bufio"
	"github.com/djherbis/times"
	"github.com/goccy/go-json"
	"os"
	"strings"
	"time"
)

type Code struct {
	Path string `json:"path"`
}

type Data struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Describe string    `json:"describe"`
	Version  string    `json:"version"`
	Time     time.Time `json:"-"`
}

func Init() *Code {
	return &Code{
		Path: Config.Get.Basic.Workspace + "/files/",
	}
}

func (code *Code) Select() []Data {
	data := make([]Data, 0)

	files, err := os.ReadDir(code.Path)
	if err != nil {
		return data
	}

	for _, file := range files {
		item := Data{}
		if !file.IsDir() {
			item.Path = file.Name()
			timeData, _ := times.Stat(code.Path + file.Name())
			item.Time = timeData.BirthTime()
			if strings.Contains(file.Name(), "demo") {
				item.Time = item.Time.Add(-10 * time.Minute)
			}
			firstLine := code.ReadFirstLine(code.Path + file.Name())
			if firstLine.Version != "" {
				item.Name = firstLine.Name
				item.Describe = firstLine.Describe
				item.Version = firstLine.Version
				data = append(data, item)
			}
		}
	}

	return data
}

func (code *Code) ReadFirstLine(path string) Data {
	data := Data{}
	file, err := os.Open(path)
	if err != nil {
		return data
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return data
	}

	line = strings.TrimSpace(line)

	if len(line) > 0 && line[0] == '(' && line[len(line)-1] == ')' {
		jsonStr := line[1 : len(line)-1]
		err := json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			return data
		}

	}

	return data
}
