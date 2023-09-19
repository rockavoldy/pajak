package kurs

import (
	"encoding/json"
	"time"

	"github.com/rockavoldy/pajak/currency"
)

type KursData struct {
	UpdatedAt  time.Time
	ValidFrom  string
	ValidTo    string
	Currencies []currency.Currency
}

func NewKursData() KursData {
	return KursData{
		UpdatedAt:  time.Now(),
		ValidFrom:  "",
		ValidTo:    "",
		Currencies: nil,
	}
}

func (k KursData) MarshalJSON() ([]byte, error) {
	var j struct {
		UpdatedAt  time.Time           `json:"updated_at"`
		ValidFrom  string              `json:"valid_from"`
		ValidTo    string              `json:"valid_to"`
		Currencies []currency.Currency `json:"currencies"`
	}

	j.UpdatedAt = k.UpdatedAt
	j.ValidFrom = k.ValidFrom
	j.ValidTo = k.ValidTo
	j.Currencies = k.Currencies

	return json.Marshal(j)
}

func (k *KursData) UnmarshalJSON(data []byte) error {
	var j struct {
		UpdatedAt  time.Time           `json:"updated_at"`
		ValidFrom  string              `json:"valid_from"`
		ValidTo    string              `json:"valid_to"`
		Currencies []currency.Currency `json:"currencies"`
	}

	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}

	k.UpdatedAt = j.UpdatedAt
	k.ValidFrom = j.ValidFrom
	k.ValidTo = j.ValidTo
	k.Currencies = j.Currencies

	return nil
}
