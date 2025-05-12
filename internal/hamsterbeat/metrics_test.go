package hamsterbeat

import (
	"testing"
)

func TestMyNewCollecor(t *testing.T) {
	mnc := NewMyCollector("test", 1, 1, &RedisCon{})
	if mnc.AnimalTypeId != 1 {
		t.Error("Неудачно создан коллектор меток")
	}
}

func TestMakePrometeusMetric(t *testing.T) {
	if MakePrometeusMetric() < 1 {
		t.Error("Не создались метрики в Prometeus")
	}
}

func TestRedisSet(t *testing.T) {
	r := &RedisCon{}
	e := r.Set(99999, 99999, "{}")
	if e != nil {
		t.Error("Не записалось в Redis")
	}
}
