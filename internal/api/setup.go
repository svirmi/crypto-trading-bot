package api

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

const (
	INCLUSIVE_MIN_PORT = 1024
	EXCLUSIVE_MAX_PORT = 65354

	API_V1             = "/api/v1"
	RESOURCE_NOT_FOUND = "resource not found"
)

var (
	exchange model.IExchange
)

func Initialize(_exchange model.IExchange) *gin.Engine {
	exchange = _exchange

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(log_request, gin.Recovery())

	// Ping pong
	r.GET(API_V1+"/ping", func(c *gin.Context) {
		status, body, apiErr := http.StatusOK, ping_pong_dto{"pong"}, error_dto{}
		send(c, status, body, apiErr)
	})

	// Start execution
	r.POST(API_V1+"/executions", func(c *gin.Context) {
		var req exe_create_req_dto
		err := c.ShouldBindJSON(&req)
		if err != nil {
			send_bad_request(c, err)
			return
		}

		err = req.Validate()
		if err != nil {
			send_bad_request(c, err)
			return
		}

		status, body, apiErr := create_execution(req)
		send(c, status, body, apiErr)
	})

	// Terminate execution
	r.PUT(API_V1+"/executions/:exeId", func(c *gin.Context) {
		exeId := c.Param("exeId")
		var req exe_update_req_dto
		err := c.ShouldBindJSON(&req)
		if err != nil {
			send_bad_request(c, err)
			return
		}

		err = req.Validate()
		if err != nil {
			send_bad_request(c, err)
			return
		}

		status, body, apiErr := update_execution(exeId, req)
		send(c, status, body, apiErr)
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, error_dto{RESOURCE_NOT_FOUND})
	})

	return start_server(r)
}

func start_server(r *gin.Engine) *gin.Engine {
	// Computing port number
	port := config.GetServerConfig().Port
	if port < 0 {
		logrus.Panicf(logger.API_ERR_NEGATIVE_PORT_NUMBER, port)
	} else if port == 0 {
		port = rand.Intn(EXCLUSIVE_MAX_PORT-INCLUSIVE_MIN_PORT) + INCLUSIVE_MIN_PORT
	} else if port < INCLUSIVE_MIN_PORT {
		logrus.Panicf(logger.API_ERR_RESERVED_PORT_NUMBER, port)
	} else if port >= EXCLUSIVE_MAX_PORT {
		logrus.Panicf(logger.API_ERR_PORT_NUMBER_OUT_OF_RANGE, port)
	}

	// Getting public IP address
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logrus.Panicf(logger.API_ERR_DIAL_UP, "udp://8.8.8.8:80")
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// Starting up server
	logrus.Infof(logger.API_SERVER_STARTUP, localAddr.IP.String(), port)
	r.Run(fmt.Sprintf(":%d", port))
	return r
}

func send_bad_request(c *gin.Context, err error) {
	send(c, http.StatusBadRequest, nil, error_dto{err.Error()})
}

func send(c *gin.Context, status int, body interface{}, errDto error_dto) {
	if !errDto.is_empty() {
		c.JSON(status, errDto)
	} else {
		c.JSON(status, body)
	}
}

var log_request = func(c *gin.Context) {
	// Request received
	id := randstr.Hex(8)
	method := c.Request.Method
	uri := c.Request.RequestURI
	agent := c.Request.UserAgent()
	clientIP := c.ClientIP()
	logrus.WithField("comp", "gin").Infof("request received | id=%s, uri=%s, method=%s, ip=%s, agent=%s",
		id,
		uri,
		method,
		clientIP,
		agent)

	// Processing request
	startTime := time.Now()
	c.Next()
	endTime := time.Now()

	// Request processed
	latency := endTime.Sub(startTime)
	status := c.Writer.Status()
	logrus.WithField("comp", "gin").Infof("request processed | id=%s, status=%d, latency=%v",
		id,
		status,
		latency)
}
