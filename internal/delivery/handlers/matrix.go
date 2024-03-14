package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"template/internal/models"
	_ "template/internal/models/swagger"
	"template/internal/service"
	"template/pkg/tracing"
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

// CreateMatrixWithoutParent @Summary Create matrix without parent
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param data body swagger.MatrixBase true "Matrix create"
// @Success 200 {object} string "Successfully created matrix"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/create_no_parent [post]
func (m MatrixHandler) CreateMatrixWithoutParent(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), tracing.CreateMatrix)
	defer span.End()

	var matrixCreate models.MatrixBase

	if err := c.ShouldBindJSON(&matrixCreate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	name, err := m.service.CreateMatrixWithoutParent(ctx, matrixCreate)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.CreateMatrixType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, name)
}

// CreateMatrix @Summary Create matrix with parent
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param data body models.MatrixDifferenceRequest true "Matrix create"
// @Success 200 {object} string "Successfully created matrix"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/create [post]
func (m MatrixHandler) CreateMatrix(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), tracing.CreateMatrix)
	defer span.End()

	var matrixCreate models.MatrixDifferenceRequest

	if err := c.ShouldBindJSON(&matrixCreate); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	name, err := m.service.CreateMatrix(ctx, matrixCreate)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.CreateMatrixType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, name)
}

// GetHistory @Summary Get matrices by time start, time end and matrix type (can be null)
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param data body swagger.GetHistoryMatrix true "Get data"
// @Success 200 {object} []swagger.ResponseHistoryMatrix "Found matrices"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/get_history [put]
func (m MatrixHandler) GetHistory(c *gin.Context) {
	var getHistoryMatrix models.GetHistoryMatrix

	ctx, span := m.tracer.Start(c.Request.Context(), tracing.GetHistory)
	defer span.End()

	if err := c.ShouldBindJSON(&getHistoryMatrix); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	matrices, err := m.service.GetHistory(ctx, getHistoryMatrix)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.GetHistoryType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, matrices)
}

// GetDifference @Summary Get difference between two matrices
// @Description Retrieves the differences between two matrices identified by their names.
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param from_name query string true "Name of the first matrix"
// @Param to_name query string true "Name of the second matrix"
// @Success 200 {object} []models.MatrixDifferenceResponse "Found matrices differences"
// @Failure 400 {object} map[string]string "Invalid input, missing matrix names"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/get_difference [get]
func (m MatrixHandler) GetDifference(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), tracing.GetDifference)
	defer span.End()

	matrixName1, ok := c.GetQuery("from_name")
	if !ok {
		er := fmt.Errorf("bad `from_name` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "matrixName1 not provided"})
		return
	}
	matrixName2, ok := c.GetQuery("to_name")
	if !ok {
		er := fmt.Errorf("bad `to_name` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "matrixName2 not provided"})
		return
	}

	span.AddEvent(tracing.CallToService)
	matrices, err := m.service.GetDifference(ctx, matrixName1, matrixName2)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.GetDifferenceType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

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

	ctx, span := m.tracer.Start(c.Request.Context(), tracing.GetTendency)
	defer span.End()

	if err := c.ShouldBindJSON(&data); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	tendency, err := m.service.GetTendency(ctx, data)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.GetTendencyType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, tendency)
}

// GetMatrix @Summary Get matrix by name and page
// @Description Retrieves a specific page of the matrix identified by its name.
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param matrix_name query string true "Name of the matrix to retrieve"
// @Param page query int true "Page number of the matrix to retrieve"
// @Param microcategory_id query int false "Microcategory ID if you need"
// @Param region_id query int false "Region ID if you need"
// @Success 200 {object} []swagger.Matrix "Successfully retrieved the specified page of the matrix"
// @Failure 400 {object} map[string]string "Invalid input, missing or incorrect parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /matrix/get_matrix [get]
func (m MatrixHandler) GetMatrix(c *gin.Context) {
	var mc, rg null.Int

	ctx, span := m.tracer.Start(c.Request.Context(), tracing.GetMatrix)
	defer span.End()

	matrixName, ok := c.GetQuery("matrix_name")
	if !ok {
		er := fmt.Errorf("bad `matrix_name` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "matrixName not provided"})
		return
	}
	pageStr, ok := c.GetQuery("page")
	if !ok {
		er := fmt.Errorf("bad `page` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "pageStr not provided"})
		return
	}
	mcStr, ok := c.GetQuery("microcategory_id")
	if ok {
		mcInt, err := strconv.Atoi(mcStr)
		if err != nil {
			span.RecordError(err, trace.WithAttributes(
				attribute.String(tracing.QueryType, err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "page can't interpret as int"})
			return
		}
		mc.Int64 = int64(mcInt)
		mc.Valid = true
	}
	pgStr, ok := c.GetQuery("region_id")
	if ok {
		pgInt, err := strconv.Atoi(pgStr)
		if err != nil {
			span.RecordError(err, trace.WithAttributes(
				attribute.String(tracing.QueryType, err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "page can't interpret as int"})
			return
		}
		rg.Int64 = int64(pgInt)
		rg.Valid = true
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.QueryType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "page can't interpret as int"})
		return
	}

	span.AddEvent(tracing.CallToService)
	matrices, err := m.service.GetMatrix(ctx, matrixName, mc, rg, page)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.GetMatrixType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, matrices)
}

// GetMatrixPages @Summary Retrieve a page of matrices
// @Description Retrieves a specific page of matrices based on the matrix name provided.
// @Tags matrix
// @Accept  json
// @Produce  json
// @Param matrix_name query string true "The name of the matrix for which to retrieve the page"
// @Success 200 {object} int "Successfully retrieved the total number of pages for the specified matrix"
// @Failure 400 {object} map[string]string "Invalid input, missing or incorrect matrix_name parameter"
// @Failure 500 {object} map[string]string "Internal server error occurred while retrieving the matrix page"
// @Router /matrix/get_matrix_pages [get]
func (m MatrixHandler) GetMatrixPages(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), tracing.GetMatrixPages)
	defer span.End()

	matrixName, ok := c.GetQuery("matrix_name")
	if !ok {
		er := fmt.Errorf("bad `matrix_name` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "matrixName not provided"})
		return
	}

	span.AddEvent(tracing.CallToService)
	matrices, err := m.service.GetMatrixPages(ctx, matrixName)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.GetMatrixPagesType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

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
	ctx, span := m.tracer.Start(c.Request.Context(), tracing.GetMatricesByDuration)
	defer span.End()

	timeFromStr, ok := c.GetQuery("time_from")
	if !ok {
		er := fmt.Errorf("bad `time_from` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "timeFrom not provided"})
		return
	}
	timeToStr, ok := c.GetQuery("time_to")
	if !ok {
		er := fmt.Errorf("bad `time_to` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "timeTo not provided"})
		return
	}

	timeFrom, err := time.Parse(time.RFC3339, timeFromStr)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.TimeFormatType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timeFrom format"})
		return
	}

	timeTo, err := time.Parse(time.RFC3339, timeToStr)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.TimeFormatType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timeTo format"})
		return
	}

	span.AddEvent(tracing.CallToService)
	matrices, err := m.service.GetMatricesByDuration(ctx, timeFrom, timeTo)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.GetMatricesByDuration, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, matrices)
}
