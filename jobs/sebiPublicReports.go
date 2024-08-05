package jobs

import (
	"context"
	"fmt"
	"sebi-scrapper/constants"
	"sebi-scrapper/entities/repositories"
	"sebi-scrapper/externals"
	modelsv1 "sebi-scrapper/models/v1"
)

func SebiPublicReports(ctx context.Context) error {

	var response []modelsv1.Report
	for _, department := range constants.Departments {

		resp, err := externals.GetSebiPublicReports(department)
		if err != nil {
			return err
		}

		response = append(response, resp...)
	}
	for _, report := range response {
		err := repositories.SaveReport(ctx, report)
		if err != nil {
			fmt.Print("Error in saving : ", err)
			return err
		}
	}
	// Save the reports which do not fall under any department category
	// Crawl all the reports irrespsetive of the departments
	fmt.Println("Crawling all the reports")
	resp, err := externals.GetSebiPublicReports(constants.AllReports)
	if err != nil {
		return err
	}
	// Save all the reports without any department in uncategorised departments
	for _, report := range resp {
		err := repositories.SaveReport(ctx, report)
		if err != nil {
			fmt.Print("Error in saving : ", err)
			return err
		}
	}

	return nil
}
