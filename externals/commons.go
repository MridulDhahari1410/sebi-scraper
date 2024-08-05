package externals

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/csv"
	"io"
	goHttp "net/http"

	"sebi-scrapper/constants"
	"sebi-scrapper/utils"
	"sebi-scrapper/utils/http"
	"sebi-scrapper/utils/metrics"

	httpclient "github.com/angel-one/go-http-client"
)

type parser func(data io.ReadCloser, value any) error

var (
	marshalJSON     = utils.MarshalJSON
	getDataAsString = utils.GetDataAsString
	getJSONData     = utils.GetJSONData
	getDataAsBytes  = utils.GetDataAsBytes
	httpGet         = http.Get
)

var (
	metricsIncrementExternalHTTPResponseCounter = metrics.IncrementExternalHTTPResponseCounter
	metricsGetExternalHTTPRequestTimer          = metrics.GetExternalHTTPRequestTimer
)

func callAndGetResponseByParser(ctx context.Context, name string, request *httpclient.Request, success int,
	data any, p parser) error {
	metricsIncrementExternalHTTPResponseCounter(name)
	// check if call is permitted
	if !isRequestPermitted(ctx, name) {
		return constants.ErrForbidden.Value()
	}

	timer := metricsGetExternalHTTPRequestTimer(name)
	// do the request
	response, err := httpGet().Request(request.SetContext(ctx))
	timer.ObserveDuration()
	if err != nil {
		return constants.ErrServiceUnavailable.WithDetails(err.Error())
	}

	// check for success status codes
	if response.StatusCode == success {
		err = p(response.Body, data)
		if err != nil {
			return constants.ErrServiceUnexpectedResponse.WithDetails(err.Error())
		}
		return nil
	}

	// handle the error body parse
	return handleNonSuccessResponse(response)
}

func callAndGetResponse(ctx context.Context, name string, request *httpclient.Request, success int, data any) error {
	return callAndGetResponseByParser(ctx, name, request, success, data, getJSONData)
}

func callAndGetGZippedResponse(ctx context.Context, name string, request *httpclient.Request,
	success int, data any) error {
	return callAndGetResponseByParser(ctx, name, request, success, data, getGZippedParser(getJSONData))
}

func callAndGetRawResponse(ctx context.Context, name string, request *httpclient.Request, success int) ([]byte, error) {
	// check if call is permitted
	if !isRequestPermitted(ctx, name) {
		return nil, constants.ErrForbidden.Value()
	}

	// do the request
	response, err := httpGet().Request(request)
	if err != nil {
		return nil, constants.ErrServiceUnavailable.WithDetails(err.Error())
	}

	// check for success status codes
	if response.StatusCode == success {
		b, err := getDataAsBytes(response.Body)
		if err != nil {
			return nil, constants.ErrServiceUnexpectedResponse.WithDetails(err.Error())
		}
		return b, nil
	}

	// handle the error body parse
	return nil, handleNonSuccessResponse(response)
}

func callAndGetResponseWithRequestBody(ctx context.Context, name string, body any, request *httpclient.Request,
	success int, data any) error {
	// marshal the body
	b, err := marshalJSON(body)
	if err != nil {
		return constants.ErrServiceUnexpectedRequestBody.WithDetails(err.Error())
	}

	// handle
	return callAndGetResponse(ctx, name, request.SetBody(bytes.NewReader(b)), success, data)
}

func parseCSV[T any](r io.Reader, shouldUseRow func(i int) bool, rowScanner func(data []string) T) ([]T, error) {
	reader := csv.NewReader(r)
	var result []T
	record, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(record); i++ {
		if shouldUseRow(i) {
			result = append(result, rowScanner(record[i]))
		}
	}
	return result, nil
}

func isRequestPermitted(_ context.Context, _ string) bool {
	return true
}

func getGZippedParser(p parser) parser {
	return func(data io.ReadCloser, value any) error {
		defer utils.CloseData(data)
		r, err := gzip.NewReader(data)
		if err != nil {
			return err
		}
		return p(r, value)
	}
}

func handleNonSuccessResponse(response *goHttp.Response) error {
	if response.StatusCode == goHttp.StatusNoContent {
		// here we need to handle it as no content case
		return constants.ErrServiceNoContentFailureInformation.Value()
	}
	s, err := getDataAsString(response.Body)
	if err != nil {
		s = err.Error()
	}
	if response.StatusCode >= goHttp.StatusBadRequest && response.StatusCode < goHttp.StatusInternalServerError {
		return constants.ErrServiceBadRequestFailureInformation.WithDetails(s)
	}
	return constants.ErrServiceUnexpectedError.WithDetails(s)
}
