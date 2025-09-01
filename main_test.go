package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestNewNetCafeService(t *testing.T) {
	service := NewNetCafeService()
	
	if service == nil {
		t.Fatal("NewNetCafeService returned nil")
	}
	
	if service.client == nil {
		t.Error("client is nil")
	}
	
	if service.client.Timeout != 10*time.Second {
		t.Errorf("expected timeout 10s, got %v", service.client.Timeout)
	}
	
	if len(service.stores) == 0 {
		t.Error("stores is empty")
	}
}

func TestGetSampleStores(t *testing.T) {
	stores := getSampleStores()
	
	if len(stores) != 5 {
		t.Errorf("expected 5 stores, got %d", len(stores))
	}
	
	expectedFirstStore := NetCafe{
		Name:     "快活CLUB 新宿西口店",
		Location: "東京都新宿区西新宿1-12-9",
		Hours:    "24時間営業",
		Phone:    "03-5321-6166",
		URL:      "https://www.kaikatsu.jp/",
	}
	
	if !reflect.DeepEqual(stores[0], expectedFirstStore) {
		t.Errorf("first store does not match expected:\ngot: %+v\nexpected: %+v", stores[0], expectedFirstStore)
	}
}

func TestNetCafeService_SearchByName(t *testing.T) {
	service := NewNetCafeService()
	
	tests := []struct {
		name     string
		keyword  string
		expected int
	}{
		{"Search for 新宿", "新宿", 2},
		{"Search for 渋谷", "渋谷", 1},
		{"Search for 池袋", "池袋", 1},
		{"Search for 秋葉原", "秋葉原", 1},
		{"Search for 快活", "快活", 1},
		{"Search case insensitive", "DICE", 1},
		{"Search not found", "横浜", 0},
		{"Search empty", "", 5},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := service.SearchByName(tt.keyword)
			if len(results) != tt.expected {
				t.Errorf("expected %d results for keyword '%s', got %d", tt.expected, tt.keyword, len(results))
			}
		})
	}
}

func TestNetCafeService_GetAll(t *testing.T) {
	service := NewNetCafeService()
	stores := service.GetAll()
	
	if len(stores) != 5 {
		t.Errorf("expected 5 stores, got %d", len(stores))
	}
	
	if !reflect.DeepEqual(stores, service.stores) {
		t.Error("GetAll should return the same stores as service.stores")
	}
}

func TestNetCafeService_FetchFromAPI(t *testing.T) {
	// APIレスポンスをモックするテストサーバーを作成
	expectedCafe := NetCafe{
		Name:     "Test Cafe",
		Location: "Test Location",
		Hours:    "10:00-22:00",
		Phone:    "03-1234-5678",
		URL:      "https://test.com",
	}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedCafe)
	}))
	defer server.Close()
	
	service := NewNetCafeService()
	cafe, err := service.FetchFromAPI(server.URL)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !reflect.DeepEqual(*cafe, expectedCafe) {
		t.Errorf("fetched cafe does not match expected:\ngot: %+v\nexpected: %+v", *cafe, expectedCafe)
	}
}

func TestNetCafeService_FetchFromAPI_Error(t *testing.T) {
	service := NewNetCafeService()
	
	// 無効なURLでエラーをテスト
	_, err := service.FetchFromAPI("http://[invalid:url]:99999")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
	
	// 無効なJSONレスポンスのテスト
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	
	_, err = service.FetchFromAPI(server.URL)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestNetCafeService_SearchByLocation(t *testing.T) {
	service := NewNetCafeService()
	
	tests := []struct {
		name     string
		keyword  string
		expected int
	}{
		{"Search for 新宿区", "新宿区", 2},
		{"Search for 渋谷区", "渋谷区", 1},
		{"Search for 豊島区", "豊島区", 1},
		{"Search for 千代田区", "千代田区", 1},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := service.SearchByName(tt.keyword)
			if len(results) != tt.expected {
				t.Errorf("expected %d results for location '%s', got %d", tt.expected, tt.keyword, len(results))
			}
		})
	}
}