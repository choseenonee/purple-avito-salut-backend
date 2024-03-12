package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"template/internal/models"
	"template/internal/service"
)

type Handler struct {
	service service.Service
	tracer  trace.Tracer
}

func InitHandler(service service.Service, tracer trace.Tracer) Handler {
	return Handler{
		service: service,
		tracer:  tracer,
	}
}

// GetPrice @Summary Get price
// @Tags price
// @Accept  json
// @Produce  json
// @Param data body models.InData true "Get price"
// @Success 200 {object} models.OutData "Successfully responsed with price"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /price [put]
func (h Handler) GetPrice(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), CreateMatrix)
	defer span.End()

	var inData models.InData

	if err := c.ShouldBindJSON(&inData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(CallToService)
	response, err := h.service.GetMicroCategoryPath(ctx, inData.MicroCategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
