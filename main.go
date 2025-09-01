package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type NetCafe struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Hours    string `json:"hours"`
	Phone    string `json:"phone"`
	URL      string `json:"url"`
}

type NetCafeService struct {
	client *http.Client
	stores []NetCafe
}

func NewNetCafeService() *NetCafeService {
	return &NetCafeService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		stores: getSampleStores(),
	}
}

func getSampleStores() []NetCafe {
	return []NetCafe{
		{
			Name:     "快活CLUB 新宿西口店",
			Location: "東京都新宿区西新宿1-12-9",
			Hours:    "24時間営業",
			Phone:    "03-5321-6166",
			URL:      "https://www.kaikatsu.jp/",
		},
		{
			Name:     "自遊空間 池袋西口ROSA店",
			Location: "東京都豊島区西池袋1-37-12",
			Hours:    "24時間営業",
			Phone:    "03-5391-7778",
			URL:      "https://jiqoo.jp/",
		},
		{
			Name:     "DiCE 秋葉原店",
			Location: "東京都千代田区外神田1-11-5",
			Hours:    "24時間営業",
			Phone:    "03-5298-1281",
			URL:      "https://www.diskcity.co.jp/",
		},
		{
			Name:     "マンボー 渋谷宮益坂店",
			Location: "東京都渋谷区渋谷1-12-1",
			Hours:    "24時間営業",
			Phone:    "03-5766-6010",
			URL:      "https://manboo.co.jp/",
		},
		{
			Name:     "アプレシオ 新宿歌舞伎町店",
			Location: "東京都新宿区歌舞伎町1-20-1",
			Hours:    "24時間営業",
			Phone:    "03-5155-4486",
			URL:      "https://www.aprecio.co.jp/",
		},
	}
}

func (s *NetCafeService) SearchByName(keyword string) []NetCafe {
	var results []NetCafe
	keyword = strings.ToLower(keyword)
	
	for _, cafe := range s.stores {
		if strings.Contains(strings.ToLower(cafe.Name), keyword) ||
			strings.Contains(strings.ToLower(cafe.Location), keyword) {
			results = append(results, cafe)
		}
	}
	return results
}

func (s *NetCafeService) GetAll() []NetCafe {
	return s.stores
}

func (s *NetCafeService) FetchFromAPI(url string) (*NetCafe, error) {
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var cafe NetCafe
	if err := json.Unmarshal(body, &cafe); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &cafe, nil
}

func printCafe(cafe NetCafe) {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("店舗名: %s\n", cafe.Name)
	fmt.Printf("場所:   %s\n", cafe.Location)
	fmt.Printf("営業時間: %s\n", cafe.Hours)
	fmt.Printf("電話番号: %s\n", cafe.Phone)
	fmt.Printf("URL:    %s\n", cafe.URL)
}

func main() {
	var (
		scrapeFlag = flag.Bool("scrape", false, "Webサイトから最新の店舗情報を取得")
		helpFlag   = flag.Bool("help", false, "ヘルプを表示")
	)
	flag.Parse()

	if *helpFlag {
		fmt.Println("ネットカフェ営業時間取得ツール")
		fmt.Println("\n使い方:")
		fmt.Println("  ./netcafe [オプション] [検索キーワード]")
		fmt.Println("\nオプション:")
		fmt.Println("  -scrape    Webサイトから最新の店舗情報を取得")
		fmt.Println("  -help      このヘルプを表示")
		fmt.Println("\n例:")
		fmt.Println("  ./netcafe                    # 登録済み店舗一覧を表示")
		fmt.Println("  ./netcafe 新宿               # 「新宿」で店舗を検索")
		fmt.Println("  ./netcafe -scrape            # Webから最新情報を取得")
		fmt.Println("  ./netcafe -scrape 渋谷       # 最新情報から「渋谷」で検索")
		return
	}

	var stores []NetCafe

	if *scrapeFlag {
		fmt.Println("Webサイトから最新の店舗情報を取得しています...")
		fmt.Println(strings.Repeat("-", 50))
		
		scraper := NewScraper()
		scrapedStores, err := scraper.ScrapeAll()
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			fmt.Println("サンプルデータを使用します。")
			stores = getSampleStores()
		} else {
			stores = scrapedStores
			fmt.Printf("\n合計 %d 店舗の情報を取得しました。\n", len(stores))
		}
	} else {
		stores = getSampleStores()
	}

	service := &NetCafeService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		stores: stores,
	}

	args := flag.Args()
	if len(args) > 0 {
		keyword := strings.Join(args, " ")
		fmt.Printf("\n「%s」で検索中...\n\n", keyword)
		
		results := service.SearchByName(keyword)
		if len(results) == 0 {
			fmt.Println("該当する店舗が見つかりませんでした。")
			return
		}
		
		fmt.Printf("%d件の店舗が見つかりました:\n", len(results))
		for _, cafe := range results {
			printCafe(cafe)
		}
	} else {
		if *scrapeFlag {
			fmt.Println("\n取得した店舗一覧:")
		} else {
			fmt.Println("ネットカフェ営業時間情報")
			fmt.Println("使い方: ./netcafe -help でヘルプを表示")
			fmt.Println("\n登録済み店舗一覧:")
		}
		
		for _, cafe := range service.GetAll() {
			printCafe(cafe)
		}
	}
}