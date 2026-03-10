package app

const validateStubOutput = "validate: ok\n"

type ValidateRequest struct{}

type ValidateResponse struct {
	Stdout string
}

func RunValidate(_ ValidateRequest) (ValidateResponse, error) {
	return ValidateResponse{
		Stdout: validateStubOutput,
	}, nil
}
