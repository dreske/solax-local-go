package solax_local_go

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJsonMapping(t *testing.T) {
	simpleResponse := "{\"sn\":\"SVRDJTTEUH\",\"ver\":\"3.003.02\",\"type\":4,\"Data\":[0,0,0,0,0,0,0,0,0,0,0,7337,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1708,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"Information\":[4.200,4,\"XB3342I2094013\",8,0.00,0.00,1.38,0.00,0.00,1]}"

	var parsedResponse RawInverterResponse
	require.NoError(t, json.Unmarshal([]byte(simpleResponse), &parsedResponse))

	assert.Equal(t, "SVRDJTTEUH", parsedResponse.SN)
	assert.Equal(t, "3.003.02", parsedResponse.Version)
	assert.Equal(t, 4, parsedResponse.Type)
	assert.Equal(t, 4.2, parsedResponse.Information[0])
	assert.Equal(t, float64(4), parsedResponse.Information[1])
	assert.Equal(t, "XB3342I2094013", parsedResponse.Information[2])
}

func TestParseSimpleX1BoostResponse(t *testing.T) {
	simpleResponse := "{\"sn\":\"SVRDJTTEUH\",\"ver\":\"3.003.02\",\"type\":4,\"Data\":[0,0,0,0,0,0,0,0,0,0,0,7337,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1708,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"Information\":[4.200,4,\"XB3342I2094013\",8,0.00,0.00,1.38,0.00,0.00,1]}"
	response, err := parseResponse([]byte(simpleResponse))
	require.NoError(t, err)

	assert.Equal(t, "SVRDJTTEUH", response.SN)
	assert.Equal(t, X1Boost, response.InverterType)
	assert.Equal(t, "XB3342I2094013", response.InverterSN)
	assert.Equal(t, 4.2, response.InverterTotalSize)

	inverterData, ok := response.Data.(X1BoostData)
	require.True(t, ok)
	assert.Equal(t, 733.7, inverterData.YieldTotal)
	assert.Equal(t, 0.0, inverterData.YieldToday)
}
