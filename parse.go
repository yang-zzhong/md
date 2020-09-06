package md

// this ia a package to extend markdown that add a header as the summary info about the markdown content
// the praser parse a markdown content like below
//  <!-- +
//  title: test title
//  urlid: test-title
//  overview: this is a sample markdown that the tool will parse
//  cate: sample of usage
//  tags: #sample, #test, #parse
//  image: test-title.png
//  published_at: 2020-09-06 12:23:00
//  updated_at: 2020-10-14 12:23:00
//  + -->
//  ## title
//  content
//  ## sub title
//  sub content
//
// the content wrapped in <!-- + content + --> will treat as head of the markdown
import (
	"errors"
	"io"
	"strings"
	"time"
)

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

const (
	shb1 = iota // <
	shb2        // !
	shb3        // -
	shb4        // -
	shb5        // \s
	she1        // +
	she2        // \s
	she3        // -
	she4        // -
	she5        // >
	sb          // body
	sh          // head
	sk          // head key
	sak         // after key
	sv          // head value

	title       = "title"
	urlid       = "urlid"
	overview    = "overview"
	cate        = "cate"
	category    = "category"
	tags        = "tags"
	image       = "image"
	img         = "img"
	publishedAt = "published_at"
	updatedAt   = "updated_at"
	lang        = "lang"
)

func Parse(r io.Reader) (MdHead, []byte, error) {
	p := parser{}
	return p.parse(r, true)
}

func ParseHead(r io.Reader) (head MdHead, err error) {
	p := parser{}
	head, _, err = p.parse(r, false)
	return
}

type parser struct {
	buf      []byte
	cache    []byte
	state    int
	key, val []byte
}

func (p *parser) parse(r io.Reader, withbody bool) (head MdHead, body []byte, err error) {
	head = MdHead{}
	head.Extra = make(map[string]string)
	body = []byte{}
	err = nil
	p.buf = make([]byte, 1)
	p.cache = []byte{}
	p.state = sb
	for {
		if _, err = r.Read(p.buf); err != nil {
			if err == io.EOF {
				if p.state != she5 && withbody {
					body = append(body, p.cache...)
				}
				err = nil
			}
			return
		}
		switch p.state {
		case sb:
			p.inbody(withbody, &body)
		case shb1:
			p.inshb('!', shb2, withbody, &body)
		case shb2:
			p.inshb('-', shb3, withbody, &body)
		case shb3:
			p.inshb('-', shb4, withbody, &body)
		case shb4:
			p.inshb(' ', shb5, withbody, &body)
		case shb5:
			p.inshb('+', sh, withbody, &body)
		case sh:
			p.insh(withbody)
		case sk:
			p.insk(withbody)
		case sak:
			p.insak(withbody)
		case sv:
			if err = p.insv(withbody, &head); err != nil {
				return
			}
		case she1:
			p.inshb(' ', she2, withbody, &body)
		case she2:
			p.inshb('-', she3, withbody, &body)
		case she3:
			p.inshb('-', she4, withbody, &body)
		case she4:
			p.inshb('>', she5, withbody, &body)
		case she5:
			if withbody {
				tmp := make([]byte, 2048)
				len := 0
				if len, err = r.Read(tmp); err != nil {
					return
				}
				body = append(body, tmp[:len]...)
			}
		}
	}
	return
}

func (p *parser) inbody(withbody bool, body *[]byte) {
	if p.buf[0] == '<' {
		if withbody {
			p.cache = []byte{'<'}
		}
		p.state = shb1
		return
	}
	if withbody {
		(*body) = append(*body, p.cache...)
	}
}

func (p *parser) inshb(b byte, s int, withbody bool, body *[]byte) {
	if p.buf[0] == b {
		if withbody {
			p.cache = append(p.cache, p.buf[0])
		}
		p.state = s
		return
	}
	if withbody {
		(*body) = append(*body, p.buf[0])
		p.cache = []byte{}
	}
	p.state = sb
}

func (p *parser) insh(withbody bool) {
	if withbody {
		p.cache = append(p.cache, p.buf[0])
	}
	if isAlpha(p.buf[0]) {
		p.key = []byte{p.buf[0]}
		p.state = sk
		return
	} else if p.buf[0] == '+' {
		p.state = she1
	}
}

func (p *parser) insk(withbody bool) {
	if withbody {
		p.cache = append(p.cache, p.buf[0])
	}
	if isW(p.buf[0]) || isC(p.buf[0]) {
		p.key = append(p.key, p.buf[0])
		return
	}
	if p.buf[0] == ':' {
		p.state = sak
		return
	}
	p.state = sh
}

func (p *parser) insak(withbody bool) {
	if withbody {
		p.cache = append(p.cache, p.buf[0])
	}
	if !isWhiteSpace(p.buf[0]) {
		p.val = []byte{p.buf[0]}
		p.state = sv
	}
}

func (p *parser) insv(withbody bool, head *MdHead) error {
	if withbody {
		p.cache = append(p.cache, p.buf[0])
	}
	if p.buf[0] == '\n' {
		e := p.handlePair(head)
		p.state = sh
		return e
	}
	p.val = append(p.val, p.buf[0])
	return nil
}

func is(k string, p string) bool {
	return strings.ToLower(k) == p
}

func (p *parser) handlePair(head *MdHead) error {
	k := string(p.key)
	if is(k, title) {
		head.Title = string(p.val)
		return nil
	} else if is(k, urlid) {
		return handleUrlid(string(p.val), head)
	} else if is(k, tags) {
		handleTags(string(p.val), head)
		return nil
	} else if is(k, cate) || is(k, category) {
		head.Cate = string(p.val)
		return nil
	} else if is(k, overview) {
		head.Overview = string(p.val)
		return nil
	} else if is(k, lang) {
		head.Lang = string(p.val)
		return nil
	} else if is(k, image) || is(k, img) {
		head.Image = string(p.val)
		return nil
	} else if is(k, publishedAt) {
		var e error
		head.PublishedAt, e = time.Parse(time.RFC3339, strings.Trim(string(p.val), " "))
		return e
	} else if is(k, updatedAt) {
		var e error
		head.UpdatedAt, e = time.Parse(time.RFC3339, strings.Trim(string(p.val), " "))
		return e
	}
	head.Extra[k] = string(p.val)
	return nil
}

func handleUrlid(v string, head *MdHead) error {
	r := strings.NewReader(v)
	buf := make([]byte, 1)
	temp := []byte{}
	for {
		if l, e := r.Read(buf); e != nil {
			if e == io.EOF {
				head.Urlid = string(temp)
				return nil
			}
			return e
		} else if l == 0 {
			continue
		} else if isW(buf[0]) || buf[0] == '_' || buf[0] == '-' {
			temp = append(temp, buf[0])
		} else {
			return errors.New("urlid only support number, alpha, _ and -")
		}
	}
}

func handleTags(v string, head *MdHead) {
	tags := strings.Split(v, "#")
	for _, tag := range tags {
		t := strings.Trim(tag, " ,")
		if t != "" {
			head.Tags = append(head.Tags, t)
		}
	}
}

func isNumber(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlpha(b byte) bool {
	return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z'
}

func isW(b byte) bool {
	return isNumber(b) || isAlpha(b)
}

func isC(b byte) bool {
	return b == '_' || b == '-'
}

func isWhiteSpace(b byte) bool {
	return b == '\t' || b == ' ' || b == '\n'
}
