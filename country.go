package feapi

import (
	"context"
	"errors"
	"io/fs"
	"sort"

	"github.com/KarpelesLab/apirouter"
	"github.com/KarpelesLab/contexter"
	"github.com/KarpelesLab/countrydb"
	"github.com/KarpelesLab/countrydb/countrynames"
	"github.com/KarpelesLab/pjson"
	"github.com/KarpelesLab/pobj"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func init() {
	pobj.RegisterActions("Country", &pobj.ObjectActions{Fetch: pobj.Static(countryGet), List: pobj.Static(countryList)})
}

type Country countrydb.Country

func (country Country) Export(curlng language.Tag) map[string]any {
	c := countrydb.Country(country)
	lngInfo, found := countrynames.LocaleByTag[curlng]
	if !found {
		lngInfo = countrynames.English
	}

	res := map[string]any{
		"Country__":      c.ISO3166_Alpha2,
		"Name":           lngInfo[c.ISO3166_Alpha2].Name,
		"ISO3166_Code":   c.ISO3166_Alpha2,
		"ISO3166_3_Code": c.ISO3166_Alpha3,
		"Phone_Prefix":   []string{c.PhonePrefix},
		"Display_Format": country.getDisplayFormat(),
		// new stuff?
		"id":   c.ISO3166_Alpha2,
		"iso2": c.ISO3166_Alpha2,
		"iso3": c.ISO3166_Alpha3,
		"name": lngInfo[c.ISO3166_Alpha2].Name,
	}
	return res
}

func (country Country) MarshalJSON() ([]byte, error) {
	ctx := contexter.Context()
	if ctx == nil {
		return nil, errors.New("could not fetch context")
	}
	return country.MarshalContextJSON(ctx)
}

func (country Country) MarshalContextJSON(ctx context.Context) ([]byte, error) {
	var curlng language.Tag
	ctx.Value(&curlng)

	return pjson.MarshalContext(ctx, country.Export(curlng))
}

func (country Country) getDisplayFormat() [][]string {
	var curlng language.Tag
	ctx := contexter.Context()
	if ctx != nil {
		ctx.Value(&curlng)
	}

	nameSuffix := "!"
	if curlng == language.Japanese {
		nameSuffix = "!様"
	}

	switch country.ISO3166_Alpha2 {
	case "US":
		return [][]string{
			[]string{"First_Name", "Middle_Name", "Last_Name", nameSuffix},
			[]string{"!(", "Nickname", "!)"},
			[]string{"Company_Name"},
			[]string{"Company_Department"},
			[]string{"Address", "!", "Address1"},
			[]string{"Address2"},
			[]string{"Address3"},
			[]string{"City", "!, ", "Province", "Zip"},
			[]string{"Country__", "!", "Country"},
		}
	case "JP":
		return [][]string{
			[]string{"!〒", "Zip", "Province", "!", "City", "!", "Address", "!", "Address1"},
			[]string{"Address2"},
			[]string{"Company_Name", "!御中"},
			[]string{"Company_Department"},
			[]string{"Last_Name", "First_Name", nameSuffix},
			[]string{"!(", "Nickname", "!)"},
			[]string{"Country__", "!", "Country"},
		}
	default:
		// most EU countries, etc
		return [][]string{
			[]string{"First_Name", "Middle_Name", "Last_Name", nameSuffix},
			[]string{"!(", "Nickname", "!)"},
			[]string{"Company_Name"},
			[]string{"Company_Department"},
			[]string{"Address", "!", "Address1"},
			[]string{"Address2"},
			[]string{"Address3"},
			[]string{"Zip", "City"},
			[]string{"Province"},
			[]string{"Country__", "!", "Country"},
		}
	}
}

func countryList(ctx *apirouter.Context) (any, error) {
	var curlng language.Tag
	ctx.Value(&curlng)

	var res Sortable
	col := collate.New(curlng, collate.Loose)
	b := &collate.Buffer{}

	// return all countries
	for _, l := range countrydb.All {
		v := (*Country)(l).Export(curlng)
		k := col.KeyFromString(b, v["Name"].(string))
		res = append(res, &SortableValue{V: v, K: k})
	}

	sort.Sort(res)

	return res, nil
}

func countryGet(ctx *apirouter.Context, in struct{ Id string }) (any, error) {
	// find a country
	c, ok := countrydb.ByAlpha2[in.Id]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return (*Country)(c), nil
}
