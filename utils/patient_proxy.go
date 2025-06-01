package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Patient struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// เพิ่ม fields ตามที่ response กลับมาจริง
}

func FetchPatientFromAPI(id string) (*Patient, error) {
	url := fmt.Sprintf("https://hospital-a.api.co.th/patient/search/%s", id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API responded with status %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var patient Patient
	if err := json.Unmarshal(body, &patient); err != nil {
		return nil, err
	}

	return &patient, nil
}
