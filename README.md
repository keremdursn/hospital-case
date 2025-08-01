# Hospital Case Backend

Bu proje, **Go Fiber** framework ile geliştirilmiş, PostgreSQL veritabanı kullanan ve JWT tabanlı kimlik doğrulama sistemi içeren, hastane yönetim sistemi backend uygulamasıdır.  
Projede **Clean Architecture** prensipleri kullanılmıştır. Ayrıca **Prometheus** ve **Grafana** ile monitoring entegrasyonu yapılmıştır. Docker ve Docker Compose desteği ile kolayca çalıştırılabilir.

---

## 🚀 Özellikler

### Temel Özellikler
- RESTful API geliştirme (Go Fiber)
- PostgreSQL entegrasyonu (GORM ORM)
- JWT tabanlı kimlik doğrulama (RSA imzalı tokenlar)
- Hastane, poliklinik, personel ve kullanıcı yönetimi
- Rol tabanlı yetkilendirme (yetkili, çalışan)
- Swagger/OpenAPI 3.0 API dokümantasyonu

### Güvenlik
- Rate limiting (endpoint bazında farklı kısıtlamalar)
- Input validation: TC kimlik numarası, telefon, e-posta formatı kontrolleri
- Parola güçlüğü kontrolü

### Monitoring ve Logging
- Prometheus ile detaylı metrik toplama (login başarısı/başarısızlığı, rate limit sayısı vb.)
- Grafana dashboard otomasyonu (dashboard import scripti)

### Containerization ve DevOps
- Docker ve Docker Compose desteği ile hızlı kurulum ve geliştirme ortamı
- Makefile ile yaygın komutların kolay kullanımı

---

## Başlangıç

### Gereksinimler
- Go 1.20+
- Docker & Docker Compose
- (Opsiyonel) golangci-lint

### Projeyi Çalıştırmak

```sh
# Docker ile tüm servisleri başlat
make docker-up

# Sadece Go uygulamasını başlatmak için
make run
```

### Migration
```sh
make migrate
```

### Test ve Lint
```sh
make test
make lint
```

## Monitoring

- Prometheus metrikleri `/metrics` endpointinde sunulur.
- Grafana ile Prometheus datasource ekleyip, örnek dashboardu `create_grafana_dashboard.sh` ile import edebilirsiniz.

## Swagger

- API dokümantasyonu `/swagger/index.html` adresinde erişilebilir.

🔧 Mimari Detaylar
- Projede Clean Architecture benimsenmiştir:
- Handler → UseCase → Repository katmanları şeklinde yapılandırılmıştır.

🗃️ Veritabanı ve Cache
- PostgreSQL (GORM ile)

- Redis (cache için)

## Ortam Değişkenleri & Konfigürasyon
- `configs/config.yml` dosyasını düzenleyerek ortam ayarlarını yapabilirsiniz.

## Kullanılan Teknolojiler
- Go (Fiber, GORM)
- PostgreSQL
- Redis
- Prometheus
- Grafana
- Docker / Docker Compose
- Swagger (OpenAPI 3.0)

## Lisans
MIT
