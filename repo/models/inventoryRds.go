package models

import "encoding/json"

type InventoryRds struct {
	InventoryId int    `json:"inventory_id"` // 物品ID
	Name        string `json:"name"`         // 物品名称
	Price       uint64 `json:"price"`        // 物品价格
	Count       uint64 `json:"count"`        // 物品库存
}

// MarshalBinary Redis序列化
func (obj *InventoryRds) MarshalBinary() ([]byte, error) {
	return json.Marshal(obj)
}

// UnmarshalBinary Redis反序列化
func (obj *InventoryRds) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &obj)
}
