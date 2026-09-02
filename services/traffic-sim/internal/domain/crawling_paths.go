package domain

import "fmt"

// CrawlingPaths never revisits a path, scraper style.
type CrawlingPaths struct {
	prefix string
	next   int
}

func NewCrawlingPaths(prefix string) *CrawlingPaths {
	return &CrawlingPaths{prefix: prefix}
}

func (c *CrawlingPaths) Next() string {
	path := fmt.Sprintf("%s%d", c.prefix, c.next)
	c.next++

	return path
}
