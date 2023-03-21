package solax_local_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type InverterType int

const (
	X1Mini  InverterType = 1
	X1Boost InverterType = 2
	X3                   = 3
	X3G4                 = 4
)

func (i InverterType) MarshalJSON() ([]byte, error) {
	switch i {
	case X1Mini:
		return json.Marshal("X1Mini")
	case X1Boost:
		return json.Marshal("X1Boost")
	}
	return []byte("Unknown"), nil
}

type RawInverterResponse struct {
	SN          string        `json:"sn"`
	Version     string        `json:"ver"`
	Type        int           `json:"type"`
	Data        []int         `json:"Data"`
	Information []interface{} `json:"Information"`
}

type InverterResult struct {
	SN                string
	Version           string
	InverterType      InverterType
	InverterSN        string
	InverterTotalSize float64
	Data              interface{}
}

type X1BoostData struct {
	YieldToday float64
	YieldTotal float64
}

var (
	inverterTotalSizeField = map[InverterType]int{X1Boost: 0}
	inverterSnField        = map[InverterType]int{X1Boost: 2}
)

func Request(host string, password string) (InverterResult, error) {
	form := url.Values{}
	form.Add("optType", "ReadRealTimeData")
	form.Add("pwd", password)

	//goland:noinspection HttpUrlsUsage
	response, err := http.PostForm(fmt.Sprintf("http://%s", host), form)
	if err != nil {
		return InverterResult{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return InverterResult{}, err
	}

	return parseResponse(responseBody)
}

func parseResponse(responseBody []byte) (InverterResult, error) {
	var rawResponse RawInverterResponse
	if err := json.Unmarshal(responseBody, &rawResponse); err != nil {
		return InverterResult{}, err
	}

	inverterType := parseInverterType(rawResponse.Type)
	inverterSN := rawResponse.Information[inverterSnField[inverterType]].(string)
	inverterTotalSize := rawResponse.Information[inverterTotalSizeField[inverterType]].(float64)

	var inverterData interface{}
	switch inverterType {
	case X1Boost:
		inverterData = parseX1Boost(rawResponse)
	default:
		return InverterResult{}, errors.New("invalid inverter type")
	}

	return InverterResult{
		SN:                rawResponse.SN,
		Version:           rawResponse.Version,
		InverterType:      inverterType,
		InverterSN:        inverterSN,
		InverterTotalSize: inverterTotalSize,
		Data:              inverterData,
	}, nil
}

func parseX1Boost(rawResponse RawInverterResponse) X1BoostData {
	return X1BoostData{
		YieldToday: float64(rawResponse.Data[13]) * 0.1,
		YieldTotal: float64(rawResponse.Data[11]) * 0.1,
	}
}

func parseInverterType(typeField int) InverterType {
	switch typeField {
	case 4:
		return X1Boost
	case 5:
	case 6:
	case 7:
	case 16:
		return X3
	case 14:
	case 15:
		return X3G4
	default:
		return X1Mini
	}
	return X1Mini
}
