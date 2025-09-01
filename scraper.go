package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	client *http.Client
}

func NewScraper() *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *Scraper) ScrapeKaikatsuClub() ([]NetCafe, error) {
	url := "https://www.kaikatsu.jp/shop/tokyo/"
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var cafes []NetCafe
	
	doc.Find(".shop-list-item").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".shop-name").Text())
		if name == "" {
			name = strings.TrimSpace(s.Find("h3").Text())
		}
		
		address := strings.TrimSpace(s.Find(".shop-address").Text())
		if address == "" {
			address = strings.TrimSpace(s.Find(".address").Text())
		}
		
		hours := "24時間営業"
		if hoursText := s.Find(".shop-hours").Text(); hoursText != "" {
			hours = strings.TrimSpace(hoursText)
		}
		
		phone := strings.TrimSpace(s.Find(".shop-tel").Text())
		if phone == "" {
			phone = strings.TrimSpace(s.Find(".tel").Text())
		}
		
		shopURL, _ := s.Find("a").Attr("href")
		if !strings.HasPrefix(shopURL, "http") && shopURL != "" {
			shopURL = "https://www.kaikatsu.jp" + shopURL
		}
		
		if name != "" {
			cafes = append(cafes, NetCafe{
				Name:     "快活CLUB " + name,
				Location: address,
				Hours:    hours,
				Phone:    phone,
				URL:      shopURL,
			})
		}
	})

	if len(cafes) == 0 {
		cafes = s.scrapeKaikatsuAlternative(doc)
	}

	return cafes, nil
}

func (s *Scraper) scrapeKaikatsuAlternative(doc *goquery.Document) []NetCafe {
	var cafes []NetCafe
	
	doc.Find("li").Each(func(i int, sel *goquery.Selection) {
		text := sel.Text()
		if strings.Contains(text, "店") || strings.Contains(text, "快活") {
			lines := strings.Split(text, "\n")
			name := ""
			address := ""
			phone := ""
			
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				
				if strings.Contains(line, "店") && name == "" {
					name = line
				} else if (strings.Contains(line, "区") || strings.Contains(line, "市")) && address == "" {
					address = line
				} else if matched, _ := regexp.MatchString(`\d{2,4}-\d{2,4}-\d{4}`, line); matched {
					phone = line
				}
			}
			
			if name != "" {
				if !strings.HasPrefix(name, "快活CLUB") {
					name = "快活CLUB " + name
				}
				cafes = append(cafes, NetCafe{
					Name:     name,
					Location: address,
					Hours:    "24時間営業",
					Phone:    phone,
					URL:      "https://www.kaikatsu.jp/",
				})
			}
		}
	})
	
	return cafes
}

func (s *Scraper) ScrapeJiqoo() ([]NetCafe, error) {
	url := "https://jiqoo.jp/shop/?pref=13"
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var cafes []NetCafe
	
	doc.Find(".shop-item, .store-item").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".shop-name, .store-name, h3").Text())
		address := strings.TrimSpace(s.Find(".shop-address, .store-address, .address").Text())
		phone := strings.TrimSpace(s.Find(".shop-tel, .store-tel, .tel").Text())
		
		hours := "24時間営業"
		if hoursText := s.Find(".shop-hours, .hours").Text(); hoursText != "" {
			hours = strings.TrimSpace(hoursText)
		}
		
		shopURL, _ := s.Find("a").Attr("href")
		if !strings.HasPrefix(shopURL, "http") && shopURL != "" {
			shopURL = "https://jiqoo.jp" + shopURL
		}
		
		if name != "" {
			if !strings.Contains(name, "自遊空間") {
				name = "自遊空間 " + name
			}
			cafes = append(cafes, NetCafe{
				Name:     name,
				Location: address,
				Hours:    hours,
				Phone:    phone,
				URL:      shopURL,
			})
		}
	})

	return cafes, nil
}

func (s *Scraper) ScrapeManboo() ([]NetCafe, error) {
	url := "https://www.manboo.co.jp/shop/"
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var cafes []NetCafe
	
	doc.Find(".shop-list-item, .store-item, li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "東京") || strings.Contains(text, "店") {
			name := ""
			address := ""
			phone := ""
			
			nameElem := s.Find(".shop-name, h3, strong").First()
			if nameElem.Length() > 0 {
				name = strings.TrimSpace(nameElem.Text())
			}
			
			addressElem := s.Find(".address, .shop-address").First()
			if addressElem.Length() > 0 {
				address = strings.TrimSpace(addressElem.Text())
			}
			
			phoneElem := s.Find(".tel, .phone").First()
			if phoneElem.Length() > 0 {
				phone = strings.TrimSpace(phoneElem.Text())
			}
			
			if name == "" && strings.Contains(text, "店") {
				lines := strings.Split(text, "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if strings.Contains(line, "店") && name == "" {
						name = line
					}
				}
			}
			
			if name != "" && strings.Contains(text, "東京") {
				if !strings.Contains(name, "マンボー") {
					name = "マンボー " + name
				}
				cafes = append(cafes, NetCafe{
					Name:     name,
					Location: address,
					Hours:    "24時間営業",
					Phone:    phone,
					URL:      "https://www.manboo.co.jp/",
				})
			}
		}
	})

	return cafes, nil
}

func (s *Scraper) ScrapeAll() ([]NetCafe, error) {
	var allCafes []NetCafe
	var errors []string

	fmt.Println("快活CLUBの店舗情報を取得中...")
	kaikatsu, err := s.ScrapeKaikatsuClub()
	if err != nil {
		errors = append(errors, fmt.Sprintf("快活CLUB: %v", err))
		log.Printf("Error scraping Kaikatsu: %v", err)
	} else {
		allCafes = append(allCafes, kaikatsu...)
		fmt.Printf("  → %d店舗を取得\n", len(kaikatsu))
	}

	fmt.Println("自遊空間の店舗情報を取得中...")
	jiqoo, err := s.ScrapeJiqoo()
	if err != nil {
		errors = append(errors, fmt.Sprintf("自遊空間: %v", err))
		log.Printf("Error scraping Jiqoo: %v", err)
	} else {
		allCafes = append(allCafes, jiqoo...)
		fmt.Printf("  → %d店舗を取得\n", len(jiqoo))
	}

	fmt.Println("マンボーの店舗情報を取得中...")
	manboo, err := s.ScrapeManboo()
	if err != nil {
		errors = append(errors, fmt.Sprintf("マンボー: %v", err))
		log.Printf("Error scraping Manboo: %v", err)
	} else {
		allCafes = append(allCafes, manboo...)
		fmt.Printf("  → %d店舗を取得\n", len(manboo))
	}

	if len(errors) > 0 {
		fmt.Println("\n取得に失敗したサイト:")
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
	}

	return allCafes, nil
}