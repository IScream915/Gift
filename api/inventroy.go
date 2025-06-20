package api

import (
	"fmt"
	"gift/dto"
	"gift/pkg/response"
	"gift/services"
	"github.com/gin-gonic/gin"
)

type Inventory interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
	Load(c *gin.Context)
}

func NewInventory(svc services.Inventory) Inventory {
	return &inventory{
		svc: svc,
	}
}

type inventory struct {
	svc services.Inventory
}

// Create POST
func (obj *inventory) Create(c *gin.Context) {
	req := &dto.CreateInventoryReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.Json(c, response.WithErr(err))
		return
	}

	if err := obj.svc.Create(c, req); err != nil {
		response.Json(c, response.WithErr(err))
		return
	}

	response.Json(c, response.WithMsg("create inventory success"))
}

// Update POST
func (obj *inventory) Update(c *gin.Context) {
	req := &dto.UpdateInventoryReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.Json(c, response.WithErr(err))
		return
	}

	if err := obj.svc.Update(c, req); err != nil {
		response.Json(c, response.WithErr(err))
		return
	}

	response.Json(c, response.WithMsg("update inventory success"))
}

// Delete DELETE
func (obj *inventory) Delete(c *gin.Context) {
	req := &dto.DeleteInventoryReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.Json(c, response.WithErr(err))
		return
	}

	if err := obj.svc.Delete(c, req); err != nil {
		response.Json(c, response.WithErr(err))
		return
	}

	response.Json(c, response.WithMsg("delete inventory success"))
}

func (obj *inventory) Get(c *gin.Context) {
	req := &dto.GetInventoriesReq{}

	if err := c.ShouldBindQuery(req); err != nil {
		response.Json(c, response.WithErr(err))
		return
	}

	resp, err := obj.svc.GetInventories(c, req)
	if err != nil {
		response.Json(c, response.WithErr(err))
		return
	}

	response.Json(c, response.WithData(resp))
}

func (obj *inventory) Load(c *gin.Context) {
	total, err := obj.svc.LoadInventories(c)
	if err != nil {
		response.Json(c, response.WithErr(err))
		return
	}

	response.Json(c, response.WithMsg(fmt.Sprintf("%d load inventory success", total)))
}
