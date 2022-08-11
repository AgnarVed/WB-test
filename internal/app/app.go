package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
	stan "github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"wbTest/internal/cache"
	"wbTest/internal/config"
	"wbTest/internal/http"
	"wbTest/internal/repository"
	"wbTest/internal/repository/client"
	"wbTest/internal/server"
	"wbTest/internal/service"
)

func Run() {
	config, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	writers := make([]io.Writer, 0)
	writers = append(writers, os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(io.MultiWriter(writers...))

	// init psql database
	db, err := sql.Open(config.DriverName, config.DBConnStr)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()
	pgClient := client.NewPostgresClient(db)

	// init repository
	repos := repository.NewRepositories(&pgClient)

	// init services
	services := service.NewService(repos, config)

	// init cache
	//expTime, err := time.ParseDuration(config.CacheExpTime)
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	//cleanupTime, err := time.ParseDuration(config.CacheCleanupInterval)
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	//cache := cache.NewCache(expTime*time.Second, cleanupTime*time.Second)
	cache := cache.NewCache(5*time.Minute, 20*time.Minute)

	// init server, router, handlers
	srv := server.NewServer(config)
	http.NewHandlers(config, services).Init(srv.App())

	// start server
	go func() {
		err := srv.Run()
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	// nats streaming listen
	go func() {
		w := sync.WaitGroup{}
		w.Add(1)
		sc, err := stan.Connect(config.ClusterName, config.NatsClient, stan.NatsURL(config.NatsURL))
		if err != nil {
			logrus.Fatal(err)
		}
		defer sc.Close()
		_, err = sc.Subscribe(config.NatsSubject, func(msg *stan.Msg) {
			// TODO insert msg into cache
			order, err := Unmarshaler(msg.Data)
			if err != nil {
				logrus.Fatal("\nUnable to unmarshal message")
			}
			ok := cache.Exist(order.OrderUID)
			if ok != 0 {
				val, _ := cache.Get(order.OrderUID)
				fmt.Printf("\nOrder by id: %s is %s", order.OrderUID, val)
			} else {
				fmt.Printf("\nCannot find order in cache by ID: %s", order.OrderUID)
				logrus.Info("Inserting into cache")
				err = insertMsg(msg.Data, cache, 3)
				if err != nil {
					logrus.Fatal("Cannot Import Message")
				}
			}

			//cache.StartGC()
			//defer cache.Delete(string(indx))

		}, stan.DeliverAllAvailable(), stan.DurableName(config.DurableName))

		if err != nil {
			logrus.Fatal("Can't subscribe to channel", err)
		}

		w.Wait()
	}()
	cache.StartGC()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if err := srv.Stop(); err != nil {
		logrus.Fatal("Server forced to shut down", err)
	}
}

func insertMsg(data []byte, c *cache.Cache, timer time.Duration) error {
	o := NewOrder()
	err := json.Unmarshal(data, o)
	if err != nil {
		return err
	}
	c.Set(o.OrderUID, o, timer*time.Minute)
	fmt.Printf("Order with id: %s inserted correctly\n", o.OrderUID)
	return nil
}

func Unmarshaler(input []byte) (*Order, error) {
	o := NewOrder()
	err := json.Unmarshal(input, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

type Order struct {
	OrderUID          string   `json:"order_uid" faker:"len=20"`
	TrackNumber       string   `json:"track_number" faker:"len=20"`
	Entry             string   `json:"entry" faker:"oneof: WBIL"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale" faker:"oneof: en"`
	InternalSignature string   `json:"internal_signature" faker:"len=5"`
	CustomerID        string   `json:"customer_id" faker:"word"`
	DeliveryService   string   `json:"delivery_service" faker:"word"`
	Shardkey          string   `json:"shardkey" faker:"oneof: 9"`
	SmID              int64    `json:"sm_id" faker:"boundary_start=0, boundary_end=100"`
	DateCreated       string   `json:"date_created" faker:"date"`
	OofShard          string   `json:"oof_shard" faker:"oneof: 1"`
}

type Delivery struct {
	Name    string `json:"name" faker:"name"`
	Phone   string `json:"phone" faker:"e_164_phone_number"`
	Zip     string `json:"zip" faker:"oneof: 2639809"`
	City    string `json:"city" faker:"oneof: Kiryat Mozkin"`
	Address string `json:"address" faker:"oneof: Ploshad Mira 15"`
	Region  string `json:"region" faker:"oneof: Kraiot"`
	Email   string `json:"email" faker:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" faker:"len=20"`
	RequestID    string `json:"request_id" faker:"len=20"`
	Currency     string `json:"currency" faker:"currency"`
	Provider     string `json:"provider" faker:"oneof: wbpay"`
	Amount       int64  `json:"amount" faker:"boundary_start=100, boundary_end=10000"`
	PaymentDt    int64  `json:"payment_dt" faker:"unix_time"`
	Bank         string `json:"bank" faker:"word"`
	DeliveryCost int64  `json:"delivery_cost" faker:"boundary_start=100, boundary_end=10000"`
	GoodsTotal   int64  `json:"goods_total" faker:"boundary_start=1, boundary_end=100"`
	CustomFee    int64  `json:"custom_fee" faker:"boundary_start=0, boundary_end=10000"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id" faker:"boundary_start=100, boundary_end=10000"`
	TrackNumber string `json:"track_number" faker:"len=20"`
	Price       int64  `json:"price" faker:"boundary_start=100, boundary_end=10000"`
	Rid         string `json:"rid" faker:"len=20"`
	Name        string `json:"name" faker:"first_name"`
	Sale        int64  `json:"sale" faker:"boundary_start=0, boundary_end=100"`
	Size        string `json:"size" faker:"oneof: 0"`
	TotalPrice  int64  `json:"total_price" faker:"boundary_start=50, boundary_end=10000"`
	NmID        int64  `json:"nm_id" faker:"boundary_start=1000, boundary_end=1000000"`
	Brand       string `json:"brand" faker:"word"`
	Status      int64  `json:"status" faker:"boundary_start=0, boundary_end=500"`
}

func NewOrder() *Order {
	return &Order{
		OrderUID:          "",
		TrackNumber:       "",
		Entry:             "",
		Delivery:          Delivery{},
		Payment:           Payment{},
		Items:             nil,
		Locale:            "",
		InternalSignature: "",
		CustomerID:        "",
		DeliveryService:   "",
		Shardkey:          "",
		SmID:              0,
		DateCreated:       "",
		OofShard:          "",
	}
}
