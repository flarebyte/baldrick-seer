package app

const reportGenerateStubOutput = "report generate: ok\n"

type ReportGenerateRequest struct{}

type ReportGenerateResponse struct {
	Stdout string
}

func RunReportGenerate(_ ReportGenerateRequest) (ReportGenerateResponse, error) {
	return ReportGenerateResponse{
		Stdout: reportGenerateStubOutput,
	}, nil
}
