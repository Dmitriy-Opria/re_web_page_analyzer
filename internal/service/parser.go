package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/log"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/model"
	"github.com/PuerkitoBio/goquery"
)

type ParserRepository interface {
	Parse(ctx context.Context) error
}

type ParserService struct {
	fetcher     Fetcher
	workerCount int

	sync.WaitGroup
}

func NewParserService(fetcher Fetcher, workerCount int) *ParserService {
	return &ParserService{
		fetcher:     fetcher,
		workerCount: workerCount,
	}
}

func (p *ParserService) Parse(ctx context.Context, url string) (*model.ParserResponse, error) {
	// Load the HTML document
	doc, err := p.fetcher.Fetch(url)
	if err != nil {
		log.Fatal(err)
	}

	internalLink, externalLink := p.setInternalLink(ctx, doc)
	return &model.ParserResponse{
		Version:       p.version(doc),
		Title:         p.title(doc),
		ListH1:        p.header(doc, 1),
		ListH2:        p.header(doc, 2),
		ListH3:        p.header(doc, 3),
		ListH4:        p.header(doc, 4),
		ListH5:        p.header(doc, 5),
		ListH6:        p.header(doc, 6),
		InternalLinks: internalLink,
		ExternalLinks: externalLink,
		Login:         p.login(doc),
	}, nil
}

func (p *ParserService) version(doc *goquery.Document) string {
	body, err := doc.Html()
	if err != nil {
		return "undefined"
	}
	var versionHTML string
	if len(body) > strings.Index(body, "<!") && len(body) > strings.Index(body, ">")+1 {
		versionHTML = body[strings.Index(body, "<!") : strings.Index(body, ">")+1]
	}

	switch {
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD HTML 4.01 Frameset//EN")):
		return "HTML 4.01 Frameset"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD HTML 4.01 Transitional//EN")):
		return "HTML 4.01 Transitional"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD HTML 4.01//EN")):
		return "HTML 4.01 Strict"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD XHTML 1.0 Frameset//EN")):
		return "XHTML 1.0 Frameset"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD XHTML 1.0 Transitional//EN")):
		return "XHTML 1.0 Transitional"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD XHTML 1.0 Strict//EN")):
		return "XHTML 1.0 Strict"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD XHTML 1.1//EN")):
		return "XHTML 1.1 DTD"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD XHTML Basic 1.1//EN")):
		return "XHTML Basic 1.1"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD XHTML 1.1 plus MathML 2.0 plus SVG 1.1//EN")):
		return "XHTML + MathML + SVG - DTD"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD XHTML 1.1 plus MathML 2.0 plus SVG 1.1//EN")):
		return "XHTML + MathML + SVG Profile (XHTML as the host language) - DTD"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD XHTML 1.1 plus MathML 2.0 plus SVG 1.1//EN")):
		return "XHTML + MathML + SVG Profile (Using SVG as the host) - DTD"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD SVG 1.1//EN")):
		return "SVG 1.1 Full - DTD"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("W3C//DTD SVG 1.0//EN")):
		return "SVG 1.0 - DTD:"
	case strings.Contains(strings.ToLower(versionHTML), strings.ToLower("!DOCTYPE HTML")):
		return "HTML5 and beyond"
	default:
		return "undefined"
	}
}

func (p *ParserService) title(doc *goquery.Document) string {
	var title string
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		// For each item found, get title
		title = s.Text()
	})
	return title
}

func (p *ParserService) header(doc *goquery.Document, level int) []string {
	var headers []string
	doc.Find(fmt.Sprintf("h%v", level)).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the header
		headers = append(headers, strings.Join(strings.Fields(strings.TrimSpace(s.Text())), " "))
	})
	return headers
}

func (p *ParserService) setInternalLink(ctx context.Context, doc *goquery.Document) ([]*model.Link, []*model.Link) {
	var internalLink []*model.Link
	var externalLink []*model.Link

	host := doc.Url.Host
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get title
		if href, _ := s.Attr("href"); !strings.Contains(href, "javascript") {
			link := model.Link{
				Name: strings.TrimSpace(
					strings.Replace(
						strings.Replace(
							s.Text(), "\n", "", -1),
						"\t", "", -1),
				),
				Url: href,
			}
			if len([]rune(href)) > 0 {
				if strings.Contains(href, host) {
					internalLink = append(internalLink, &link)
				} else if []rune(href)[0] == '/' || []rune(href)[0] == '#' || []rune(href)[0] == ':' {

					link.Url = fmt.Sprintf("http://%s/%s", host,
						strings.Replace(
							strings.Replace(link.Url,
								"#", "/", -1),
							"/", "", -1))
					internalLink = append(internalLink, &link)
				} else {
					externalLink = append(externalLink, &link)
				}
			}
		}
	})

	p.checkLinks(ctx, internalLink)
	p.checkLinks(ctx, externalLink)

	return internalLink, externalLink
}

func (p *ParserService) checkLinks(ctx context.Context, links []*model.Link) {
	jobPool := make(chan struct{}, p.workerCount)

	for index, link := range links {
		p.Add(1)
		jobPool <- struct{}{}
		go func() {
			defer func() {
				<-jobPool
				p.Done()
			}()
			result, err := p.fetcher.IsAccessible(ctx, &model.WorkerWrapper{
				Index: index,
				Url:   link.Url,
			})
			if err != nil {
				log.Error(err)
				return
			}
			links[result.Index].Accessible = result.Result
		}()
		time.Sleep(50 * time.Millisecond)
	}
	p.Wait()
}

func (p *ParserService) login(doc *goquery.Document) bool {
	var formFields []string
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		formFields = append(formFields, processChildren(s)...)
	})
	return model.IsLoginForm(formFields)
}

func processChildren(s *goquery.Selection) []string {
	result := make([]string, 0)
	if s.Children().Length() > 0 {
		s.Children().Each(func(i int, si *goquery.Selection) {
			result = append(result, processChildren(si)...)
		})
	} else {
		if text := s.Text(); len(strings.TrimSpace(text)) > 0 {
			result = append(result, text)
		}
	}
	return result
}
