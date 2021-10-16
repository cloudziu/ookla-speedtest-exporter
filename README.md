# Ookla CLI Speedtest Exporter

Speedtest tool: https://www.speedtest.net/apps/cli wraped in Golang code to expose Internet connection metrics in Prometheus format:
- Download Speed
- Upload Speed
- Latency
- Jitter
- PacketLost

### Run
```bash
docker-compose up
```

### Change sccrape interval

Scrape interval is set in `./conf/prometheus/prometheus.yml` for 30 min.

```yaml
#./conf/prometheus/prometheus.yml
scrape_configs:
  - job_name: speedtest-ookla
    scrape_interval: 30m # <-- 
```
