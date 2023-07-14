package item

import (
	"errors"
	"github.com/antoniokichaev/hezzl-collector/internal/repo/pgdb"
	"github.com/antoniokichaev/hezzl-collector/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ResponseRemove struct {
	Id         int  `json:"id"`
	CampaignId int  `json:"campaignId"`
	Removed    bool `json:"removed"`
}
type RequestCreate struct {
	Name string `json:"name" binding:"required"`
}

type RequestUpdate struct {
	RequestCreate
	Description string `json:"description,omitempty"`
}

type Handlers struct {
	service service.Item
}

func New(repo service.Item) *Handlers {
	return &Handlers{service: repo}
}

// Create принимает запрос ввида POST /item/?campaignId=int
func (h *Handlers) Create(c *gin.Context) {
	campaignId, err := strconv.Atoi(c.Query("campaignId"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req := RequestCreate{}
	err = c.ShouldBindJSON(&req)

	if err != nil || req.Name == "" {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	item, err := h.service.CreateItem(c.Request.Context(), req.Name, campaignId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// Update принимает запрос ввида PATCH /item/?campaignId=int
func (h *Handlers) Update(c *gin.Context) {
	itemId, campaignId, err := h.parseParam(c)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req := RequestUpdate{}
	err = c.ShouldBindJSON(&req)
	if err != nil || req.Name == "" {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	item, err := h.service.UpdateItem(c.Request.Context(), req.Name, req.Description, itemId, campaignId)
	if errors.Is(err, pgdb.ErrNotFoundItem) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "errors.item.notFound", "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Delete принимает запрос ввида Delete /item/remove
func (h *Handlers) Delete(c *gin.Context) {
	itemId, campaignId, err := h.parseParam(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"itemId": itemId, "campaignId": campaignId, "err": err})
		return
	}

	it, err := h.service.DeleteItem(c.Request.Context(), itemId, campaignId)
	if errors.Is(err, pgdb.ErrNotFoundItem) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "errors.item.notFound", "code": 3, "detail": "{}"})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ResponseRemove{Removed: it.Removed, Id: itemId, CampaignId: campaignId})
}

func (h *Handlers) GetItem(c *gin.Context) {
	itemId, campaignId, err := h.parseParam(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	allItems, err := h.service.GetItem(c.Request.Context(), itemId, campaignId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, allItems)
}

func (h *Handlers) GetItems(c *gin.Context) {
	allItems, err := h.service.GetItems(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, allItems)
}

func (h *Handlers) parseParam(c *gin.Context) (itemId, campaignId int, err error) {
	campaignId, err = strconv.Atoi(c.Query("campaignId"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	itemId, err = strconv.Atoi(c.Query("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	return
}
