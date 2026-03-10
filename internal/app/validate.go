package app

const validateStubOutput = "validate: ok\n"

type ValidateRequest struct {
	ConfigPath string
}

type ValidateResponse struct {
	Stdout string
}

func RunValidate(req ValidateRequest) (ValidateResponse, error) {
	_, err := LoadConfig(ConfigRequest{Path: req.ConfigPath})
	if err != nil {
		return ValidateResponse{}, err
	}

	return ValidateResponse{
		Stdout: validateStubOutput,
	}, nil
}
