package currency

import "encoding/json"

type Currency struct {
	Name    string
	Symbol  string
	Value   int
	Changes int
}

func NewCurrency(name, symbol string, value, changes int) (Currency, error) {
	if err := validateName(name); err != nil {
		return Currency{}, err
	}
	if err := validateSymbol(symbol); err != nil {
		return Currency{}, err
	}
	if err := validateValue(value); err != nil {
		return Currency{}, err
	}
	if err := validateChanges(changes); err != nil {
		return Currency{}, err
	}

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
		Name    string `json:"name"`
		Symbol  string `json:"symbol"`
		Value   int    `json:"value"`
		Changes int    `json:"changes"`
	}

	j.Name = c.Name
	j.Symbol = c.Symbol
	j.Value = c.Value / 100
	j.Changes = c.Changes / 100

	return json.Marshal(j)
}
