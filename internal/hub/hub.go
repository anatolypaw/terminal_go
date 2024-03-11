package hub

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"terminal/internal/entity"
	pb "terminal/internal/hub/grpcapi"
)

type Hub struct {
	HubAddres string
	tname     string
	client    pb.HubClient

	Connected bool // Флаг, что связь с сервером есть
}

func New(hubAddr string, tname string) Hub {
	// Set up a connection to the server.
	conn, err := grpc.Dial(hubAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("gRPC did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewHubClient(conn)

	return Hub{
		HubAddres: hubAddr,
		client:    client,
		tname:     tname,
	}
}

// Выполняет слежебные операции
func (h *Hub) Run() {

	for {

		time.Sleep(1000 * time.Millisecond)
	}
}

// Возвращает код для печати
func (h *Hub) GetCodeForPrint(gtin string) (entity.Code, error) {
	timeout := 50 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request := pb.GetCodeForPrintRequest{
		Tname: h.tname,
		Gtin:  gtin,
	}

	resp, err := h.client.GetCodeForPrint(ctx, &request)

	return entity.Code{
			Gtin:   resp.GetGtin(),
			Serial: resp.GetSerial(),
			Crypto: resp.GetCrypto(),
		},
		err
}
