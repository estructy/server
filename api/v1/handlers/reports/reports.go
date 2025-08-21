// Package reportshandler provides a handler for report-related operations.
package reportshandler

import (
	"net/http"

	reportbycategory "github.com/nahtann/controlriver.com/internal/domain/reports/use_cases/by_category"
	bycategoryrequest "github.com/nahtann/controlriver.com/internal/domain/reports/use_cases/by_category/request"
	contexthelper "github.com/nahtann/controlriver.com/internal/helpers/context"
	jsonhelper "github.com/nahtann/controlriver.com/internal/helpers/json"
	requesthelper "github.com/nahtann/controlriver.com/internal/helpers/request"
	"github.com/nahtann/controlriver.com/internal/infra/database/repository"
)

type ReportsHandler struct {
	repository *repository.Queries
}

func NewReportsHandler(repository *repository.Queries) *ReportsHandler {
	return &ReportsHandler{
		repository: repository,
	}
}

func (uc *ReportsHandler) GetReportByCategory(w http.ResponseWriter, r *http.Request) {
	accountID, ok := contexthelper.AccountIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Account ID not found in context", http.StatusInternalServerError)
		return
	}

	categoryType := r.URL.Query().Get("type")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	withSubCategories := r.URL.Query().Get("with-sub-categories") == "true"
	withTransactions := r.URL.Query().Get("with-transactions") == "true"

	request := &bycategoryrequest.Request{
		Type:              categoryType,
		From:              from,
		To:                to,
		WithSubCategories: withSubCategories,
		WithTransactions:  withTransactions,
	}

	errorMessages := requesthelper.ValidateRequest(request)
	if errorMessages != "" {
		jsonhelper.HTTPError(w, http.StatusBadRequest, errorMessages)
		return
	}

	reportByCategoryUseCase := reportbycategory.NewReportBycategoryUseCase(uc.repository)
	response, err := reportByCategoryUseCase.Execute(accountID, request)
	if err != nil {
		errMappings := map[error]jsonhelper.ErrorMappings{
			reportbycategory.ErrInvalidDateFormat: {Code: http.StatusBadRequest, Message: "Invalid date format, expected YYYY-MM-DD"},
			reportbycategory.ErrFailedToGetReport: {Code: http.StatusInternalServerError, Message: "Failed to get report by category"},
		}
		jsonhelper.HandleError(w, err, errMappings)
		return
	}

	jsonhelper.HTTPResponse(w, http.StatusOK, response)
}
