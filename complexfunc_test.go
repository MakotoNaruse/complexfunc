package complexfunc_test

import (
	"github.com/MakotoNaruse/complexfunc"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, complexfunc.Analyzer, "a")
}
