package pipeline

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func largeValidConfigPath() string {
	return filepath.Join("..", "..", "testdata", "config", "large_valid.cue")
}

func largeInvalidConfigPath() string {
	return filepath.Join("..", "..", "testdata", "config", "large_invalid.cue")
}

func TestLargeFixturePipelineStressDeterminism(t *testing.T) {
	t.Parallel()

	runner := NewDefaultRunner()
	command := domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  largeValidConfigPath(),
	}

	first, err := runner.RunReportGenerate(context.Background(), command)
	if err != nil {
		t.Fatalf("first RunReportGenerate() error = %v", err)
	}

	for i := 0; i < 12; i++ {
		got, err := runner.RunReportGenerate(context.Background(), command)
		if err != nil {
			t.Fatalf("RunReportGenerate() iteration %d error = %v", i, err)
		}

		if !reflect.DeepEqual(first, got) {
			t.Fatalf("iteration %d result drifted", i)
		}
	}
}

func TestLargeFixtureValidateStressDeterminism(t *testing.T) {
	t.Parallel()

	runner := NewDefaultRunner()
	command := domain.CommandRequest{
		CommandName: domain.CommandNameValidate,
		ConfigPath:  largeValidConfigPath(),
	}

	first, err := runner.RunValidate(context.Background(), command)
	if err != nil {
		t.Fatalf("first RunValidate() error = %v", err)
	}

	for i := 0; i < 12; i++ {
		got, err := runner.RunValidate(context.Background(), command)
		if err != nil {
			t.Fatalf("RunValidate() iteration %d error = %v", i, err)
		}

		if !reflect.DeepEqual(first, got) {
			t.Fatalf("iteration %d result drifted", i)
		}
	}
}

func TestLargeInvalidFixtureValidationStressDeterminism(t *testing.T) {
	t.Parallel()

	runner := NewDefaultRunner()
	command := domain.CommandRequest{
		CommandName: domain.CommandNameValidate,
		ConfigPath:  largeInvalidConfigPath(),
	}

	_, firstErr := runner.RunValidate(context.Background(), command)
	first := domain.PresentError(firstErr)

	for i := 0; i < 12; i++ {
		_, err := runner.RunValidate(context.Background(), command)
		got := domain.PresentError(err)

		if !reflect.DeepEqual(first, got) {
			t.Fatalf("iteration %d error drifted: first = %#v, got = %#v", i, first, got)
		}
	}
}

func TestLargeInvalidFixtureFailsDeterministically(t *testing.T) {
	t.Parallel()

	_, err := NewDefaultRunner().RunReportGenerate(context.Background(), domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  largeInvalidConfigPath(),
	})

	failure := assertCommandFailure(
		t,
		err,
		ErrValidationFailed,
		"unknown criterion name in evaluation values: unknown",
	)

	if failure.Category != domain.FailureCategoryValidation {
		t.Fatalf("Category = %q, want %q", failure.Category, domain.FailureCategoryValidation)
	}
}
