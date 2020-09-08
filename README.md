<!-- +
title: about md head parser
urlid: about-md-head-parser
overview: md-head-parser is a extention of the markdown format. it defines a header to describe the markdown
cate: tools
tags: #md head parser, #markdown
lang: en
published_at: 2020-09-06T12:00:00Z
updated_at: 2020-09-06T12:00:00Z
+ -->

## introduction

sometimes we need more information to describe a markdown content. such as it's title, language of the content and so on. but there is still no any solution to carry such information. this is a package that define a header of the markdown file. the header can contain all the information you wanna carry but don't wanna make them visible. we will introduce the header specified.

## syntax of header

it looks like below

```
<!-- + 
title: test title
overview: overview of the content
cate: sample
tags: #hello world, #tag1, #tag2
lang: en
published_at: 2020-08-08T08:00:00Z
updated_at: 2020-08-08T08:00:00Z
+ -->

```

I choose `<!-- + + -->` as a wrapper because it can make the header hidden. the content of the header just like a http header presentation. the content is composed by key-value pairs which seperated by `:`, the whitespace around `:` will be ignored. the key of the key-value pair only accept number, alpha, -, _ and must lead with alpha. normally the has no limitation but some exception. the value limitation of specified key list below.

* urlid: only accept alpha, number, _, -
* tags: tags will begin with a `#`, and seperated by `,`, whitespace around `,` will be ignored
* published_at: must a `RFC3339` time format string

## usage

```golang

import md

f, e := os.Open("./blog.md")
if e != nil {
	log.Print(e)
	return
}

// this will parse the input then output the head and body without head
head, body, err := md.Parse(f)

// this will parse the input then output the head.
head, err := md.ParseHead(f)
```

## API

```golang

type MdHead struct {
	Title       string
	Urlid       string
	Overview    string
	Cate        string
	Lang        string
	Tags        []string
	Image       string
	PublishedAt time.Time
	UpdatedAt   time.Time
	Extra       map[string]string
}

func Parse(r io.Reader) (head MdHead, body []byte, error);

func ParseHead(r io.Reader) (head MdHead, error);

```
