// Package reportbycategory provides the use case for generating a report by category.
package reportbycategory

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	bycategoryrequest "github.com/nahtann/controlriver.com/internal/domain/reports/use_cases/by_category/request"
	bycategoryresponse "github.com/nahtann/controlriver.com/internal/domain/reports/use_cases/by_category/response"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

var (
	ErrInvalidDateFormat = fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	ErrFailedToGetReport = fmt.Errorf("failed to get report by category")
)

type ReportBycategoryUseCase struct {
	repository *repository.Queries
}

func NewReportBycategoryUseCase(repository *repository.Queries) *ReportBycategoryUseCase {
	return &ReportBycategoryUseCase{
		repository: repository,
	}
}

func (uc *ReportBycategoryUseCase) Execute(accountID uuid.UUID, request *bycategoryrequest.Request) (*bycategoryresponse.Response, error) {
	fromDate, err := time.Parse("2006-01-02", request.From)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDateFormat, request.From)
	}
	toDate, err := time.Parse("2006-01-02", request.To)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDateFormat, request.To)
	}

	ctx := context.Background()

	var response *bycategoryresponse.Response
	if request.WithSubCategories {
		reportyByCategory, err := uc.repository.GetReportByCategoriesWithParentCategory(ctx, repository.GetReportByCategoriesWithParentCategoryParams{
			AccountID: &accountID,
			Type:      request.Type,
			From:      fromDate,
			To:        toDate,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToGetReport, err)
		}

		response = bycategoryresponse.NewResponseWithSubCategories(request.From, request.To, request.Type, reportyByCategory)
	} else {
		reportyByCategory, err := uc.repository.GetReportByCategories(ctx, repository.GetReportByCategoriesParams{
			AccountID: &accountID,
			Type:      request.Type,
			From:      fromDate,
			To:        toDate,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToGetReport, err)
		}

		response = bycategoryresponse.NewResponse(request.From, request.To, request.Type, reportyByCategory)
	}

	return response, nil
}
