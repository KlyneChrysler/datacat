package domain

import "fmt"

// CrawlingPaths never revisits: each call yields the next item in an
// endless sweep — scraper-like navigation. Next is O(1).
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
