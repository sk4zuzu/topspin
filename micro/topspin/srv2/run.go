package srv2

//go:generate go generate github.com/sk4zuzu/topspin/micro/topspin/proto

import (
	"context"

	log "github.com/sirupsen/logrus"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"

	"github.com/sk4zuzu/topspin/micro/topspin/proto/topspin"

	_ "github.com/micro/go-plugins/registry/nats"
	_ "github.com/micro/go-plugins/transport/nats"
)

const serviceName = "topspin.srv2"

type TopSpin struct{}

var (
	cl topspin.TopSpinClient
)

func (ts *TopSpin) Return(ctx context.Context, req *topspin.Ping, rsp *topspin.Pong) error {
	log.WithFields(log.Fields{
		"srv":            serviceName,
		"handler":        "TopSpin.Return",
		"req.GetMessage": req.GetMessage(),
	}).Info()

	rsp2, err := cl.Return(context.TODO(), &topspin.Ping{
		Message: req.GetMessage(),
	})

	if err != nil {
		return err
	}

	rsp.Message = rsp2.GetMessage() + ",pong!"

	log.WithFields(log.Fields{
		"srv":            serviceName,
		"handler":        "TopSpin.Return",
		"rsp.GetMessage": rsp.GetMessage(),
	}).Debug()

	return nil
}

func Run() {
	log.WithFields(log.Fields{
		"srv": serviceName,
	}).Info("Starting...")

	service := micro.NewService(
		micro.Name(serviceName),
	)

	service.Init()

	cl = topspin.NewTopSpinClient("topspin.srv1", client.DefaultClient)

	topspin.RegisterTopSpinHandler(service.Server(), new(TopSpin))

	log.WithFields(log.Fields{
		"srv": serviceName,
	}).Error(service.Run())
}
