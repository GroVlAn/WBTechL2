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

/*
indexHTML - имя файла поумолчанию
*/
const (
	indexHTML = "index.html"
)

/*
Parser - структура для парсинга сайта
mt      sync.Mutex
output  string - директория где будет лежать скачанный сайт
host    string - хост сайта
urls    []*url.URL - слайс ссылок найденых на странице
dip     int - глубина парсинга, 0 - только переданную страницу, 1 - все ссылки, что будут найдены на странице
baseUrl *url.URL - страница переданная в качестве аргумента
*/
type Parser struct {
	mt      sync.Mutex
	output  string
	host    string
	urls    []*url.URL
	dip     int
	baseUrl *url.URL
}

func NewParser(dir string, dip int) *Parser {
	return &Parser{
		output: dir,
		urls:   make([]*url.URL, 0),
		dip:    dip,
	}
}

// DownloadSite - метод запускающий парсинг текущей страинци
func (p *Parser) DownloadSite() {
	err := p.download(p.baseUrl)

	if err != nil {
		fmt.Println(err)
	}
}

// DownloadParts - метод запускающий скачивание внутрениз ссылок
func (p *Parser) DownloadParts() {
	var wg sync.WaitGroup
	var countG int

	if len(p.urls) > 10 {
		countG = 10
	} else {
		countG = len(p.urls)
	}
	// запускаем и ждём максимум 10 горутин, которые будут скачивать файлы с сайта
	wg.Add(countG)
	for i := 0; i < countG; i++ {
		go func() {
			curUrl := p.popUrl()
			fmt.Println(":::", p.getDip())
			fmt.Println(":::", len(p.urls))
			defer wg.Done()
			for curUrl != nil {
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

// ReadURL - метод читающий url из консоли
func (p *Parser) ReadURL() error {
	var errURL error
	args := os.Args
	if len(args) < 2 {
		return fmt.Errorf("url: url is empty\n")
	}
	p.baseUrl, errURL = url.Parse(args[len(args)-1])

	if errURL != nil {
		p.baseUrl, errURL = url.Parse(args[len(args)-2])
		if errURL != nil {
			return fmt.Errorf("url: %s", errURL.Error())
		}
	}

	p.host = p.baseUrl.Host

	return nil
}

/*
writeBodyToFile - метод пишущий в файл и читающий ссылки на странице
url *url.URL - ссылка на файл
path string - путь до файла
fileName string - имя файла
*/
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
	body, _ := io.ReadAll(resp.Body)
	reader := strings.NewReader(string(body))

	_, errW := io.Copy(file, reader)

	if errW != nil {
		return fmt.Errorf("write to file: %s\n", errW.Error())
	}

	if p.getDip() > 0 {
		p.degreeDip()
		links := p.parseHref(string(body))
		p.resolveURL(links)
		p.DownloadParts()
	}

	return err
}

/*
resolveURL - метод добавляющий к списку дочерних ссылок, только те, что являются ссылкой и принадлежат данному хосту
links []string - все найденые ссылки
*/
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

/*
parseHref - метод для парсингла href
htmlContent string -
*/
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

/*
baseOutputPath - метод для формирования базового пути до файла
path string - текущий путь
host string - хост ссылки
*/
func (p *Parser) baseOutputPath(path, host string) string {
	return fmt.Sprintf("%s/%s%s", p.output, host, path)
}

/*
mkdir - метод для создания директории
path string - путь до директории
*/
func (p *Parser) mkdir(path string) {
	if errMk := os.MkdirAll(path, os.ModePerm); errMk != nil {
		fmt.Println(errMk)
	}
}

/*
degreeDip - метод для уменьшения значения глубины
*/
func (p *Parser) degreeDip() {
	p.mt.Lock()
	defer p.mt.Unlock()
	if p.dip > 1 {
		p.dip = 0
		return
	}
	p.dip--
}

/*
getDip - метод для получения текущей глубины
*/
func (p *Parser) getDip() int {
	p.mt.Lock()
	defer p.mt.Unlock()
	return p.dip
}

// addUrl - метод для добавления ссылки к текущему списку
func (p *Parser) addUrl(url *url.URL) {
	p.mt.Lock()
	p.urls = append(p.urls, url)
	p.mt.Unlock()
}

// popUrl - метод достаёт первую ссылку в списке
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

// download - метод для скачивания файла, он поределяет путь по которому нужно записать файл и запускает запись в файл
func (p *Parser) download(curUrl *url.URL) error {
	host := curUrl.Host

	path := curUrl.Path
	if len(path) == 0 || path == "/" {
		path = ""
	}
	filePath := p.baseOutputPath(path, host)
	lastSlash := strings.LastIndex(filePath, "/")

	nameFile := indexHTML

	if lastSlash != -1 {
		if strings.Contains(filePath[lastSlash+1:], ".") && filePath[lastSlash+1:] != p.host {
			nameFile = filePath[lastSlash+1:]
			filePath = filePath[:lastSlash]
		}
	}

	errW := p.writeBodyToFile(curUrl, fmt.Sprintf("%s", filePath), nameFile)
	if errW != nil {
		return fmt.Errorf("download site: %s", errW.Error())
	}

	return nil
}

func main() {
	otpFlg := flag.String("output", "res", "path to output result")
	dipFlg := flag.Int("dip", 0, "path to output result")
	flag.Parse()
	prs := NewParser(*otpFlg, *dipFlg)
	errRURL := prs.ReadURL()

	if errRURL != nil {
		fmt.Print(errRURL.Error())
		return
	}

	prs.DownloadSite()
}
