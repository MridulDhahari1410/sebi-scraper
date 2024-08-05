package apiv1

import (
	businessv1 "sebi-scrapper/business/v1"
	"sebi-scrapper/constants"
	"sebi-scrapper/jobs"
	"sebi-scrapper/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HandleGetAllStrategies godoc
//
//	@Summary		Get all strategies
//	@Description	Get all strategies
//	@ID				getAllStrategies
//	@Tags			strategy
//	@Produce		json
//	@Param			AccessToken	header		string	true	"${JWT}"
//	@Param			X-source	header		string	true	"source"
//	@Success		200			{object}	models.Response{data=[]modelsv1.Strategy}
//	@Failure		400			{object}	models.Response
//	@Failure		500			{object}	models.Response
//	@Router			/v1/strategies [get].
func HandleCrawlSebiReports(ctx *gin.Context) {
	err := jobs.SebiPublicReports(ctx)
	if err != nil {
		ctx.JSON(models.GetErrorResponse(err))
	}
	ctx.JSON(models.GetNotFoundResponse())
}

func HandleGetPublicReports(ctx *gin.Context) {
	department := ctx.DefaultQuery(constants.DepartmentQuery, constants.AllDepartment)
	order := ctx.DefaultQuery(constants.OrderQuery, constants.Descending)
	if order != constants.Descending && order != constants.Ascending {
		ctx.JSON(models.GetErrorResponse(&constants.ErrInvalidParamOrder))
		return
	}
	page := ctx.DefaultQuery(constants.Page, constants.DefaultPage)
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(models.GetErrorResponse(err))
		return
	}
	response, err := businessv1.GetPublicReports(ctx, department, order, pageInt)
	if err != nil {
		ctx.JSON(models.GetErrorResponse(err))
		return
	}
	ctx.JSON(models.GetOKSuccessResponse(response))

}

func HandleGetPublicReportsCount(ctx *gin.Context) {
	department := ctx.DefaultQuery(constants.DepartmentQuery, constants.AllDepartment)
	response, err := businessv1.GetPublicReportsCount(ctx, department)
	if err != nil {
		ctx.JSON(models.GetErrorResponse(err))
		return
	}
	ctx.JSON(models.GetOKSuccessResponse(response))

}

func HandleGetReportsDepartmentsList(ctx *gin.Context) {
	response := constants.Departments
	response = append(response, constants.Uncategorised)
	response = append(response, constants.AllDepartment)
	ctx.JSON(models.GetOKSuccessResponse(response))
}
