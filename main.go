package main 

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"context"
	"strings"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v3"
)

const(
	TorProxyAddress = "127.0.0.1:9150"
	LogFilePath = "logs/scan_report.log"
	OutputDir = "output"
)

type Site struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	Targets []Site `yaml:"targets"`
}

func writeLog(level string, message string) {
if _, err := os.Stat("logs"); os.IsNotExist(err) {
os.Mkdir("logs", 0755)
}	
file, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
fmt.Printf("[SİSTEM HATASI] log dosyası açılamadı : %v\n", err)
return
}
defer file.Close()

writer := bufio.NewWriter(file)
timestamp := time.Now().Format("2006-01-02 15:04:05")

LogLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

writer.WriteString(LogLine)
writer.Flush()
}

func getTorHttpClient() *http.Client {
dialer, err := proxy.SOCKS5("tcp", TorProxyAddress, nil, proxy.Direct)
if err != nil {
log.Fatalf("Tor bağlantı hatası : %v", err)
}
tr := &http.Transport {
Dial : dialer.Dial,
}
return &http.Client {
Transport : tr,
Timeout : time.Second * 45,
}
}

func saveHTML(siteName string, body []byte){
filename := fmt.Sprintf("%s_%d.html", siteName, time.Now().Unix())
path := filepath.Join(OutputDir, filename)

err := os.WriteFile(path, body, 0644)
if err != nil {
writeLog("ERROR", fmt.Sprintf("%s HTML kaydedilemedi : %v", siteName, err))
} else {
msg := fmt.Sprintf("HTML verisi kaydedildi : %s", filename)
fmt.Println("[HTML]" + msg)
writeLog("SUCCESS", msg+" ("+siteName+")")
}
}

func takeScreenshot(site Site) {
if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
os.Mkdir(OutputDir, 0755)
}

opts := append(chromedp.DefaultExecAllocatorOptions[:],
chromedp.ProxyServer("socks5://"+TorProxyAddress),
chromedp.WindowSize(1920, 1080),
chromedp.Flag("headless", true), 
chromedp.DisableGPU,
)

ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
defer cancel()
ctx, cancel = chromedp.NewContext(ctx)
defer cancel()
ctx, cancel = context.WithTimeout(ctx, 90*time.Second)
defer cancel()

var buf []byte
fmt.Printf("[GÖRSEL] %s için screenshot alınıyor...\n", site.Name)
writeLog("INFO", fmt.Sprintf("Screenshot işlemi başlatıldı : %s", site.Name))

err := chromedp.Run(ctx,
	chromedp.Navigate(site.URL),
	chromedp.Sleep(5*time.Second), 
	chromedp.CaptureScreenshot(&buf),
)

if err != nil {
msg := fmt.Sprintf("Screenshot başarısız (%s) : %v", site.Name, err)
fmt.Println("[HATA] " + msg)
writeLog("ERROR", msg)
return
}

filename := fmt.Sprintf("%s_%d.png", site.Name, time.Now().Unix())
path := filepath.Join(OutputDir, filename)
if err := os.WriteFile(path, buf, 0644); err != nil {
writeLog("ERROR", fmt.Sprintf("Resim dosyası yazılamadı : %v", err))
} else {
msg := fmt.Sprintf("Screenshot kaydedildi : %s", filename)
fmt.Println("[FOTO]" + msg)
writeLog("SUCCESS", msg)
}
}

func main() {
	
fmt.Print("\033[H\033[2J")
fmt.Println("=== THOR SCRAPER ===")
	
if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
os.Mkdir(OutputDir, 0755)
}

writeLog("INFO", "=== YENİ TARAMA OTURUMU BAŞLATILDI ===")

targetFile := "targets.yaml"
if len(os.Args) > 1 {
targetFile = os.Args[1]
}

data, err := os.ReadFile(targetFile)
if err != nil {
log.Fatalf("[KRİTİK] Hedef dosyası okunamadı : %v", err)
}

var cfg Config
if err := yaml.Unmarshal(data, &cfg); err != nil {
log.Fatalf("[KRİTİK] YAML format hatası : %v", err)
}

httpClient := getTorHttpClient()
	
fmt.Println("Tor bağlantısı ve IP kontrol ediliyor...")
checkIPResp, err := httpClient.Get("http://check.torproject.org/")
if err == nil {
fmt.Println("Bağlantı Başarılı. Tor Ağı üzerindeyiz.")
checkIPResp.Body.Close()
} else {
fmt.Println("[UYARI] Tor ağına bağlanılamadı, lütfen Tor Browser'ı kontrol edin!")
}
fmt.Printf("\n[BİLGİ] %d adet hedef yüklendi. Tarama Başlıyor...\n\n", len(cfg.Targets))

for i, site := range cfg.Targets {
site.URL = strings.TrimSpace(site.URL)

fmt.Printf("[%d/%d] Taranıyor : %s\n", i+1, len(cfg.Targets), site.Name)
writeLog("INFO", fmt.Sprintf("Scanning : %s -> %s", site.Name, site.URL))

resp, err := httpClient.Get(site.URL)
		
if err != nil {
msg := fmt.Sprintf("Erişim sağlanamadı : %v", err)
fmt.Println("[HATA]" + msg)
writeLog("FAIL", msg+" ("+site.Name+")")
continue
}

if resp.StatusCode == 200 {
body, _ := io.ReadAll(resp.Body)
saveHTML(site.Name, body)
resp.Body.Close() 
takeScreenshot(site)
} else {
msg := fmt.Sprintf("Site aktif ama hata kodu döndü : %d", resp.StatusCode)
fmt.Printf("[UYARI] %s\n", msg)
writeLog("WARN", msg+" ("+site.Name+")")
}

time.Sleep(2 * time.Second)
}

fmt.Println("\n=== TARAMA TAMAMLANDI ===")
fmt.Printf("Rapor : %s\n", LogFilePath)
fmt.Printf("Veriler : %s klasöründe.\n", OutputDir)
writeLog("INFO", "=== OTURUM SONLANDI ===")
}

