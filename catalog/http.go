package catalog

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type httpServer struct {
	service Service
}

func (hS *httpServer) GetBookByID(c *gin.Context) error {
	id := c.Param("id")
	b, err := hS.service.GetBookByID(c, id)
	if err != nil {
		return err
	}
	test := &b
	return c.JSON(http.StatusOK, test)
}
func NewHttpServer(s Service) *httpServer {
	return &httpServer{s}
}
