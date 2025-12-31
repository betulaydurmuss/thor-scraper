# THOR Scraper 🛡️

THOR Scraper, **Tor ağı üzerinden anonim web taraması** yapabilen, hedef sitelerin **HTML içeriklerini ve ekran görüntülerini (screenshot)** otomatik olarak toplayan, **Go (Golang)** dili ile geliştirilmiş bir araçtır.

Bu proje özellikle **OSINT (Open Source Intelligence)** ve **CTI (Cyber Threat Intelligence)** toplama süreçleri için tasarlanmıştır.

---

## 📌 Proje Amacı

Siber tehdit aktörleri, altyapılarını ve sızıntı platformlarını gizlemek amacıyla sıklıkla **Tor ağı** kullanmaktadır.  
Yüzlerce `.onion` adresinin manuel olarak takip edilmesi hem zaman alıcı hem de verimsizdir.

Bu projenin amacı:
- Tor ağı üzerinden **anonim veri toplamak**
- Çoklu hedefleri **otomatik olarak taramak**
- HTML ve görsel delil (screenshot) toplamak
- Aktif ve pasif siteleri **log dosyaları ile raporlamak**
- CTI süreçlerinde **Collection** ve **Automation** yetkinliği kazandırmaktır

---

## 🚀 Özellikler

- ✅ Tor SOCKS5 proxy desteği (127.0.0.1:9150)
- ✅ YAML tabanlı hedef yönetimi
- ✅ Otomatik HTML içerik kaydı
- ✅ Headless Chrome ile ekran görüntüsü alma
- ✅ Zaman damgalı loglama sistemi
- ✅ IP sızıntısını önleyen özel HTTP client
- ✅ Hata toleranslı tarama (dead siteler programı durdurmaz)

---

## 🛠 Kullanılan Teknolojiler

- **Go (Golang)**
- **Tor Network**
- **net/http**
- **golang.org/x/net/proxy**
- **chromedp**
- **YAML**


