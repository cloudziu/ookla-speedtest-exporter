package speedtest_ookla

import (
	"encoding/json"
	"log"
	"os/exec"
)

const (
	binaryName = "speedtest"
)

var (
	binaryArgs = []string{"-f", "json"}
)

type Client struct {
	Download     float64
	Upload       float64
	Latency      float64
	Jitter       float64
	PacketLoss   float64
}

type SpeedtestResult struct {
	Ping struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
	} `json:"ping"`
	Download struct {
		Bandwidth int `json:"bandwidth"`
	} `json:"download"`
	Upload struct {
		Bandwidth int `json:"bandwidth"`
	} `json:"upload"`
	PacketLoss int `json:"packetLoss"`
}

func (c *Client) RunSpeedtest() {

	log.Printf("Run speedtest...")
	out, err := exec.Command(binaryName, binaryArgs...).Output()

	if err != nil {
		log.Println(err)
	}

	r := SpeedtestResult{}

	if err = json.Unmarshal(out, &r); err != nil {
		log.Println("Error executing speedtest, set metrics to 0")
		c.setMetricsToZero()
	} else {

		c.setMetrics(
			bytesToMegabits(r.Download.Bandwidth),
			bytesToMegabits(r.Upload.Bandwidth),
			r.Ping.Latency,
			r.Ping.Jitter,
			float64(r.PacketLoss),
		)
		log.Printf("Speedtest complete => download=%.0fMbp/s "+
			"upload=%.0fMbp/s "+
			"latency=%.2fms "+
			"jitter=%.2fms "+
			"packet_loss=%.1f%%",
			c.Download,
			c.Upload,
			c.Latency,
			c.Jitter,
			c.PacketLoss)
	}
}

func (c *Client) setMetrics(d float64, u float64, l float64, j float64, p float64) {
	c.Download = d
	c.Upload = u
	c.Latency = l
	c.Jitter = j
	c.PacketLoss = p
}

func (c *Client) setMetricsToZero() {
	c.Download = 0
	c.Upload = 0
	c.Latency = 0
	c.Jitter = 0
	c.PacketLoss = 0
}

// Convert Bytes to Megabits
func bytesToMegabits(bytes int) float64 {
	result := ((bytes * 8) / 1024) / 1024
	return float64(result)
}

// check if speedtest is installed
func InitialCheck() {
	// check if speedtest is installed
	log.Println("Run dependency check")
	if _, err := exec.LookPath(binaryName); err != nil {
		log.Fatal(err)
	}

	if _, err := exec.Command(binaryName, "--version").Output(); err != nil {
		log.Fatal(err)
	}
}
