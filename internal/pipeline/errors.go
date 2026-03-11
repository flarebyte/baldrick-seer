package pipeline

import "errors"

var (
	ErrConfigPathRequired     = errors.New("config flag is required")
	ErrConfigPathDoesNotExist = errors.New("config path does not exist")
	ErrConfigPathIsDirectory  = errors.New("config path is a directory")
	ErrConfigLoadInvalid      = errors.New("config could not be loaded")
	ErrConfigNotConcrete      = errors.New("config must evaluate to a concrete value")
)
