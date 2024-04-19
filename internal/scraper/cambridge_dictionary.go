package scraper

import (
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/manipulation"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gocolly/colly/v2"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
)

const cambridgeDictionaryURL = "https://dictionary.cambridge.org/dictionary/english/%s"
const cambridgeDictionaryDomainName = "dictionary.cambridge.org"

type CambridgeEnglishDictionaryData struct {
	Level        string
	Phonetic     string
	PartOfSpeech string
	AudioURL     string
	ReferenceURL string
}

func (s *Scraper) GetCambridgeEnglishDictionaryData(ctx *appcontext.AppContext, term string) (*CambridgeEnglishDictionaryData, error) {
	ctx.Logger().Info("new Cambridge English Dictionary scraping", appcontext.Fields{"term": term})

	return s.scrapeCambridgeEnglishDictionary(ctx, term, 0)
}

func (s *Scraper) scrapeCambridgeEnglishDictionary(ctx *appcontext.AppContext, term string, retryTimes int) (*CambridgeEnglishDictionaryData, error) {
	var (
		headwordFound = false
		result        = &CambridgeEnglishDictionaryData{
			ReferenceURL: fmt.Sprintf(cambridgeDictionaryURL, manipulation.Slugify(term)),
		}
	)

	c := colly.NewCollector(
		colly.AllowedDomains(cambridgeDictionaryDomainName),
	)
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  cambridgeDictionaryDomainName,
		Parallelism: 10,
		Delay:       1 * time.Second,
	})

	c.OnHTML(".headword", func(e *colly.HTMLElement) {
		if !headwordFound {
			headwordFound = true
		}
	})

	c.OnHTML(".posgram .pos", func(e *colly.HTMLElement) {
		if result.PartOfSpeech == "" {
			result.PartOfSpeech = e.Text
		}
	})

	c.OnHTML(".us .daud audio source[type='audio/mpeg']", func(e *colly.HTMLElement) {
		if result.AudioURL == "" {
			result.AudioURL = fmt.Sprintf("https://%s%s", cambridgeDictionaryDomainName, e.Attr("src"))
		}
	})

	c.OnHTML(".us .ipa", func(e *colly.HTMLElement) {
		if result.Phonetic == "" {
			result.Phonetic = e.Text
		}
	})

	c.OnHTML(".def-info .epp-xref", func(e *colly.HTMLElement) {
		if result.Level == "" {
			result.Level = e.Text
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", gofakeit.UserAgent())
		ctx.Logger().Info("Visiting url ...", appcontext.Fields{"url": r.URL.String()})
	})

	c.OnScraped(func(r *colly.Response) {
		ctx.Logger().Text("Url scraped successfully")
	})

	if err := c.Visit(result.ReferenceURL); err != nil {
		if retryTimes < 3 {
			return s.scrapeCambridgeEnglishDictionary(ctx, term, retryTimes+1)
		}
		return nil, err
	}

	if !headwordFound {
		return nil, nil
	}

	return result, nil
}
