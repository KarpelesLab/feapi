package feapi

import (
	"context"
	"errors"
	"io/fs"
	"strings"

	"github.com/KarpelesLab/apirouter"
	"github.com/KarpelesLab/contexter"
	"github.com/KarpelesLab/lngdb"
	"github.com/KarpelesLab/pjson"
	"github.com/KarpelesLab/pobj"
	"github.com/KarpelesLab/putil"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func init() {
	pobj.RegisterStatic("Language:local", languageLocal)
	pobj.RegisterActions("Language", &pobj.ObjectActions{Fetch: pobj.Static(languageGet), List: pobj.Static(languageList)})
}

type Language lngdb.Lng

func (ln Language) MarshalJSON() ([]byte, error) {
	ctx := contexter.Context()
	if ctx == nil {
		return nil, errors.New("could not fetch context")
	}
	return ln.MarshalContextJSON(ctx)
}

func (ln Language) MarshalContextJSON(ctx context.Context) ([]byte, error) {
	var curlng language.Tag
	ctx.Value(&curlng)

	l := lngdb.Lng(ln)
	lngStr := l.Tag.String()

	nameShort := putil.LookupI18N(ctx, "language_"+lngStr[:2], nil, true)
	countryName := putil.LookupI18N(ctx, "country_"+strings.ToLower(lngStr[3:]), nil, true)

	res := map[string]any{
		"Language__": lngStr,
		"Local_Name": l.LocalName,
		"Locale":     l.Locale,
		"Selected":   l.Tag == curlng,
		"Name_Short": nameShort,
		"Name_Med":   putil.Concat(nameShort, " (", countryName, ")"),
		"Name_Long":  putil.Concat(nameShort, " (", countryName, ")"),
	}

	return pjson.MarshalContext(ctx, res)
}

func languageLocal(ctx *apirouter.Context) (any, error) {
	var lngs []language.Tag
	ctx.Value(&lngs)

	var res []any
	for _, lng := range lngs {
		// grab from lngdb
		lngStr := lng.String()
		l := lngdb.Languages[lngStr]

		res = append(res, (*Language)(l))
	}
	return res, nil
}

func languageList(ctx *apirouter.Context) (any, error) {
	curlng := language.English
	ctx.Value(&curlng)
	var res []*Language

	// return all languages
	for _, l := range lngdb.Languages {
		res = append(res, (*Language)(l))
	}

	col := collate.New(curlng, collate.Loose)
	col.Sort(LanguagesSort(res))

	return res, nil
}

func languageGet(ctx *apirouter.Context, in struct{ Id string }) (any, error) {
	if in.Id == "@" {
		// grab current language
		var curlng language.Tag
		ctx.Value(&curlng)
		in.Id = curlng.String()
	}

	// find a language
	l, ok := lngdb.Languages[in.Id]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return (*Language)(l), nil
}

type LanguagesSort []*Language

func (l LanguagesSort) Len() int {
	return len(l)
}

func (l LanguagesSort) Bytes(i int) []byte {
	return []byte(l[i].LocalName)
}

func (l LanguagesSort) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
