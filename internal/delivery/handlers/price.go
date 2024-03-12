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
	update  service.Update
	tracer  trace.Tracer
}

func InitHandler(service service.Service, update service.Update, tracer trace.Tracer) Handler {
	return Handler{
		service: service,
		update:  update,
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
	ctx, span := h.tracer.Start(c.Request.Context(), GetPrice)
	defer span.End()

	var inData models.InData

	if err := c.ShouldBindJSON(&inData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(CallToService)
	response, err := h.service.GetPrice(ctx, inData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RecalculateRedis @Summary recalculates hops and updates price information
// @Tags price
// @Accept  json
// @Produce  json
// @Param name body models.UpdateMatrixName true "Get price"
// @Success 204 "Successfully recalculated and updated price information"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /recalculate [put]
func (h Handler) RecalculateRedis(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), RecalculateRedis)
	defer span.End()

	var newBaseMatrixName models.UpdateMatrixName
	if err := c.ShouldBindJSON(&newBaseMatrixName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cringe data provided"})
		return
	}

	span.AddEvent(CallToService)
	h.update.ReRunInit(ctx, newBaseMatrixName.Name)

	c.Status(http.StatusNoContent)
}
