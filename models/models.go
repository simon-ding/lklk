package models

import "github.com/jinzhu/gorm"

type TVSeries struct {
	gorm.Model
	NameCN   string
	NameEN   string
	YyetsID   string
	DoubanID string
	ImdbID   string
}

type Episode struct {
	gorm.Model
	TVSeriesID int
	Season     int
	Episode    int
	Name       string
}

type Task struct {
	gorm.Model
	MagnetLink string
	Status Status
}

type Status int

const (
	StatusAdded Status = iota
	StatusRunning
	StatusStopped
	StatusSeeding
)


type RSS struct {
	gorm.Model
	URL string
}