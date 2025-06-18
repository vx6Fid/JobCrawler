package sites

var parsers []SiteParser

func Register(p SiteParser) {
	parsers = append(parsers, p)
}

func GetParser(url string) SiteParser {
	for _, p := range parsers {
		if p.Matches(url) {
			return p
		}
	}
	return nil
}
