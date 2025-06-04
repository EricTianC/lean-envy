package main

import (
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/erictianc/lean-envy/provider"
)

func main() {
	fmt.Println("Hello lean-envy v0.0.1")

	fmt.Println("检查 git 与 elan 的安装情况") // TODO: Turning to gogit soon
	// 检查 git, elan
	path, err := exec.LookPath("git")
	if err != nil {
		fmt.Println("git 未安装: ", err)
		fmt.Println("可访问：https://pc.qq.com/detail/13/detail_22693.html (不要下载电脑管家，点 '直接下载')")
		return
	}
	fmt.Println("git: \t", path)

	path, err = exec.LookPath("elan")
	if err != nil {
		fmt.Println("elan 未安装: ", err)
		fmt.Println("可访问：https://s3.jcloud.sjtu.edu.cn/899a892efef34b1b944a19981040f55b-oss01/elan/elan/releases/download/mirror_clone_list.html，下载最新版本")
		return
	}
	fmt.Println("elan: \t", path)

	elanInfo, err := provider.CheckElan()
	if err != nil {
		fmt.Println("检查 elan 环境失败: ", err)
	}
	fmt.Println(elanInfo.ToString())

	lakeFile, err := provider.GetLakeFile(".")
	if err != nil {
		fmt.Println("解析 lakefile 失败", err)
	}
	fmt.Println(lakeFile)

	err = lakeFile.SolveRequirements()
	if err != nil {
		slog.Error("处理依赖项失败", "err", err)
	}
	slog.Info("依赖项处理完成")
}
