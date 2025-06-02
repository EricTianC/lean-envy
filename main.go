package main

import (
	"fmt"
	"os/exec"

	"github.com/erictianc/lean-envy/provider"
)

func main() {
	fmt.Println("Hello lean-envy v0.0.1")

	// 检查 git, elan
	path, err := exec.LookPath("git")
	if err != nil {
		fmt.Println("git 未安装: ", err)
		return
	}
	fmt.Println("git: \t", path)

	path, err = exec.LookPath("elan")
	if err != nil {
		fmt.Println("elan 未安装: ", err)
		return
	}
	fmt.Println("elan: \t", path)

	elanInfo, err := provider.CheckElanInfo()
	if err != nil {
		fmt.Println("检查 elan 环境失败: ", err)
	}
	fmt.Println(elanInfo.ToString())
}
