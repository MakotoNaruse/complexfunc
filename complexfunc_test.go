package complexfunc_test

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"practice/complexfunc"
	"testing"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, complexfunc.Analyzer, "a")
}
