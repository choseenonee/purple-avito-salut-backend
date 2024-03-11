package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"template/internal/models"
	_ "template/internal/models/swagger"
	"template/internal/service"
	"time"
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

// GetTendency @Summary Get price tendency
// @Description Retrieves price difference in time span
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param data body models.GetTendencyNode true "Get data"
// @Success 200 {object} []models.ResponseTendencyNode "Found prices in time span and one before it"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/get_tendency [put]
func (m MatrixHandler) GetTendency(c *gin.Context) {
	var data models.GetTendencyNode

	ctx, span := m.tracer.Start(c.Request.Context(), GetHistory)
	defer span.End()

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(CallToService)
	tendency, err := m.service.GetTendency(ctx, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tendency)
}

// GetMatrix @Summary Get matrix by name and page
// @Description Retrieves a specific page of the matrix identified by its name.
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param matrix_name query string true "Name of the matrix to retrieve"
// @Param page query int true "Page number of the matrix to retrieve"
// @Success 200 {object} []models.Matrix "Successfully retrieved the specified page of the matrix"
// @Failure 400 {object} map[string]string "Invalid input, missing or incorrect parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/get_matrix [get]
func (m MatrixHandler) GetMatrix(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), GetHistory)
	defer span.End()

	matrixName, ok := c.GetQuery("matrix_name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "matrixName not provided"})
		return
	}
	pageStr, ok := c.GetQuery("page")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pageStr not provided"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page can't interpret as int"})
		return
	}

	span.AddEvent(CallToService)
	matrices, err := m.service.GetMatrix(ctx, matrixName, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matrices)
}

// GetMatricesByDuration @Summary Get matrices by duration
// @Description Retrieves matrices that fall within the specified time duration.
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param time_from query string true "Start time of the duration (RFC3339 format)"
// @Param time_to query string true "End time of the duration (RFC3339 format)"
// @Success 200 {object} []models.Matrix "Successfully retrieved matrices within the specified duration"
// @Failure 400 {object} map[string]string "Invalid input, missing or incorrect parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/get_matrices_by_duration [get]
func (m MatrixHandler) GetMatricesByDuration(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), "GetMatricesByDuration")
	defer span.End()

	timeFromStr, ok := c.GetQuery("time_from")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "timeFrom not provided"})
		return
	}
	timeToStr, ok := c.GetQuery("time_to")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "timeTo not provided"})
		return
	}

	timeFrom, err := time.Parse(time.RFC3339, timeFromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timeFrom format"})
		return
	}

	timeTo, err := time.Parse(time.RFC3339, timeToStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timeTo format"})
		return
	}

	span.AddEvent("CallToService")
	matrices, err := m.service.GetMatricesByDuration(ctx, timeFrom, timeTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matrices)
}
