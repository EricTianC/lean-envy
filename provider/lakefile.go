package provider

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pelletier/go-toml/v2"
)

/*
LakeFile 解析得到的数据类，

目前只处理 require 部分
按照 lake 的逻辑，所有的依赖项目将 clone 至 .lake/packages/ 下
*/
type LakeFile struct {
	IsExist      bool
	FileType     string        // LakeFile 类型：".lean" 或 ".toml" （带点号，不带引号)
	Requirements []Requirement // TODO: 创建依赖对应的数据类型
}

var lakePackagesPath = filepath.Join(".lake", "packages")

/*
	Requirement

lakefile 中的单个依赖项

Lean 的 Lake 配置文件中 require 命令的基本语法如下：

require ["<scope>" /] <pkg-name> [@ <version>]

	[from <source>] [with <options>]

from 从句告知 Lake 依赖的地址。没有 from 从句，Lake 将从默认注册表（例如 Reservoir）中查找包，并使用获得的信息下载指定 version 的包。
可选的 scope 用来消除同名包的歧义。在 Reservoir 中，scope 部分是包的所有者（例如，@leanprover/doc-gen4 中的 leanprover）。
*/
type Requirement struct {
	Scope, PkgName, Version, Source string
	Options                         string // TODO: decide on this field's type
}

func fileExist(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func findLakeFile(baseDir string) (string, error) {
	// 优先检查 .lean 格式
	leanLakeFileExists, err := fileExist(filepath.Join(baseDir, "lakefile.lean"))
	if err != nil {
		return "", err
	}
	if leanLakeFileExists {
		return ".lean", nil
	}

	tomlLakeFileExists, err := fileExist(filepath.Join(baseDir, "lakefile.toml"))
	if err != nil {
		return "", err
	}
	if tomlLakeFileExists {
		return ".toml", nil
	}

	return "", nil
}

func handleLeanLakeFileRequirements(baseDir string) ([]Requirement, error) {
	/* 格式为：
	require ["<scope>" /] <pkg-name> [@ <version>]
	  [from <source>] [with <options>]
	示例:
	require "leanprover-community" / "mathlib"
	或
	require "leanprover-community" / "mathlib" @ "git#20c73142afa995ac9c8fb80a9bb585a55ca38308"
	*/
	file, err := os.ReadFile(filepath.Join(baseDir, "lakefile.lean"))
	if err != nil {
		return nil, err
	}
	content := string(file)
	return parseLeanLakeFileRequirements(content)
}

func parseLeanLakeFileRequirements(content string) ([]Requirement, error) {
	// 去除 -- 开头的注释行 TODO: 更加准确的注释判断逻辑，可能需要手写 Parser 了
	lines := strings.Split(content, "\n")
	var newLines []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			newLines = append(newLines, line)
		}
	}
	content = strings.Join(newLines, "\n")

	var requirements []Requirement

	re, err := regexp.Compile(`require(\s*"([\w-]*)"\s*/)?\s*"([\w-]*)"(\s*@\s*(git\s*)?"([\w\.#-]*)")?`) // TODO: 抛弃正则仙人写法
	if err != nil {
		return nil, err
	}
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		// fmt.Println(strings.Join(match, "😊"))
		requirements = append(requirements, Requirement{Scope: match[2], PkgName: match[3], Version: match[6]})
	}
	return requirements, nil
}

func handleTomlLakeFileRequirements(baseDir string) ([]Requirement, error) {
	// return []string{}, nil
	file, err := os.ReadFile(filepath.Join(baseDir, "lakefile.toml"))
	if err != nil {
		return nil, err
	}
	content := string(file)
	return parseTomlLakeFileRequirements(content)
}

type tomlRequire struct {
	Name    string            `toml:"name"`
	Scope   string            `toml:"scope,omitempty"`
	Version string            `toml:"version,omitempty"`
	Options map[string]string `toml:"options,omitempty"`
	Path    string            `toml:"path,omitempty"`
	Git     string            `toml:"git,omitempty"`
	Rev     string            `toml:"rev,omitempty"`
	SubDir  string            `toml:"subDir,omitempty"`
}

type tomlLakeFile struct {
	Requires []tomlRequire `toml:"require"`
}

func parseTomlLakeFileRequirements(content string) ([]Requirement, error) {
	// lines := strings.Split(content, "\n")
	var lakeFile tomlLakeFile
	toml.Unmarshal([]byte(content), &lakeFile)
	var requirements []Requirement
	for _, item := range lakeFile.Requires {
		if item.Rev != "" {
			requirements = append(requirements, Requirement{item.Scope, item.Name, item.Rev, item.Git, ""})
		} else {
			requirements = append(requirements, Requirement{item.Scope, item.Name, item.Version, item.Git, ""})
		}
	}
	return requirements, nil
}

func GetLakeFile(baseDir string) (LakeFile, error) {
	fileType, err := findLakeFile(baseDir)
	if err != nil {
		return LakeFile{}, err
	}
	if fileType == "" {
		return LakeFile{}, fmt.Errorf("项目下未发现 lakeFile 文件")
	}
	var requirements []Requirement
	switch fileType {
	case ".lean":
		requirements, err = handleLeanLakeFileRequirements(baseDir)
		if err != nil {
			return LakeFile{IsExist: true, FileType: fileType}, err
		}
	case ".toml":
		requirements, err = handleTomlLakeFileRequirements(baseDir)
		if err != nil {
			return LakeFile{IsExist: true, FileType: fileType}, err
		}
	}
	return LakeFile{IsExist: true, FileType: fileType, Requirements: requirements}, nil
}

func errorBrief(err error) string {
	if len(err.Error()) > 50 {
		return fmt.Sprintf("%s...", err.Error()[0:50])
	} else {
		return err.Error()
	}
}

func (lakefile *LakeFile) solveRequirements(solvedMap map[Requirement]bool) error {
	for _, item := range lakefile.Requirements { // TODO: just go it
		if solvedMap[item] {
			continue
		}
		err := item.Solve()
		if err != nil {
			slog.Error("依赖项下载失败", "name", item.RepoString(), "error", errorBrief(err))
		} else {
			solvedMap[item] = true
			subLakefile, err := GetLakeFile(filepath.Join(lakePackagesPath, item.PkgName))
			slog.Debug("Solve subdependencies", "repo", item.RepoString(), "lakefile", subLakefile.Requirements)
			if err != nil {
				// return err
				slog.Error("处理依赖项失败", "err", errorBrief(err))
			}
			err = subLakefile.solveRequirements(solvedMap)
			if err != nil {
				// return err
				slog.Error("处理依赖项失败", "err", err)
			}
		}

	}
	return nil
}

func (lakefile *LakeFile) SolveRequirements() error {
	var solvedMap = make(map[Requirement]bool)
	err := lakefile.solveRequirements(solvedMap)
	// var justSolved = make(chan Requirement) // 伏笔
	return err
}

func (requirement *Requirement) Solve() error { // TODO: 递归发现 requirement
	r, err := git.PlainOpen(filepath.Join(lakePackagesPath, requirement.PkgName))
	if err == nil {
		slog.Debug("Repo already exists, skipping clone", "name", requirement.RepoString(), "path", filepath.Join(lakePackagesPath, requirement.PkgName))
	} else {
		mirrorUrl, err := ToGithubRepoMirrorUrl(requirement.RepoString())
		if err != nil {
			slog.Error("未查询到对应的 Git repo 镜像", "repo", requirement.RepoString())
		}
		slog.Debug("Cloning Repo: ", "name", requirement.RepoString(), "URL", mirrorUrl)
		r, err = git.PlainClone(filepath.Join(lakePackagesPath, requirement.PkgName), false, &git.CloneOptions{
			URL:      mirrorUrl,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	}
	rev := requirement.Version
	if rev != "" {
		rev, _ = strings.CutPrefix(rev, "git#")
		worktree, err := r.Worktree()
		if err != nil {
			return err
		}
		refHash, err := r.ResolveRevision(plumbing.Revision(rev))
		if err != nil {
			return err
		}
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash: *refHash,
		})
		if err != nil {
			return err
		}
	}

	// remotes, err := r.Remotes() // TODO: 修改远程库
	// if err != nil {
	// 	return err
	// }
	// slog.Info("Repo Solved", "remotes", remotes)
	slog.Debug("Repo Solved")
	return nil
}

func (requirement *Requirement) RepoString() string {
	return fmt.Sprintf("%s/%s", requirement.Scope, requirement.PkgName)
}
