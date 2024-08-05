package externals

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sebi-scrapper/constants"
	modelsv1 "sebi-scrapper/models/v1"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Define the struct to store data

func GetSebiPublicReports(department string) ([]modelsv1.Report, error) {
	c := colly.NewCollector()
	var reports []modelsv1.Report
	var totalPages int
	// Callback for when a HTML element is visited
	cutoffDate, _ := time.Parse("Jan 2, 2006", "May 10, 1999")
	crawledAllLatestReports := false
	fmt.Println("Departmet id : ", constants.DepartmentToValue[department], "  department :", department)
	// Callback for when a HTML element is visited

	c.OnHTML(".pagination_inner p", func(e *colly.HTMLElement) {
		text := e.Text
		// Use regular expression to extract the total number of records
		re := regexp.MustCompile(`(\d+) records`)
		match := re.FindStringSubmatch(text)
		if len(match) > 1 {
			totalRecordsStr := match[1]
			totalRecords, err := strconv.Atoi(totalRecordsStr)
			if err != nil {
				fmt.Printf("Error parsing total records: %v\n", err)
				return
			}
			totalPages = (totalRecords + 24) / 25 // Each page has up to 25 records
		} else {
			fmt.Println("Error: Total records not found in text")
		}
	})

	c.OnHTML("tr[role=row].odd", func(e *colly.HTMLElement) {
		// fmt.Println("HTML electemt :", e.Text)
		date := e.ChildText("td:nth-child(1)")

		// Extract link from the <a> tag within the <td> element
		link := e.ChildAttr("td:nth-child(2) a", "href")

		// Extract title from the title attribute of the <a> tag
		titleAttr := e.ChildAttr("td:nth-child(2) a", "title")
		title := strings.Split(titleAttr, "<a")[0]

		// Print extracted information for debugging
		// fmt.Printf("Date: %s\n", date)
		// fmt.Printf("Link: %s\n", link)
		// fmt.Printf("Title: %s\n", title)

		// Convert date string to time.Time
		reportDate, err := time.Parse("Jan 2, 2006", date)
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err)
			return
		}

		fmt.Println("Date : ", reportDate, "   Department :", department)
		// Check if the date is before the cutoff date
		if !reportDate.After(cutoffDate) {
			fmt.Println("No more relevant reports found, stopping the crawl.")
			crawledAllLatestReports = true
			return // Exit the crawl
		}

		// Visit the link to get the PDF link
		pdfCollector := colly.NewCollector()

		pdfCollector.OnHTML("iframe", func(e *colly.HTMLElement) {
			// Extract the src attribute from the iframe
			iframeSrc := e.Attr("src")
			if strings.HasPrefix(iframeSrc, "../../../web/?file=") {
				// Construct the full PDF URL
				pdfURL := strings.TrimPrefix(iframeSrc, "../../../web/?file=")
				baseURL := pdfURL[0:11]
				if baseURL == "/sebi_data/" {
					pdfURL = "https://www.sebi.gov.in" + pdfURL
				}
				// Download the PDF
				fmt.Println("PDF URL :", pdfURL)
				pdfData, err := downloadPDF(pdfURL)
				if err != nil {
					fmt.Printf("Error downloading PDF: %v\n", err)
					return
				}

				// Store report
				departmentCategory := department
				if departmentCategory == constants.AllReports {
					departmentCategory = "uncategorised"
				}
				report := modelsv1.Report{
					Date:       date,
					Title:      title,
					Content:    pdfData,
					Department: departmentCategory,
				}
				reports = append(reports, report)
			}
		})

		// Visit the link for PDF collection
		pdfCollector.Visit(link)
	})

	// Start scraping

	for pageNumber := 1; pageNumber <= totalPages || totalPages == 0; pageNumber++ {
		if crawledAllLatestReports {
			break
		}
		nextValue := fmt.Sprintf("%d", pageNumber)
		next := "s"
		if pageNumber > 1 {
			next = "n"
			nextValue = fmt.Sprintf("%d", pageNumber-1)
		}
		headers := map[string]string{
			"Accept":             "*/*",
			"Accept-Encoding":    "gzip, deflate, br, zstd",
			"Accept-Language":    "en-GB,en-US;q=0.9,en;q=0.8",
			"Cache-Control":      "no-cache",
			"Connection":         "keep-alive",
			"Content-Type":       "application/x-www-form-urlencoded",
			"Host":               "www.sebi.gov.in",
			"Origin":             "https://www.sebi.gov.in",
			"Pragma":             "no-cache",
			"Referer":            "https://www.sebi.gov.in/sebiweb/home/HomeAction.do?doListing=yes&sid=4&ssid=38&smid=35",
			"Sec-Fetch-Mode":     "cors",
			"Sec-Fetch-Site":     "same-origin",
			"User-Agent":         "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36",
			"sec-ch-ua":          `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`,
			"sec-ch-ua-platform": `"Android"`,
			"Cookie":             "JSESSIONID=E79C2C769468AA65B687F393F38F31AF",
		}
		departmentId := "-1"
		if department != constants.AllReports {
			departmentId = constants.DepartmentToValue[department]
		}
		// Define form data
		formData := map[string]string{
			"nextValue":  nextValue,
			"next":       next,
			"search":     "",
			"fromDate":   "",
			"toDate":     "",
			"fromYear":   "",
			"toYear":     "",
			"deptId":     departmentId,
			"sid":        "4",
			"ssid":       "38",
			"smid":       "35",
			"ssidhidden": "38",
			"intmid":     "-1",
			"sText":      "Reports ",
			"Statistics": "",
			"ssText":     "Reports",
			"smText":     "Reports for Public Comments",
			"doDirect":   "-1",
		}

		c.OnRequest(func(r *colly.Request) {
			for key, value := range headers {
				r.Headers.Set(key, value)
			}
		})
		fmt.Println("Pagenumber : ", pageNumber)
		oldsize := len(reports)
		err := c.Post("https://www.sebi.gov.in/sebiweb/ajax/home/getnewslistinfo.jsp", formData)
		if err != nil {
			fmt.Printf("Error starting scraping: %v\n", err)
			return reports, nil
		}

		c.OnResponse(func(r *colly.Response) {
			// fmt.Printf("Response Status: %d\n", r.StatusCode)
			// fmt.Printf("Response Body: %s\n", string(r.Body))
		})
		// Stop condition if no more relevant reports are found
		if len(reports) == oldsize {
			break
		}

	}
	return reports, nil
}

func downloadPDF(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
