package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/ini.v1"
)

// VARS
var workDir string
var useSystemIcons bool
var catShow bool
var catPrefix string
var catSuffix string

func readINI() string {
	binPath, _ := os.Executable()
	// fmt.Println(filepath.Dir(binPath) + "/")
	workDir = filepath.Dir(binPath) + "/"
	cfg, err := ini.Load(workDir + "/config.ini")

	if err != nil {
		fmt.Printf("error loading ini file: \n%v", err)
		os.Exit(1)
	}

	configFilePath := cfg.Section("main").Key("configfile").String()
	useSystemIcons, _ = cfg.Section("icons").Key("usesystemicons").Bool()

	catShow, _ = cfg.Section("misc").Key("showcat").Bool()
	catPrefix = cfg.Section("misc").Key("catprefix").String()
	catSuffix = cfg.Section("misc").Key("catsuffix").String()
	// fmt.Println(configFilePath)
	// fmt.Println(useSystemIcons)

	return configFilePath
}

func getIcons(icon string) string {

	var iconWithPath string

	if !useSystemIcons {
		iconWithPath = workDir + "icons/" + icon
	} else {
		iconWithPath = icon
	}
	// fmt.Println(iconWithPath)
	return iconWithPath
}

func printFavList() {
	favList, err := os.Open(readINI())

	if err != nil {
		panic(err)
	}

	defer favList.Close()

	scanner := bufio.NewScanner(favList)

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Print(string(scanner.Text()))
		if !strings.HasPrefix(line, "#") {
			lineArr := strings.Split(line, ";")
			cmd := lineArr[0]
			cat := lineArr[1]
			cmdOpts := lineArr[2]
			path := lineArr[3]
			name := lineArr[4]
			icon := getIcons(lineArr[5])

			if !catShow {
				fmt.Printf("%s%s\x00icon\x1f%s\x1finfo\x1fexec;%s;%s;%s\n", cat, name, icon, cmd, cmdOpts, path)
			} else {
				fmt.Printf("%s%s%s%s\x00icon\x1f%s\x1finfo\x1fexec;%s;%s;%s\n", catPrefix, cat, catSuffix, name, icon, cmd, cmdOpts, path)
			}
		}
	}

	// fmt.Print(string(favList), err)
}

// ANCHOR - parseOptions
func parseOptions(options string) []string {

	optCount := len(strings.Split(options, ","))

	if optCount != 0 {
		optArr := strings.Split(options, ",")
		return optArr
	} else {
		return nil
	}
}

func main() {

	rofiArg := os.Getenv("ROFI_RETV")
	rofiInfo := strings.Split(os.Getenv("ROFI_INFO"), ";")
	// args := os.Args[1:]

	// rofi first starts with ROFI_RETV = 0
	if rofiArg == "0" {
		printFavList()
	}

	if rofiArg == "1" {
		if rofiInfo[0] == "exec" {

			var processArgs []string

			cmd := rofiInfo[1]
			cmdOpts := parseOptions(rofiInfo[2])
			path := rofiInfo[3]

			if len(cmdOpts) != 0 {
				processArgs = slices.Concat(cmdOpts, []string{path})
			} else {
				processArgs = []string{path}
			}

			process := exec.Command(cmd, processArgs...)
			process.Start()
			// fmt.Println(cmd, "\n", cmdOpts, "\n", processArgs, "\n", path+"\n")
		}
	}
	// selected input value
	// if rofiArg == "2" {
	// 	fmt.Println(rofiArg)
}
