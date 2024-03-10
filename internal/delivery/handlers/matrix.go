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


// CreateMatrix @Summary Create matrix
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

// GetHistory @Summary Get matrices by time start, time end and matrix type (can be null)
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

	ctx, span := m.tracer.Start(c.Request.Context(), GetHistory)
	defer span.End()

	if err := c.ShouldBindJSON(&getHistoryMatrix); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(CallToService)
	matrices, err := m.service.GetHistory(ctx, getHistoryMatrix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matrices)
}

// GetDifference @Summary Get difference between two matrices
// @Description Retrieves the differences between two matrices identified by their names.
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param from_name query string true "Name of the first matrix"
// @Param to_name query string true "Name of the second matrix"
// @Success 200 {object} []models.MatrixDifference "Found matrices differences"
// @Failure 400 {object} map[string]string "Invalid input, missing matrix names"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/get_difference [get]
func (m MatrixHandler) GetDifference(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), GetHistory)
	defer span.End()

	matrixName1, ok := c.GetQuery("from_name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "matrixName1 not provided"})
		return
	}
	matrixName2, ok := c.GetQuery("to_name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "matrixName2 not provided"})
		return
	}

	span.AddEvent(CallToService)
	matrices, err := m.service.GetDifference(ctx, matrixName1, matrixName2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matrices)
}