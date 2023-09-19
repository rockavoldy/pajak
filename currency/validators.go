package currency

import "errors"

var (
	ErrNameEmpty    = errors.New("name is empty")
	ErrSymbolEmpty  = errors.New("symbol is empty")
	ErrValueEmpty   = errors.New("value is zero")
	ErrChangesEmpty = errors.New("changes is zero")
)

func validateName(name string) error {
	if len(name) == 0 {
		return ErrNameEmpty
	}
	return nil
}

func validateSymbol(symbol string) error {
	if len(symbol) == 0 {
		return ErrSymbolEmpty
	}
	return nil
}

func validateValue(value float64) error {
	if value == 0 {
		return ErrValueEmpty
	}
	return nil
}

func validateChanges(changes float64) error {
	if changes == 0 {
		return ErrChangesEmpty
	}
	return nil
}
