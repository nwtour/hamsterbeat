package hamsterbeat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
)

type MyCollector struct {
	counterDesc  *prometheus.Desc
	Counter      string
	AnimalTypeId int64
	AnimalNumber int64
	Redis        *RedisCon
}

func (c *MyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.counterDesc
}

func (c *MyCollector) Collect(ch chan<- prometheus.Metric) {
	value := c.Redis.Get(c.AnimalTypeId, c.AnimalNumber)
	fmt.Printf("Start Collect %s %d\n", c.Counter, value)

	ch <- prometheus.MustNewConstMetric(
		c.counterDesc,
		prometheus.CounterValue,
		float64(value),
	)
}

func NewMyCollector(counter string, animalTypeId int64, animalNumber int64, redis *RedisCon) *MyCollector {
	return &MyCollector{
		counterDesc:  prometheus.NewDesc(counter, "Help string", nil, nil),
		Counter:      counter,
		AnimalTypeId: animalTypeId,
		AnimalNumber: animalNumber,
		Redis:        redis,
	}
}

type RedisCon struct {
	redis     *redis.Client
	connected bool
}

func (c *RedisCon) Get(typeId int64, animalId int64) int64 {
	c.Connect()
	ctx := context.Background()
	val, err := c.redis.Get(ctx, fmt.Sprintf(REDIS_KEY_TEMPLATE, typeId, animalId)).Result()
	if err != nil && err != redis.Nil {
		fmt.Printf("Redis error %s\n", err)
		return 50
	}
	var data struct{ Heartbeat string }
	json.Unmarshal([]byte(val), &data)
	r, _ := strconv.Atoi(data.Heartbeat)
	return int64(r)
}

func (c *RedisCon) Set(animalTypeId int64, animalId int64, protobuf string) error {
	c.Connect()
	ctx := context.Background()
	return c.redis.Set(ctx, fmt.Sprintf(REDIS_KEY_TEMPLATE, animalTypeId, animalId), protobuf, 0).Err()
}

func (c *RedisCon) Connect() {
	if !c.connected {
		c.redis = redis.NewClient(&redis.Options{Addr: REDIS_ADDR})
		c.connected = true
	}
}

func MakePrometeusMetric() (ret int) {
	for animalTypeId := range Zoopark {
		typeLimit, _ := strconv.ParseInt(Zoopark[animalTypeId][1], 10, 64)
		r := &RedisCon{}
		r.Connect()
		for animalNumber := int64(1); animalNumber <= typeLimit; animalNumber++ {
			ret++
			prometheus.MustRegister(
				NewMyCollector(
					fmt.Sprintf(REDIS_KEY_TEMPLATE, animalTypeId, animalNumber),
					animalTypeId,
					animalNumber,
					r,
				),
			)
		}
	}
	return ret
}

func ServeMetrics() {
	MakePrometeusMetric()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(PROMETEUS_ADDR, nil)
}
