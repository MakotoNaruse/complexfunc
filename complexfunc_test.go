package complexfunc_test

import (
	"testing"

	"practice/complexfunc"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, complexfunc.Analyzer, "a")
}

