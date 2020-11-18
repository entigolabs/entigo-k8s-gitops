package util

import (
	"testing"
)

func TestComposeFeatureBranchName(t *testing.T) {
	testBranchNameStriping(t)
	testBranchNameHashing(t)
}

func testBranchNameStriping(t *testing.T) {
	thirtyTwoChars := "11111111110111111111101111111111"
	over32Chars := thirtyTwoChars + "x"

	if composeFeatureBranchName(over32Chars) != thirtyTwoChars {
		t.Error("Striping didn't work")
	}
}

func testBranchNameHashing(t *testing.T) {
	invalidBranchName := "#invalidBranchName"
	expectedHash := "97312217f0dd75c8c6075080c23c80f3"

	if composeFeatureBranchName(invalidBranchName) != expectedHash {
		t.Errorf("Invalid branch name (%s) could not be converted into expected hash: %s", invalidBranchName, expectedHash)
	}
}
