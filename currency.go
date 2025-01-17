package feapi

import (
	"io/fs"
	"strings"

	"github.com/KarpelesLab/apirouter"
	"github.com/KarpelesLab/currencydb"
	"github.com/KarpelesLab/pjson"
	"github.com/KarpelesLab/pobj"
)

func init() {
	pobj.RegisterActions[Currency]("Currency", &pobj.ObjectActions{
		Fetch: pobj.Static(currencyGet),
		List:  pobj.Static(currencyList),
	})
}

type Currency currencydb.Currency

func (cur Currency) MarshalJSON() ([]byte, error) {
	c := currencydb.Currency(cur)
	res := map[string]any{
		"Currency__":       c.ISO,
		"Country_ISO":      c.Country,
		"Name":             c.Name,
		"Symbol":           c.Symbol,
		"Decimals":         c.Decimals + 3,
		"Display_Decimals": c.Decimals,
		"Symbol_Position":  strings.ToLower(c.SymbolPosition.String()),
		"Virtual":          "N",
		"Visible":          "Y", // TODO only visible if part of Currency_List
	}

	return pjson.Marshal(res)
}

func currencyList(ctx *apirouter.Context) (any, error) {
	// let the system fallback
	return nil, fs.ErrNotExist
}

func currencyGet(ctx *apirouter.Context, id string) (any, error) {
	cur, ok := currencydb.All[id]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return (*Currency)(cur), nil
}
