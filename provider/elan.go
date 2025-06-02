package provider

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type ElanInfo struct {
	// Path             string
	Version          string
	Toolchains       []string
	DefaultToolchain string
}

// var _defaultElanInfo = ElanInfo{}

func CheckElan() (ElanInfo, error) {
	// _defaultElanInfo = ElanInfo{Path: "elan"}
	// return _defaultElanInfo, nil
	cmd := exec.Command("elan", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ElanInfo{}, err
	}
	var elanVersion = out.String()
	elanVersion = strings.TrimSpace(elanVersion)
	out.Reset()

	cmd = exec.Command("elan", "toolchain", "list")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ElanInfo{}, err
	}
	elanToolchainsRaw := strings.Fields(out.String())
	var elanToolchains []string
	var elanDefaultToolchain = ""
	for index, item := range elanToolchainsRaw {
		if item == "(default)" {
			elanDefaultToolchain = elanToolchainsRaw[index-1]
		} else {
			elanToolchains = append(elanToolchains, item)
		}
	}

	return ElanInfo{Version: elanVersion, Toolchains: elanToolchains, DefaultToolchain: elanDefaultToolchain}, nil
}

func (info *ElanInfo) ToString() string {
	var toolchainString string
	for _, toolchain := range info.Toolchains {
		if toolchain == info.DefaultToolchain {
			toolchainString = fmt.Sprint(toolchainString, toolchain, " (默认)\n")
		} else {
			toolchainString = fmt.Sprint(toolchainString, toolchain, "\n")
		}
	}
	return fmt.Sprintf("%s\n%s", info.Version, toolchainString)
}

// func GetToolchains() ([]string, error) {
// 	cmd := exec.Command("elan", "toolchain", "list")
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	if err := cmd.Run(); err != nil {
// 		return nil, err
// 	}

// 	outString := out.String()
// 	return strings.Fields(outString), nil
// }
