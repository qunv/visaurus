package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qunv/visaurus/core/responses"
	"net/http"
)

type SymnonymController struct {
}

func NewSymnonymController() *SymnonymController {
	return &SymnonymController{}
}

func (s *SymnonymController) List(ctx *gin.Context) {
	responses.OfSucceed(ctx, http.StatusOK, "ok")
}
