package api

//go:generate go generate github.com/sk4zuzu/topspin/micro/topspin/proto

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	web "github.com/micro/go-web"

	"github.com/sk4zuzu/topspin/micro/topspin/proto/topspin"

	_ "github.com/micro/go-plugins/registry/nats"
	_ "github.com/micro/go-plugins/transport/nats"
)

const serviceName = "topspin.api"

type TopSpin struct{}

var (
	cl topspin.TopSpinClient
)

func (ts *TopSpin) Health(c *gin.Context) {
	log.WithFields(log.Fields{
		"srv":     serviceName,
		"handler": "TopSpin.Health",
	}).Info()

	c.String(200, "")
}

func (ts *TopSpin) Ping(c *gin.Context) {
	log.WithFields(log.Fields{
		"srv":     serviceName,
		"handler": "TopSpin.Ping",
	}).Info()

	c.String(200, "pong")
}

func (ts *TopSpin) Spin(c *gin.Context) {
	log.WithFields(log.Fields{
		"srv":     serviceName,
		"handler": "TopSpin.Spin",
	}).Info()

	rsp, err := cl.Return(context.TODO(), &topspin.Ping{
		Message: "ping?",
	})
	if err != nil {
		c.JSON(500, err)
		return
	}

	log.WithFields(log.Fields{
		"srv":            serviceName,
		"handler":        "TopSpin.Spin",
		"rsp.GetMessage": rsp.GetMessage(),
	}).Debug()

	c.JSON(200, rsp)
}

func (ts *TopSpin) Pods(c *gin.Context) {
	log.WithFields(log.Fields{
		"srv":     serviceName,
		"handler": "TopSpin.Pods",
	}).Info()

	pods, err := GetPods()
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, pods)
}

func NewAPIRouter() *gin.Engine {
	ts := new(TopSpin)
	router := gin.Default()
	router.GET("/ping", ts.Ping)
	router.GET("/spin", ts.Spin)
	router.GET("/pods", ts.Pods)
	router.GET("/", ts.Health)
	return router
}

func Run() {
	log.WithFields(log.Fields{
		"srv": serviceName,
	}).Info("Starting...")

	service := web.NewService(
		web.Name(serviceName),
		web.Address(":8080"),
	)

	service.Init()

	cl = topspin.NewTopSpinClient("topspin.srv2", client.DefaultClient)

	service.Handle("/", NewAPIRouter())

	log.WithFields(log.Fields{
		"srv": serviceName,
	}).Error(service.Run())
}
