package currency

import (
	"encoding/json"
)

type Currency struct {
	Name    string
	Symbol  string
	Value   float64
	Changes float64
}

func NewCurrency(name, symbol string, value, changes float64) (Currency, error) {
	if err := validateName(name); err != nil {
		return Currency{}, err
	}
	if err := validateSymbol(symbol); err != nil {
		return Currency{}, err
	}
	if err := validateValue(value); err != nil {
		return Currency{}, err
	}
	// FIXME: check later, if no changes should be considered or not
	// if err := validateChanges(changes); err != nil {
	// 	return Currency{}, err
	// }

	currency := Currency{
		Name:    name,
		Symbol:  symbol,
		Value:   value,
		Changes: changes,
	}

	return currency, nil
}

func (c Currency) MarshalJSON() ([]byte, error) {
	var j struct {
		Name    string  `json:"name"`
		Symbol  string  `json:"symbol"`
		Value   float64 `json:"value"`
		Changes float64 `json:"changes"`
	}

	j.Name = c.Name
	j.Symbol = c.Symbol
	j.Value = c.Value / 100
	j.Changes = c.Changes / 100

	return json.Marshal(j)
}

func (c *Currency) UnmarshalJSON(data []byte) error {
	var j struct {
		Name    string  `json:"name"`
		Symbol  string  `json:"symbol"`
		Value   float64 `json:"value"`
		Changes float64 `json:"changes"`
	}

	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}

	c.Name = j.Name
	c.Symbol = j.Symbol
	c.Value = j.Value * 100
	c.Changes = j.Changes * 100

	return nil
}
