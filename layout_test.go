package layout_test

import (
	"testing"

	layout "github.com/gomatic/yze-go-layout"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLayoutCorrespondence(t *testing.T) {
	analysistest.Run(
		t, analysistest.TestData(), layout.Analyzer,
		"m/internal/app/commands/greet",
		"m/internal/domain/greet",
		"m/internal/app/commands/orphan",
		"m/internal/domain/lonely",
		"m/internal/util",
	)
}

func TestRegistrationIsWellFormed(t *testing.T) {
	assert.NoError(t, layout.Registration.Validate())
	assert.Equal(t, "yze/go/layout", layout.Registration.RuleID())
	assert.Same(t, layout.Analyzer, layout.Registration.Analyzer)
}
