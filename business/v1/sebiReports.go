package businessv1

import (
	"context"
	"sebi-scrapper/entities/repositories"
	modelsv1 "sebi-scrapper/models/v1"
)

func GetPublicReports(ctx context.Context, department, order string, page int) ([]modelsv1.GetPublicReportsResponse, error) {
	size := 20
	offset := (page - 1) * size
	response, err := repositories.GetSebiReports(ctx, department, order, offset, size)
	if err != nil {
		return nil, err
	}
	result := make([]modelsv1.GetPublicReportsResponse, len(response))
	for i, report := range response {
		result[i] = modelsv1.GetPublicReportsResponse{ID: report.ID, Date: report.Date, Title: report.Title, Department: report.Department, Status: report.Status}
	}

	return result, err
}

func GetPublicReportsCount(ctx context.Context, department string) (*modelsv1.RecordsSize, error) {
	count, err := repositories.GetSebiReportsCount(ctx, department)
	if err != nil {
		return nil, err
	}
	return &modelsv1.RecordsSize{Size: count}, nil
}
