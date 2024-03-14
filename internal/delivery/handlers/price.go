package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"template/internal/models"
	"template/internal/service"
	"template/pkg/tracing"
)

type UpdateHandler struct {
	service      service.Update
	tracer       trace.Tracer
	priceApiURLs []string
}

func InitMUpdateHandler(service service.Update, tracer trace.Tracer, priceApiURLs []string) UpdateHandler {
	return UpdateHandler{
		service:      service,
		tracer:       tracer,
		priceApiURLs: priceApiURLs,
	}
}

// PrepareAndSendStorage @Summary Prepare and send storage
// @Tags storage
// @Accept  json
// @Produce  json
// @Param data body models.StorageBase true "Storage create"
// @Success 200 {object} string "Successfully created matrix"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /storage/send [post]
func (m UpdateHandler) PrepareAndSendStorage(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), tracing.PrepareAndSendStorage)
	defer span.End()

	var storage models.StorageBase

	if err := c.ShouldBindJSON(&storage); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(tracing.CallToService)
	preparedStorage, err := m.service.PrepareStorage(ctx, storage.BaseLineMatrixName, storage.DiscountMatrixNames)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.PrepareAndSendStorageType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	discountMatrixNames := make([]string, 0, len(preparedStorage.DiscountMatrices))

	for _, i := range preparedStorage.DiscountMatrices {
		discountMatrixNames = append(discountMatrixNames, i.Name)
	}

	//TODO: change hardcoded url
	for _, url := range m.priceApiURLs {
		err = m.service.SendUpdatedStorage(fmt.Sprintf("http://%v:8000", url), models.PreparedStorageSend{
			StorageBase: models.StorageBase{
				BaseLineMatrixName:  preparedStorage.BaseLineMatrix.Name,
				DiscountMatrixNames: discountMatrixNames,
			},
			MicroCategoryHops: preparedStorage.MicroCategoryHops,
			RegionHops:        preparedStorage.RegionHops,
			DiscountHops:      preparedStorage.DiscountHops,
			SegmentDiscount:   preparedStorage.SegmentDiscount,
		})
		if err != nil {
			span.RecordError(err, trace.WithAttributes(
				attribute.String(tracing.MakeRequestType, err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, gin.H{"detail": "successfully!"})
}

// SwitchStorageToNext @Summary Switch storage on nodes
// @Tags storage
// @Accept  json
// @Produce  json
// @Success 200 {object} string "Successfully switched matrix"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /storage/switch [post]
func (m UpdateHandler) SwitchStorageToNext(c *gin.Context) {
	_, span := m.tracer.Start(c.Request.Context(), tracing.SwitchStorageToNext)
	defer span.End()

	span.AddEvent(tracing.CallToService)
	for _, url := range m.priceApiURLs {
		err := m.service.SwitchStorage(fmt.Sprintf("http://%v:8000", url))
		if err != nil {
			span.RecordError(err, trace.WithAttributes(
				attribute.String(tracing.MakeRequestType, err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, gin.H{"detail": "successfully!!!"})
}

// GetCurrentStorage @Summary Get current storage (nodes use this handler on startup)
// @Tags storage
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PreparedStorageSend "Successfully switched matrix"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /storage/current [get]
func (m UpdateHandler) GetCurrentStorage(c *gin.Context) {
	_, span := m.tracer.Start(c.Request.Context(), tracing.GetCurrentStorage)
	defer span.End()

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, m.service.GetCurrentStorage())
}
