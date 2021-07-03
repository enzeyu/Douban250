package parse

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Movie struct {
	Rank string
	Link string
	Title    string
	Subtitle string
	Other    string
	People     string
	Year     string
	Country     string
	Tag      string
	Star     string
	Commentpeople string
	Quote    string
}

type Page struct {
	Page int
	Url  string
}

// 获取页
func GetPages(url string) ([]Page,[]Movie) {
	// 构建请求头，使用client.Do发送请求
	client := &http.Client{}
	request, _ := http.NewRequest("GET",url,nil)
	request.Header.Add("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 OPR/66.0.3515.115")
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 字节数组s，将其转换为*goquery.Document类型
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil { log.Fatal(err) }
	dom , _ := goquery.NewDocumentFromReader(strings.NewReader(string(s)))

	pages := ParsePages(dom)
	movies25 := ParseMovies(dom)

	return pages,movies25
}

func ParsePages(doc *goquery.Document) (pages []Page) {
	pages = append(pages, Page{Page: 1, Url: ""})
	doc.Find("#wrapper>#content>div>.article>.paginator>a").Each(func(i int, s *goquery.Selection) {
		page, _ := strconv.Atoi(s.Text())
		url, _ := s.Attr("href")

		// fmt.Println(page,url)
		pages = append(pages, Page{
			Page: page,
			Url:  url,
		})
	})
	return pages
}

// 获得电影数据
func ParseMovies(doc *goquery.Document) (movie []Movie) {
	var movies []Movie
	doc.Find("#wrapper>#content>div>.article>ol>li>.item").Each(func(i int, s *goquery.Selection) {
		var people string
		var year string
		var country string
		var tag string
		rank := s.Find(".pic>em").Eq(0).Text()
		link,_ := s.Find(".pic a").Eq(0).Attr("href")
		title := s.Find(".hd a span").Eq(0).Text()
		subtitle := s.Find(".hd a span").Eq(1).Text()
		other := s.Find(".hd a span").Eq(2).Text()
		infos := s.Find(".bd>p").Eq(0).Text()
		info := strings.Split(infos, "/")

		if len(info) == 3{
			temps := strings.Split(info[0],"...")
			if len(temps) == 1{ // 序号为208，无...
				temps := strings.Split(info[0],"主演:")
				fmt.Println(temps)
				people = temps[0]
				year = temps[1]
				country = info[1]
				tag = info[2]
			}else{
				people = temps[0]
				year = temps[1]
				country = info[1]
				tag = info[2]
			}
		} else if len(info) == 4 {
			people = info[0]
			year = info[1]
			country = info[2]
			tag = info[3]
		}
		star := s.Find(".rating_num").Eq(0).Text()
		Commentpeople := s.Find(".bd .star span").Eq(3).Text()
		reg := regexp.MustCompile(`[0-9]`)
		Commentpeople = strings.Join(reg.FindAllString(Commentpeople, -1),"")
		quote := s.Find(".bd .quote").Eq(0).Text()
		movie := Movie{
			Rank: rank,
			Link: link,
			Title: title,
			Subtitle: subtitle,
			Other: other,
			People: people,
			Year: year,
			Country: country,
			Tag: tag,
			Star: star,
			Commentpeople:Commentpeople,
			Quote: quote,
		}
		log.Printf("i: %s, movie: %s", rank, movie)
		movies = append(movies,movie)
	})

	return movies
}