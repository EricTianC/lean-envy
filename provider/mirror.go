package provider

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

// type Mirror interface {
// 	CheckAvailability() (bool, error)
// 	GitRepo(string) string
// }

/* require "leanprover-community" / "mathlib"
-- require "leanprover-community" / "mathlib" @ "git#20c73142afa995ac9c8fb80a9bb585a55ca38308"
require "leanprover-community" / "batteries" @ git "main"
require "leanprover-community" / "Qq" @ git "master"
require "leanprover-community" / "aesop" @ git "master"
require "leanprover-community" / "proofwidgets" @ git "v0.0.60" -- ProofWidgets should always be pinned to a specific version
  with NameMap.empty.insert `errorOnBuild
    "ProofWidgets not up-to-date. \
    Please run `lake exe cache get` to fetch the latest ProofWidgets. \
    If this does not work, report your issue on the Lean Zulip."
require "leanprover-community" / "importGraph" @ git "main"
require "leanprover-community" / "LeanSearchClient" @ git "main"
require "leanprover-community" / "plausible" @ git "main" */

// sjtu 镜像有点过于容易挂了

func gitCloneUrl(repoStr string) string {
	return fmt.Sprintf("https://gitclone.com/github.com/%s", repoStr)
}

// 第一步转换, 只储存特殊规则
var mapLeanPkgName = map[string]string{
	"leanprover-community/mathlib":      "leanprover-community/mathlib4",
	"leanprover-community/Qq":           "leanprover-community/quote4",
	"leanprover-community/proofwidgets": "leanprover-community/ProofWidgets4",
	"leanprover-community/importGraph":  "leanprover-community/import-graph",
	"leanprover/Cli":                    "leanprover/lean4-cli",
}

var stjuRepoMirrors = map[string]string{
	"leanprover-community/mathlib4":      "https://mirror.sjtu.edu.cn/git/lean4-packages/mathlib4",      // ✔️ ❓ stju 不定时抽风是这样的
	"leanprover-community/batteries":     "https://mirror.sjtu.edu.cn/git/lean4-packages/batteries",     // ✔️ ❓
	"leanprover-community/quote4":        "https://mirror.sjtu.edu.cn/git/lean4-packages/quote4",        // ✔️ ❓
	"leanprover-community/aesop":         "https://mirror.sjtu.edu.cn/git/lean4-packages/aesop",         // ✔️ ❓
	"leanprover-community/ProofWidgets4": "https://mirror.sjtu.edu.cn/git/lean4-packages/ProofWidgets4", // ✔️ ❓
	"leanprover-community/import-graph":  "https://mirror.sjtu.edu.cn/git/lean4-packages/import-graph",  // ✔️ ❓
	// "leanprover-community/LeanSearchClient": gitCloneUrl("leanprover-community/LeanSearchClient"),          // ✔️
	// "leanprover-community/plausible":        gitCloneUrl("leanprover-community/plausible"),                 // ✔️
	"leanprover/lean4-cli": "https://mirror.sjtu.edu.cn/git/lean4-packages/lean4-cli", // ❓
}

func CheckGithubRepoAvailibility(url string) bool {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	refs, err := rem.List(&git.ListOptions{
		PeelingOption: git.AppendPeeled,
	})

	slog.Debug("Checking Repo", "git-url", url, "refs", refs)

	return err == nil
}

func ToGithubRepoMirrorUrl(repoStr string) (string, error) {
	newRepoStr, ok := mapLeanPkgName[repoStr]
	if ok {
		repoStr = newRepoStr
	}

	// Check STJU mirror: 经常抽风
	url, ok := stjuRepoMirrors[repoStr]
	if ok && CheckGithubRepoAvailibility(url) {
		return url, nil
	}

	// Check GitClone: 不太稳定
	if CheckGithubRepoAvailibility(gitCloneUrl(repoStr)) {
		return gitCloneUrl(repoStr), nil
	}

	return "", errors.New("no available github mirror for the repo found")
}
