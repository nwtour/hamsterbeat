package hamsterbeat

import (
	"context"
	//	"errors"
	//	"fmt"
	//"github.com/golang/mock/gomock"
	gengrpc "hamsterbeat/gen/hamsterbeat.grpc"
	//	"hamsterbeat/internal/hamsterbeat"
	//"github.com/stretchr/testify/mock"
	"testing"
)

func TestPulse(t *testing.T) {
	s := &ServerStruct{worker: NewWorker()}
	var ctx = context.Background()
	req := gengrpc.HamsterbeatRequest{Animaltypeid: int64(len(Zoopark) + 1)}
	res, _ := s.Pulse(ctx, &req)
	if res.Result {
		t.Error("Тип животного за пределами доступного не вернул ошибку")
	}
	req = gengrpc.HamsterbeatRequest{Animaltypeid: 1, Animalid: 99999999}
	res, _ = s.Pulse(ctx, &req)
	if !res.Result {
		t.Errorf("Корректное животное не отметилось в системе %v", res)
	}
	close(*s.worker.Channel)
}

func TestMakeNewHearbeat(t *testing.T) {
	r := &RedisCon{}
	h, err := MakeNewHeartbeat(0, 0, r)
	if err == nil {
		t.Error("Тип животного за пределами доступного не вернул ошибку")
	}
	h, err = MakeNewHeartbeat(int64(len(Zoopark)+1), 1, r)
	if err == nil {
		t.Error("Тип животного за пределами доступного не вернул ошибку")
	}
	h, err = MakeNewHeartbeat(1, 1, r)
	if h > 100 {
		t.Error("Пульс животного слишком велик")
	}
}

/*
type T2 struct {
	mock.Mock
}

func (d *T2) GetRPCListner() {
	return
}
*/

/*
func TestServer(t *testing.T) {
	m := NewMockGrpcListener(t)
	m.listner = nil
	m.On("GetRPCListner").Return(nil)
	err := StartServer(m)
	if err == nil {
		t.Error("При ошибке биндинга сервер не вернул ошибку")
	}
}
*/
