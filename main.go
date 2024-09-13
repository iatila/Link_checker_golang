package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// URL için yapı (formda gelecek)
type LinkCheckRequest struct {
	Link string `json:"link"`
}

// Yanıt için yapı
type LinkCheckResponse struct {
	Link        string `json:"link"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// URL kontrolü yapar ve HTTP kodunu basar
func checkLink(w http.ResponseWriter, r *http.Request) {
	var req LinkCheckRequest

	// Gelen yanıt JSON mu
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	// URL alanı boşsa
	if req.Link == "" {
		resp := LinkCheckResponse{Link: req.Link, Status: 0, Description: "URL boş"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Protokol yoksa ekler
	if !strings.HasPrefix(req.Link, "http://") && !strings.HasPrefix(req.Link, "https://") {
		req.Link = "https://" + req.Link
	}

	// URL Kontrolü
	parsedURL, err := url.Parse(req.Link)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		resp := LinkCheckResponse{Link: req.Link, Status: 0, Description: "Geçersiz URL formatı"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Domain Regex
	domainRegex := `^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	domain := parsedURL.Hostname()

	// Domain geçerli mi?
	matched, err := regexp.MatchString(domainRegex, domain)
	if err != nil || !matched {
		resp := LinkCheckResponse{Link: req.Link, Status: 0, Description: "Geçersiz domain"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Yönlendirmeleri  devre dışı bırak 10 saniye timeout eklendi
	client := http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(req.Link)
	// Eğer istek başarısız olursa
	if err != nil {
		resp := LinkCheckResponse{Link: req.Link, Status: 0, Description: fmt.Sprintf("Bağlantı hatası: %v", err)}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	// Yanıtları sonlandırma
	defer resp.Body.Close()

	// HTTP Kodları
	var description string
	var status int
	switch resp.StatusCode {
	case 200:
		description = "Başarılı"
		status = 200
	case 301:
		description = "Kalıcı Yönlendirme"
		status = 301
	case 302:
		description = "Geçici Yönlendirme"
		status = 302
	case 403:
		description = "Erişim Engellenmiş"
		status = 403
	case 404:
		description = "Bulunamadı"
		status = 404
	case 500:
		description = "Sunucu Hatası"
		status = 500
	case 502:
		description = "Kötü Ağ Yönlendirmesi"
		status = 502
	case 503:
		description = "Hizmet Kullanılamıyor"
		status = 503
	default:
		description = fmt.Sprintf("Durum Kodu: %d", resp.StatusCode)
		status = resp.StatusCode
	}

	// Yanıtı JSON formatında döndürür
	response := LinkCheckResponse{
		Link:        req.Link,
		Status:      status,
		Description: description,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Şablon dosyasını çek ve bas.
func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Şablon çözümleme hatası", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "Şablon çalıştırma hatası", http.StatusInternalServerError)
	}
}

// Ana sayfayı şablon dosyasını kullanarak render eder.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/index.html")
}

func main() {
	// Ana sayfa handler
	http.HandleFunc("/", indexHandler)
	// Link kontrol handler
	http.HandleFunc("/check", checkLink)
	// Statik dosyalar için handler
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Sunucu http://localhost:8080 adresinde çalışıyor")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
