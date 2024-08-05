package utils_test

import (
	"context"
	"testing"
	"time"

	"sebi-scrapper/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetRunAtWithOffset(t *testing.T) {
	assert.Equal(t, "05:30", utils.GetRunAtWithOffset(context.Background(), "05:30", 0))
	for i := 0; i < 100; i++ {
		// Running multiple iterations
		result := utils.GetRunAtWithOffset(context.Background(), "05:30", 15)
		parsedTime, err := time.Parse("15:04", result) //converting to 24hr format
		if err != nil {
			t.Fatalf("Failed to parse time: %v", err)
		}
		expectedMinTime, err := time.Parse("15:04", "05:15") //converting to 24hr format
		if err != nil {
			t.Fatalf("Failed to parse minimum expected time: %v", err)
		}
		expectedMaxTime, err := time.Parse("15:04", "05:45") //converting to 24hr format
		if err != nil {
			t.Fatalf("Failed to parse maximum expected time: %v", err)
		}
		assert.True(t, parsedTime.After(expectedMinTime) || parsedTime.Equal(expectedMinTime), "Expected time after or equal to 05:15")
		assert.True(t, parsedTime.Before(expectedMaxTime) || parsedTime.Equal(expectedMaxTime), "Expected time before or equal to 05:45")
	}
}
