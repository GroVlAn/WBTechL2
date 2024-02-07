package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const (
	indexHTML = "index.html"
)

type Parser struct {
	mt      sync.Mutex
	output  string
	host    string
	urls    []*url.URL
	file    *os.File
	baseUrl *url.URL
}

func NewParser(dir string) *Parser {
	return &Parser{
		output: dir,
		urls:   make([]*url.URL, 0),
	}
}

func (p *Parser) addUrl(url *url.URL) {
	p.mt.Lock()
	p.urls = append(p.urls, url)
	p.mt.Unlock()
}

func (p *Parser) popUrl() *url.URL {
	p.mt.Lock()
	defer p.mt.Unlock()
	if len(p.urls) == 0 {
		return nil
	}
	var pURL []*url.URL

	if len(p.urls) == 1 {
		pURL = []*url.URL{p.urls[0]}
		p.urls = []*url.URL{}

		return pURL[0]
	}
	pURL, p.urls = p.urls[:1], p.urls[1:]

	return pURL[0]
}

func (p *Parser) download(curUrl *url.URL) error {
	fmt.Println(":", curUrl)
	host := curUrl.Host
	fmt.Println(host)

	path := curUrl.Path
	if len(path) == 0 || path == "/" {
		path = ""
	}
	filePath := p.baseOutputPath(path, host)

	errW := p.writeBodyToFile(curUrl, fmt.Sprintf("%s", filePath), indexHTML)
	if errW != nil {
		return fmt.Errorf("download site: %s", errW.Error())
	}
	curUrl = p.popUrl()

	return nil
}

func (p *Parser) DownloadSite() {
	err := p.download(p.baseUrl)

	if err != nil {
		fmt.Println(err)
	}
}

func (p *Parser) parseLinks(body io.Reader) ([]string, error) {
	var links []string
	buf := make([]byte, 4096)
	for {
		n, err := body.Read(buf)
		if err != nil && err != io.EOF {
			return links, fmt.Errorf("parse links: %s", err.Error())
		}

		if n == 0 {
			break
		}

		htmlContent := string(buf[:n])
		links = append(links, p.parseHref(htmlContent)...)

	}

	return links, nil
}

func (p *Parser) parseHref(htmlContent string) []string {
	var links []string

	startIndex := 0
	for {

		hrefIndex := strings.Index(htmlContent[startIndex:], "href=")
		if hrefIndex == -1 {
			break
		}

		startIndex += hrefIndex + 6
		endIndex := startIndex + strings.IndexAny(htmlContent[startIndex:], "\"'")
		if endIndex == -1 {
			break
		}

		var link string
		if startIndex < endIndex && endIndex < len(htmlContent[startIndex:]) {
			link = htmlContent[startIndex:endIndex]
		}
		links = append(links, strings.TrimSpace(link))
		startIndex = endIndex
	}

	return links
}

func (p *Parser) DownloadParts() {
	var wg sync.WaitGroup
	var countG int

	if len(p.urls) > 10 {
		countG = 10
	} else {
		countG = len(p.urls)
	}

	wg.Add(countG)
	for i := 0; i < countG; i++ {
		go func() {
			curUrl := p.popUrl()
			defer wg.Done()
			for curUrl != nil {
				fmt.Println("COUNT::::::", len(p.urls))
				fmt.Println("COUNT::::::", curUrl)
				err := p.download(curUrl)
				if err != nil {
					fmt.Println(err)
				}
				curUrl = p.popUrl()
			}
			return
		}()
	}

	wg.Wait()
}

func (p *Parser) writeBodyToFile(url *url.URL, path, fileName string) error {
	p.mt.Lock()
	resp, err := http.Get(url.String())
	p.mt.Unlock()
	defer func(resp *http.Response) {
		if resp == nil {
			return
		}
		p.mt.Lock()
		defer p.mt.Unlock()
		if err := resp.Body.Close(); err != nil {
			log.Panicf("can no close response body: %s", err.Error())
		}
	}(resp)

	if err != nil {
		return fmt.Errorf("body: %s", err.Error())
	}

	p.mkdir(path)
	file, err := os.Create(path + "/" + fileName)
	defer func(file *os.File) {
		p.mt.Lock()
		defer p.mt.Unlock()
		if err := file.Close(); err != nil {
			log.Panic("can not close file")
		}
	}(file)

	if err != nil {
		return fmt.Errorf("write to file: %s\n", err.Error())
	}
	//body, _ := io.ReadAll(resp.Body)
	//reader := strings.NewReader(string(body))

	_, errW := io.Copy(file, resp.Body)

	if errW != nil {
		return fmt.Errorf("write to file: %s\n", errW.Error())
	}

	//links := p.parseHref(string(body))
	//p.resolveURL(links)
	//p.DownloadParts()

	return err
}

func (p *Parser) resolveURL(links []string) {
	for _, link := range links {
		if len(link) == 0 {
			continue
		}
		if link[0] == '/' {
			link = p.host + link[1:]
		}
		lURL, errURL := url.Parse(link)

		if errURL != nil {
			continue
		}

		if lURL.Host == p.host {
			p.addUrl(lURL)
		}
	}
}

func (p *Parser) ReadURL() error {
	var errURL error
	args := os.Args
	if len(args) < 2 {
		return fmt.Errorf("url: url is empty\n")
	}
	p.baseUrl, errURL = url.Parse(args[len(args)-1])

	if errURL != nil {
		return fmt.Errorf("url: %s", errURL.Error())
	}

	p.host = p.baseUrl.Host

	return nil
}

func (p *Parser) baseOutputPath(path, host string) string {
	return fmt.Sprintf("%s/%s%s", p.output, host, path)
}

func (p *Parser) mkdir(path string) {
	if errMk := os.MkdirAll(path, os.ModePerm); errMk != nil {
		fmt.Println(errMk)
	}
}

func main() {
	otpFlg := flag.String("output", "res", "path to output result")
	flag.Parse()
	prs := NewParser(*otpFlg)
	errRURL := prs.ReadURL()

	if errRURL != nil {
		fmt.Print(errRURL.Error())
		return
	}

	prs.DownloadSite()
}
