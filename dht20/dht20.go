package dht20

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/d2r2/go-i2c"
)

type DHT20 struct {
	client *i2c.I2C
}

func New() (*DHT20, error) {
	// Create new connection to I2C bus on 1 line with address 0x38
	c, err := i2c.NewI2C(0x38, 1)
	if err != nil {
		return nil, err
	}
	return &DHT20{
		client: c,
	}, nil
}

func (d *DHT20) Get() (float64, float64, error) {
	// Need to wait 100ms after launched
	time.Sleep(100 * time.Millisecond)
	var initial byte = 0x71
	ret, err := d.client.ReadRegU8(initial)
	if err != nil {
		return 0, -273, err
	}
	if ret != 0x1c {
		return 0, -273, errors.New("Initial code is not 28 (0x1c)")
	}

	time.Sleep(10 * time.Millisecond)
	// Start measure
	_, err = d.client.WriteBytes([]byte{0x00, 0xAC, 0x33, 0x00})
	if err != nil {
		return 0, -273, err
	}
	// Need to wait after sending ac3300
	time.Sleep(80 * time.Millisecond)
	dat := make([]byte, 7)
	_, err = d.client.ReadBytes(dat)
	if err != nil {
		return 0, -273, err
	}
	// byte is uint8 (8bits), it is not enough to shift
	// So cast []uint8 to []uint32
	var long_dat []uint32
	for _, d := range dat {
		long_dat = append(long_dat, uint32(d))
	}

	// Get humidity and tempreature data
	hum := long_dat[1]<<12 | long_dat[2]<<4 | ((long_dat[3] & 0xF0) >> 4)
	tmp := ((long_dat[3] & 0x0F) << 16) | long_dat[4]<<8 | long_dat[5]

	// Calcurate real data
	real_hum := float64(hum) / math.Pow(2, 20) * 100
	real_tmp := float64(tmp)/math.Pow(2, 20)*200 - 50

	return real_hum, real_tmp, nil
}

func (d *DHT20) Clean() {
	log.Println("cleaning i2c client...")
	d.client.Close()
}
