package hamsterbeat

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	gengrpc "hamsterbeat/gen/hamsterbeat.grpc"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	mu             sync.Mutex
	requestCounter int
	failedRequest  int
)

type ServerStruct struct {
	gengrpc.UnimplementedHamsterbeatServer
	worker *Worker
}

func (s *ServerStruct) Pulse(ctx context.Context, in *gengrpc.HamsterbeatRequest) (*gengrpc.HamsterbeatResponse, error) {
	if in.Animaltypeid > int64(len(Zoopark)) {
		return &gengrpc.HamsterbeatResponse{Result: false}, nil
	}
	fmt.Printf("Received %s (%d):%d\n", Zoopark[in.Animaltypeid][0], in.Animalid, in.Heartbeat)

	s.worker.Publish(in)

	return &gengrpc.HamsterbeatResponse{Result: true}, nil
}

func Connect() {

	conn, _ := grpc.NewClient(GRPC_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))

	for animalTypeId, val := range Zoopark {
		var wg sync.WaitGroup
		var typeLimit int64
		var err error
		typeLimit, err = strconv.ParseInt(val[1], 10, 64)
		if err != nil {
			fmt.Printf("Incorrect value in Zoopark: %s", err)
			continue
		}
		var redis RedisCon = RedisCon{}
		for animalNumber := int64(1); animalNumber <= typeLimit; animalNumber++ {

			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() {
					if p := recover(); p != nil {
						fmt.Printf("Panic recovery %v\n", p)
					}
				}()

				heartbeat, _ := MakeNewHeartbeat(animalTypeId, animalNumber, &redis)
				proccessAnimal(animalTypeId, animalNumber, heartbeat, conn)
			}()
		}
		wg.Wait()
		fmt.Printf("Request %s: %d (%d)\n", Zoopark[animalTypeId][0], requestCounter, failedRequest)
	}
}

func MakeNewHeartbeat(animalTypeId int64, animalNumber int64, redis *RedisCon) (int64, error) {
	if animalTypeId < 1 {
		return int64(50), errors.New("animalTypeId < 0")
	} else if animalTypeId > int64(len(Zoopark)) {
		return int64(50), errors.New("animalTypeId limit")
	}
	/*
		Хак. Обращаемся во внутренний контур для предыдущего значения
		Это позволяет рисовать сглаженные графики
	*/
	var heartbeat = redis.Get(animalTypeId, animalNumber)
	heartbeat += int64(rand.Int31n(3) - 1)
	/*
		Хомячки не умирают
		Высокий пульс принудительно опускаем ниже/выше лимита
	*/
	if heartbeat < 1 {
		heartbeat = 1
	} else if heartbeat > 100 {
		heartbeat = 100
	}
	return heartbeat, nil
}

func proccessAnimal(animalTypeId int64, animalNumber int64, heartbeat int64, conn *grpc.ClientConn) {
	var sendedRequest = 0
	var failRequest = 0
	client := gengrpc.NewHamsterbeatClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	r, err := client.Pulse(ctx, &gengrpc.HamsterbeatRequest{Animaltypeid: animalTypeId, Animalid: animalNumber, Heartbeat: heartbeat})
	if err != nil {
		failedRequest++
	} else {
		sendedRequest++
		if !r.Result {
			fmt.Printf("Failed\n")
		}
	}
	mu.Lock()
	requestCounter += sendedRequest
	failedRequest += failRequest
	mu.Unlock()
}

type Worker struct {
	Channel *chan *gengrpc.HamsterbeatRequest
	Redis   *RedisCon
}

func (w *Worker) Publish(in *gengrpc.HamsterbeatRequest) {
	*w.Channel <- in
}

func (w *Worker) Reader() {
	for {
		x, ok := <-*w.Channel
		if !ok {
			break
		}

		var str string = protojson.Format(x)
		w.Redis.Set(x.GetAnimaltypeid(), x.GetAnimalid(), str)
	}
}

type GrpcListener interface {
	GetTCPListner() net.Listener
	GetRPCListner() *grpc.Server
}

type NetStruct struct {
	listner net.Listener
	s       *grpc.Server
}

func (l *NetStruct) GetTCPListner() net.Listener {
	return l.listner
}

func (l *NetStruct) GetRPCListner() *grpc.Server {
	return l.s
}

func Server() error {
	listner, err := net.Listen("tcp4", GRPC_ADDR)
	if err != nil {
		return err
	}
	g := &NetStruct{listner: listner, s: grpc.NewServer()}
	return StartServer(g)
}

func NewWorker() *Worker {
	ch := make(chan *gengrpc.HamsterbeatRequest, 1024)
	return &Worker{Channel: &ch, Redis: &RedisCon{}}
}

func StartServer(g GrpcListener) error {
	worker := NewWorker()
	defer close(*worker.Channel)
	go worker.Reader()
	gengrpc.RegisterHamsterbeatServer(g.GetRPCListner(), &ServerStruct{worker: worker})
	fmt.Printf("Server listening at %s\n", g.GetTCPListner().Addr())
	if err := g.GetRPCListner().Serve(g.GetTCPListner()); err != nil {
		return err
	}
	return nil
}
