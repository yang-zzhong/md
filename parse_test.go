package md

import (
	"log"
	"os"
	"testing"
)

func TestParseBlog(t *testing.T) {
	f, e := os.Open("./blog.md")
	if e != nil {
		log.Print(e)
		return
	}
	head, body, err := Parse(f)

	if err != nil {
		log.Printf("err: %s\n", err.Error())
	}

	log.Printf("title: %s\n", head.Title)
	log.Printf("overview: %s\n", head.Overview)
	log.Printf("urlid: %s\n", head.Urlid)
	log.Printf("tags: %v\n", head.Tags)
	log.Printf("lang: %s\n", head.Lang)
	log.Printf("cate: %s\n", head.Cate)
	log.Printf("published_at: %s\n", head.PublishedAt)
	log.Printf("updated_at: %s\n", head.UpdatedAt)
	log.Printf("body: %s\n", string(body))
}
