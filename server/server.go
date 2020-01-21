package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/simon-ding/lklk/controller"
	"github.com/simon-ding/lklk/models"
	"net/http"
)

type Server struct {
	db *gorm.DB
	engine *gin.Engine
}

func New() *Server {
	s := &Server{}
	db := models.OpenDB()
	models.Migrate(db)
	s.db = db
	s.engine = gin.Default()

	s.engine.Use(s.inject())
	s.RegisterHandler()
	return s
}

func (s *Server) RegisterHandler() {
	s.Register("GET", "/import_yyets", &controller.ImportYYetsAPI{})
}

func (s Server) Register(method string, url string, handler controller.APIHandler) {
	switch method {
	case http.MethodGet:
		s.engine.GET(url, controller.Wrap(handler))
	case http.MethodPost:
		s.engine.POST(url, controller.Wrap(handler))
	}
}

func (s *Server) Run() error {
	return s.engine.Run(":8080")
}

func (s *Server) inject() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", s.db)
	}
}