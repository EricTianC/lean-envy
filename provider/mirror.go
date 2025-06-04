package provider

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

var GitRepoMirrors = map[string]string{
	"leanprover-community/mathlib":          "https://mirror.sjtu.edu.cn/git/lean4-packages/mathlib4",                // ✔️
	"leanprover-community/batteries":        "https://mirror.sjtu.edu.cn/git/lean4-packages/batteries",               // ✔️
	"leanprover-community/Qq":               "https://mirror.sjtu.edu.cn/git/lean4-packages/quote4",                  // ✔️
	"leanprover-community/aesop":            "https://mirror.sjtu.edu.cn/git/lean4-packages/aesop",                   // ✔️
	"leanprover-community/proofwidgets":     "https://mirror.sjtu.edu.cn/git/lean4-packages/ProofWidgets4",           // ✔️
	"leanprover-community/importGraph":      "https://mirror.sjtu.edu.cn/git/lean4-packages/import-graph",            // ✔️
	"leanprover-community/LeanSearchClient": "https://gitclone.com/github.com/leanprover-community/LeanSearchClient", // ✔️
	"leanprover-community/plausible":        "https://gitclone.com/github.com/leanprover-community/plausible",        // ✔️
	"leanprover/Cli":                        "https://mirror.sjtu.edu.cn/git/lean4-packages/lean4-cli",               // ❓
}
