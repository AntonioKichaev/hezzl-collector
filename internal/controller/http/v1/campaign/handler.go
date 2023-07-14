package campaign

import (
	"errors"
	"github.com/antoniokichaev/hezzl-collector/internal/repo/pgdb"
	"github.com/antoniokichaev/hezzl-collector/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handlers struct {
	campaignService service.Campaign
}

func New(campaignService service.Campaign) *Handlers {
	return &Handlers{campaignService: campaignService}
}

type RequestCreate struct {
	Name string `json:"name" binding:"required"`
}
type RequestUpdate struct {
	RequestCreate
}

// Create принимает запрос ввида POST /item/create
func (h *Handlers) Create(c *gin.Context) {
	req := &RequestCreate{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	comp, err := h.campaignService.CreateCampaign(c.Request.Context(), req.Name)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, comp)
}

// Update принимает запрос ввида PATCH /item/update?id=int
func (h *Handlers) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	req := &RequestCreate{}
	err = c.ShouldBindJSON(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	comp, err := h.campaignService.UpdateCampaign(c.Request.Context(), req.Name, id)
	if errors.Is(err, pgdb.ErrNotFoundCampaign) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "errors.item.notFound", "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comp)
}

// Delete принимает запрос ввида DELETE /item/remove?id=int
func (h *Handlers) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	isDeleted, err := h.campaignService.DeleteCampaign(c.Request.Context(), id)
	if errors.Is(err, pgdb.ErrNotFoundCampaign) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "errors.item.notFound", "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "removed": isDeleted})
}

// GetCampaigns принимает запрос ввида GET /campaigns/list
func (h *Handlers) GetCampaigns(c *gin.Context) {
	allItems, err := h.campaignService.GetCampaigns(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, allItems)
}

// GetCampaigns принимает запрос ввида GET /campaigns/list
func (h *Handlers) GetCampaign(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return

	}
	allItems, err := h.campaignService.GetCampaign(c.Request.Context(), id)
	if errors.Is(err, pgdb.ErrNotFoundCampaign) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "errors.item.notFound", "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, allItems)
}
