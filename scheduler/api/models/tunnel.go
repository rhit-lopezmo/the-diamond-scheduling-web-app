package models

type Tunnel struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}
