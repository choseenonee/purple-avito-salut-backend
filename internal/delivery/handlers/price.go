package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"template/internal/models"
	"template/internal/service"
)

type UpdateHandler struct {
	service service.Update
	tracer  trace.Tracer
}

func InitMUpdateHandler(service service.Update, tracer trace.Tracer) UpdateHandler {
	return UpdateHandler{
		service: service,
		tracer:  tracer,
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
	ctx, span := m.tracer.Start(c.Request.Context(), CreateMatrix)
	defer span.End()

	var storage models.StorageBase

	if err := c.ShouldBindJSON(&storage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.AddEvent(CallToService)
	preparedStorage, err := m.service.PrepareStorage(ctx, storage.BaseLineMatrixName, storage.DiscountMatrixNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// TODO: change url

	discountMatrixNames := make([]string, 0, len(preparedStorage.DiscountMatrices))

	for _, i := range preparedStorage.DiscountMatrices {
		discountMatrixNames = append(discountMatrixNames, i.Name)
	}

	err = m.service.SendUpdatedStorage("http://localhost:8000", models.PreparedStorageSend{
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
	_, span := m.tracer.Start(c.Request.Context(), CreateMatrix)
	defer span.End()

	span.AddEvent(CallToService)
	err := m.service.SwitchStorage("http://localhost:8000")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
	_, span := m.tracer.Start(c.Request.Context(), CreateMatrix)
	defer span.End()

	c.JSON(http.StatusOK, m.service.GetCurrentStorage())
}
