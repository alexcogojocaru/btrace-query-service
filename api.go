package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alexcogojocaru/query/config"
	storage "github.com/alexcogojocaru/query/proto-gen/btrace_storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type KeyValue struct {
	Type  string
	Value string
}

type Timestamp struct {
	Started  string
	Ended    string
	Duration float32
}

type Span struct {
	ID       string
	ParentID string
	Name     string
	Time     Timestamp
	Logs     []KeyValue
}

func NewStorageClient(host string, port int) storage.StorageClient {
	storageServerHost := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(storageServerHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot dial %s", storageServerHost)
	}

	return storage.NewStorageClient(conn)
}

func main() {
	conf, err := config.ParseConfig("config/config.yml")
	if err != nil {
		log.Fatal("Error while parsing the config file")
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	sc := NewStorageClient(conf.Storage.Hostname, int(conf.Storage.Port))

	log.Printf("Connecting to %s:%d", conf.Storage.Hostname, conf.Storage.Port)

	r.GET("/api/services", func(ctx *gin.Context) {
		var services []string
		stream, err := sc.GetServices(context.Background(), &storage.Empty{})
		if err != nil {
			log.Fatalf("GetServices err %v", err)
		}

		for {
			service, err := stream.Recv()
			if err == io.EOF {
				break
			}

			services = append(services, service.Name)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"services": services,
		})
	})

	r.GET("/api/service/:name/data", func(ctx *gin.Context) {
		var serviceData map[string][]Span = make(map[string][]Span)

		stream, err := sc.GetServiceData(context.Background(), &storage.ServiceName{
			Name: ctx.Param("name"),
		})
		if err != nil {
			log.Fatal(err)
		}

		for {
			var spans []Span
			data, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if data == nil {
				break
			}

			for _, span := range data.Spans {
				var logs []KeyValue

				for _, log := range span.Logs {
					logs = append(logs, KeyValue{
						Type:  log.Type,
						Value: log.Value,
					})
				}

				log.Printf("%s %v", span.SpanID, logs)
				spans = append(spans, Span{
					ID:       span.SpanID,
					ParentID: span.ParentSpanID,
					Name:     span.SpanName,
					Time: Timestamp{
						Started:  span.Timestamp.Started,
						Ended:    span.Timestamp.Ended,
						Duration: span.Timestamp.Duration,
					},
					Logs: logs,
				})
			}

			serviceData[data.TraceId] = spans
		}

		ctx.JSON(http.StatusOK, gin.H{
			"traces": serviceData,
		})
	})

	r.Run()
}
