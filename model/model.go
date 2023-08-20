package model

type VolumeInfo struct {
	Title    string `json:"title"`
	Language string `json:"language"`
}

type BookInformation struct {
	Id         string     `json:"id"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
}

type Items struct {
	Items []BookInformation `json:"items"`
}
