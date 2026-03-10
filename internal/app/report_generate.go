package app

const reportGenerateStubOutput = "report generate: ok\n"

type ReportGenerateRequest struct {
	ConfigPath string
}

type ReportGenerateResponse struct {
	Stdout string
}

func RunReportGenerate(req ReportGenerateRequest) (ReportGenerateResponse, error) {
	_, err := LoadConfig(ConfigRequest{Path: req.ConfigPath})
	if err != nil {
		return ReportGenerateResponse{}, err
	}

	return ReportGenerateResponse{
		Stdout: reportGenerateStubOutput,
	}, nil
}
