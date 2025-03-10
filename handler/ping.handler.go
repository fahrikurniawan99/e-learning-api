package handler

import "github.com/gin-gonic/gin"

type PingHandler interface {
	Ping(c *gin.Context)
}

type PingHandlerImplement struct {
}

func NewPingHandler() PingHandler {
	return &PingHandlerImplement{}
}

// PingHandler contoh handler
// @Summary Ping API
// @Description API untuk mengecek apakah server berjalan
// @Tags Health Check
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func (ph PingHandlerImplement) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})

}
