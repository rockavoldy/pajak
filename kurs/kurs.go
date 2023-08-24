package kurs

import (
	"encoding/json"
	"time"

	"github.com/rockavoldy/pajak/currency"

	"gopkg.in/guregu/null.v4"
)

type KursData struct {
	UpdatedAt  time.Time
	ValidFrom  null.String
	ValidTo    null.String
	Currencies []currency.Currency
}

func NewKursData() KursData {
	return KursData{
		UpdatedAt:  time.Now(),
		ValidFrom:  null.String{},
		ValidTo:    null.String{},
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
	j.ValidFrom = k.ValidFrom.String
	j.ValidTo = k.ValidTo.String
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

	k = &KursData{
		UpdatedAt:  j.UpdatedAt,
		ValidFrom:  null.StringFrom(j.ValidFrom),
		ValidTo:    null.StringFrom(j.ValidTo),
		Currencies: j.Currencies,
	}

	return nil
}
