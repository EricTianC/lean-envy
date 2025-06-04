package provider

import (
	"reflect"
	"testing"
)

/*
require "leanprover-community" / "mathlib"
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
require "leanprover-community" / "plausible" @ git "main"
*/
func TestParseLeanLakeFileRequirements(t *testing.T) {
	mathlibReqs, err := parseLeanLakeFileRequirements(
		`require "leanprover-community" / "mathlib"
require "leanprover-community" / "mathlib" @ "git#20c73142afa995ac9c8fb80a9bb585a55ca38308"
require "leanprover-community" / "Qq" @ git "master"
require "leanprover-community" / "batteries" @ git "main"
require "leanprover-community" / "proofwidgets" @ git "v0.0.60"
`)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(mathlibReqs,
		[]Requirement{
			{Scope: "leanprover-community", PkgName: "mathlib"},
			{Scope: "leanprover-community", PkgName: "mathlib", Version: "git#20c73142afa995ac9c8fb80a9bb585a55ca38308"},
			{Scope: "leanprover-community", PkgName: "Qq", Version: "master"},
			{Scope: "leanprover-community", PkgName: "batteries", Version: "main"},
			{Scope: "leanprover-community", PkgName: "proofwidgets", Version: "v0.0.60"},
		}) {
		t.Error("parse 失败: ", "requirements", mathlibReqs)
	}
}

func TestParseTomlLakeFileRequirements(t *testing.T) {
	content := `
[[require]]
name = "Cli"
scope = "leanprover"
rev = "main"

[[require]]
name = "batteries"
scope = "leanprover-community"
rev = "v4.21.0-rc1"`
	reqs, _ := parseTomlLakeFileRequirements(content) // 不可能发生异常
	if !reflect.DeepEqual(reqs,
		[]Requirement{
			{Scope: "leanprover", PkgName: "Cli", Version: "main"},
			{Scope: "leanprover-community", PkgName: "batteries", Version: "v4.21.0-rc1"},
		}) {
		t.Error("test: ", "reqs", reqs)
	}
}
