package assemble

import (
	"gift/dto"
	"gift/repo/models"
)

func Model2Info(list []*models.Inventory) []*dto.InventoryInfo {
	infoList := make([]*dto.InventoryInfo, 0)
	for _, v := range list {
		infoList = append(infoList, &dto.InventoryInfo{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Picture:     v.Picture,
			Price:       v.Price,
			Count:       v.Count,
			CreatedAt:   uint64(v.CreatedAt.Unix()),
			UpdatedAt:   uint64(v.UpdatedAt.Unix()),
		})
	}
	return infoList
}
