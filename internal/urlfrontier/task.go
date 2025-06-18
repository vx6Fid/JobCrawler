package urlfrontier

type CrawlTask struct {
	URL  string
	Type string
	Meta map[string]string
}
