package main

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

type SMSReq struct {
	SysCode    string `json:"sysCode"`
	SysName    string `json:"sysName"`
	EventType  string `json:"eventType"`
	EventLevel string `json:"eventLevel"`
	// 系统负责人
	ProjectManager string `json:"projectManager"`
	// DB-AGILEX+年月日时分秒+序列号
	EventID        string    `json:"eventid" validate:"required"`
	EventSummary   string    `json:"eventSummary" validate:"required"`
	EventDetail    string    `json:"eventDetail"`
	EventDatetime  time.Time `json:"eventDatetime" validate:"required"`
	EventIP        string    `json:"eventIp"`
	AlertType      string    `json:"alertType"`
	AlertCellphone string    `json:"alertCellphone"`
	AlertEhr       string    `json:"alertEhr"`
}

type SMSRecover struct {
	Address           string    `json:"-"`
	EventID           string    `json:"eventid" validate:"required"`
	EventRecoveryType string    `json:"eventRecoveryType" `
	EventRecoveryTime time.Time `json:"eventRecoveryTime" validate:"required"`
}

func main() {
	r := gin.Default()
	r.POST("/api/monitor/event/add", func(c *gin.Context) {
		req := &SMSReq{}
		if err := c.ShouldBind(req); err != nil {
			c.JSON(400, gin.H{
				"code": 400,
				"msg":  err.Error(),
			})
			return
		}
		slog.Info("req", slog.Any("req", req))
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "成功",
		})
	})
	r.POST("/api/monitor/event/recover", func(c *gin.Context) {
		req := &SMSRecover{}
		if err := c.ShouldBind(req); err != nil {
			c.JSON(400, gin.H{
				"code": 400,
				"msg":  err.Error(),
			})
			return
		}
		slog.Info("req", slog.Any("req", req))
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "成功",
		})
	})
	// 启动服务
	if err := r.Run(":9090"); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
	}
}
