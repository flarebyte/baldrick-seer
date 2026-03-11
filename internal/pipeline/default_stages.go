package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueformat "cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type DefaultConfigLoader struct{}

func (DefaultConfigLoader) LoadConfig(ctx context.Context, input LoadConfigInput) (LoadConfigOutput, error) {
	if err := checkContext(ctx, input.ConfigPath); err != nil {
		return LoadConfigOutput{}, err
	}

	if input.ConfigPath == "" {
		return LoadConfigOutput{}, NewInputFailure("config.required", "", "config flag is required", ErrConfigPathRequired)
	}

	configPath := filepath.Clean(input.ConfigPath)

	info, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return LoadConfigOutput{}, NewInputFailure("config.not_found", input.ConfigPath, "config path does not exist", ErrConfigPathDoesNotExist)
		}

		return LoadConfigOutput{}, WrapStageFailure(domain.FailureCategoryInternal, "config.stat_failed", input.ConfigPath, "command failed", err)
	}

	loadArgs, loadDir, err := cueLoadTarget(configPath, info)
	if err != nil {
		return LoadConfigOutput{}, err
	}

	instances := load.Instances(loadArgs, &load.Config{
		Dir: loadDir,
	})
	if len(instances) != 1 || instances[0] == nil {
		return LoadConfigOutput{}, NewInputFailure("config.load_invalid", input.ConfigPath, "config could not be loaded", ErrConfigLoadInvalid)
	}

	instance := instances[0]
	if instance.Err != nil {
		return LoadConfigOutput{}, NewInputFailure("config.load_invalid", input.ConfigPath, "config could not be loaded", ErrConfigLoadInvalid)
	}

	value := cuecontext.New().BuildInstance(instance)
	if err := value.Err(); err != nil {
		return LoadConfigOutput{}, NewInputFailure("config.load_invalid", input.ConfigPath, "config could not be loaded", ErrConfigLoadInvalid)
	}

	if err := value.Validate(cue.Concrete(true)); err != nil {
		return LoadConfigOutput{}, NewInputFailure("config.not_concrete", input.ConfigPath, "config must evaluate to a concrete value", ErrConfigNotConcrete)
	}

	syntax := value.Syntax(
		cue.Concrete(true),
		cue.Definitions(false),
		cue.Hidden(false),
		cue.Optional(false),
	)
	formatted, err := cueformat.Node(syntax)
	if err != nil {
		return LoadConfigOutput{}, WrapStageFailure(domain.FailureCategoryInternal, "config.format_failed", input.ConfigPath, "command failed", err)
	}

	fields, err := cueTopLevelFields(value)
	if err != nil {
		return LoadConfigOutput{}, WrapStageFailure(domain.FailureCategoryInternal, "config.fields_failed", input.ConfigPath, "command failed", err)
	}

	var executionConfig *ExecutionConfig
	var configFields []string
	configValue := value.LookupPath(cue.ParsePath("config"))
	if configValue.Exists() {
		if err := checkContext(ctx, input.ConfigPath); err != nil {
			return LoadConfigOutput{}, err
		}
		configFields, err = cueTopLevelFields(configValue)
		if err != nil {
			return LoadConfigOutput{}, NewInputFailure("config.decode_invalid", input.ConfigPath, "config could not be loaded", ErrConfigLoadInvalid)
		}

		decoded := new(ExecutionConfig)
		if err := configValue.Decode(decoded); err != nil {
			return LoadConfigOutput{}, NewInputFailure("config.decode_invalid", input.ConfigPath, "config could not be loaded", ErrConfigLoadInvalid)
		}
		executionConfig = decoded
	}

	return LoadConfigOutput{
		Config: LoadedConfig{
			Path:           configPath,
			Evaluated:      string(formatted),
			TopLevelFields: domain.CanonicalNames(fields),
			ConfigFields:   domain.CanonicalNames(configFields),
			Config:         executionConfig,
		},
	}, nil
}

type DefaultCriteriaWeighter struct{}

type DefaultScenarioRanker struct{}

type DefaultScenarioAggregator struct{}

type DefaultReportRenderer struct{}

func cueLoadTarget(configPath string, info os.FileInfo) ([]string, string, error) {
	if !info.IsDir() {
		if filepath.Ext(configPath) != ".cue" {
			return nil, "", NewInputFailure("config.invalid_file_type", configPath, "config file must have .cue extension", ErrConfigFileExtension)
		}

		return []string{filepath.Base(configPath)}, filepath.Dir(configPath), nil
	}

	entries, err := os.ReadDir(configPath)
	if err != nil {
		return nil, "", WrapStageFailure(domain.FailureCategoryInternal, "config.read_dir_failed", configPath, "command failed", err)
	}

	var cueFiles []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".cue" {
			continue
		}
		cueFiles = append(cueFiles, entry.Name())
	}

	sort.Strings(cueFiles)
	if len(cueFiles) == 0 {
		return nil, "", NewInputFailure("config.directory_empty", configPath, "config directory does not contain any .cue files", ErrConfigDirectoryEmpty)
	}

	return cueFiles, configPath, nil
}

func cueTopLevelFields(value cue.Value) ([]string, error) {
	iterator, err := value.Fields(
		cue.Definitions(false),
		cue.Hidden(false),
		cue.Optional(false),
	)
	if err != nil {
		return nil, err
	}

	var fields []string
	for iterator.Next() {
		fields = append(fields, iterator.Selector().Unquoted())
	}

	sort.Strings(fields)
	return fields, nil
}
