package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"template/internal/models"
	_ "template/internal/models/swagger"
	"template/internal/service"
)

type MatrixHandler struct {
	service service.Matrix
	tracer  trace.Tracer
}

func InitMatrixHandler(service service.Matrix, tracer trace.Tracer) MatrixHandler {
	return MatrixHandler{
		service: service,
		tracer:  tracer,
	}
}

// @Summary Create matrix
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param data body swagger.MatrixBase true "Matrix create"
// @Success 200 {object} string "Successfully created matrix"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/create [post]
func (m MatrixHandler) CreateMatrix(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), CreateMatrix)
	defer span.End()

	var matrixCreate models.MatrixBase

	if err := c.ShouldBindJSON(&matrixCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(CallToService)
	name, err := m.service.Create(ctx, matrixCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, name)
}

// @Summary Get matrixes by time start, time end and matrix type (can be null)
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param data body swagger.GetHistoryMatrix true "Get data"
// @Success 200 {object} []swagger.ResponseHistoryMatrix "Found matrixes"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/get_history [put]
func (m MatrixHandler) GetHistory(c *gin.Context) {
	var getHistoryMatrix models.GetHistoryMatrix

	if err := c.ShouldBindJSON(&getHistoryMatrix); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	matrixes, err := m.service.GetHistory(ctx, getHistoryMatrix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matrixes)
}
