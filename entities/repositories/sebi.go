package repositories

import (
	"context"
	"fmt"
	"sebi-scrapper/constants"
	"sebi-scrapper/entities"
	modelsv1 "sebi-scrapper/models/v1"
	"sebi-scrapper/utils/database"
	"time"
)

func SaveReport(ctx context.Context, data modelsv1.Report) error {
	date, err := time.Parse("Jan 2, 2006", data.Date)
	if err != nil {
		return err
	}
	fmt.Println("date : ", date)
	EntityData := &entities.Reports{Title: data.Title, Department: data.Department, Content: data.Content, Date: date}
	err = database.Get().QueryRaw(ctx, EntityData, entities.CheckReportExists)
	if err != nil {
		return err
	}
	if !EntityData.Exists {
		fmt.Println("Didnt Exist !")
		err = database.Get().QueryRaw(ctx, EntityData, entities.SaveReport)
		if err == nil {
			fmt.Println("ID :", EntityData.ID)
		}
	} else {
		fmt.Println("Already Exist !")
	}

	return err
}

func GetSebiReports(ctx context.Context, department, order string, offset, size int) ([]*entities.Reports, error) {

	EntityData := &entities.Reports{Offset: offset, Size: size}
	code := entities.GetAllPublicReportsDesc
	if order == constants.Ascending {
		code = entities.GetAllPublicReportsAsc
	}
	if department != constants.AllDepartment {
		EntityData.Department = department
		code = entities.GetPublicReportsWithDepartmentDesc
		if order == constants.Ascending {
			code = entities.GetAllPublicReportsAsc
		}
	}
	response, err := database.Get().QueryMultiRaw(ctx, EntityData, code)

	result := make([]*entities.Reports, len(response))
	for i, b := range response {
		result[i], _ = b.(*entities.Reports)
	}

	return result, err
}

func GetSebiReportsCount(ctx context.Context, department string) (int, error) {

	EntityData := &entities.Reports{}
	code := entities.GetAllPublicReportsCount
	if department != constants.AllDepartment {
		EntityData.Department = department
		code = entities.GetPublicReportsWithDepartmentCount
	}
	err := database.Get().QueryRaw(ctx, EntityData, code)
	return EntityData.Count, err
}
