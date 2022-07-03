package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	storage "github.com/alexcogojocaru/query/proto-gen/btrace_storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type KeyValue struct {
	Type  string
	Value string
}

type Span struct {
	ID       string
	ParentID string
	Name     string
	Logs     []KeyValue
}

func NewStorageClient() storage.StorageClient {
	storageServerHost := "localhost:50051"
	conn, err := grpc.Dial(storageServerHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot dial %s", storageServerHost)
	}

	return storage.NewStorageClient(conn)
}

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	sc := NewStorageClient()

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

			log.Print(data)
			for _, span := range data.Spans {
				var logs []KeyValue

				for _, log := range span.Logs {
					fmt.Print(log)
					logs = append(logs, KeyValue{
						Type:  log.Type,
						Value: log.Value,
					})
				}

				spans = append(spans, Span{
					ID:       span.SpanID,
					ParentID: span.ParentSpanID,
					Name:     span.SpanName,
					Logs:     logs,
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
