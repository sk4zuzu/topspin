package srv1

//go:generate go generate github.com/sk4zuzu/topspin/micro/topspin/proto

import (
	"context"

	log "github.com/sirupsen/logrus"

	micro "github.com/micro/go-micro"

	"github.com/sk4zuzu/topspin/micro/topspin/proto/topspin"

	_ "github.com/micro/go-plugins/registry/nats"
	_ "github.com/micro/go-plugins/transport/nats"
)

const serviceName = "topspin.srv1"

type TopSpin struct{}

func (ts *TopSpin) Return(ctx context.Context, req *topspin.Ping, rsp *topspin.Pong) error {
	log.WithFields(log.Fields{
		"srv":            serviceName,
		"handler":        "TopSpin.Return",
		"req.GetMessage": req.GetMessage(),
	}).Info()

	rsp.Message = req.GetMessage() + ",pong?"

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

	topspin.RegisterTopSpinHandler(service.Server(), new(TopSpin))

	log.WithFields(log.Fields{
		"srv": serviceName,
	}).Error(service.Run())
}
