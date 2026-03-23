package models

import "time"

// JumpType represents valid skydiving jump discipline types
type JumpType string

const (
	JumpTypeFun       JumpType = "Fun"
	JumpTypeFormation JumpType = "Formation (FS/VFS)"
	JumpTypeAFF       JumpType = "AFF"
	JumpTypeTandem    JumpType = "Tandem"
	JumpTypeWingsuit  JumpType = "Wingsuit"
	JumpTypeCRW       JumpType = "CRW"
	JumpTypeFreefly   JumpType = "Freefly"
	JumpTypeAngle     JumpType = "Angle"
	JumpTypeTracking  JumpType = "Tracking"
	JumpTypeHopNPop   JumpType = "Hop & Pop"
	JumpTypeDemo      JumpType = "Demo"
	JumpTypeOther     JumpType = "Other"
)

// Jump represents a single skydive entry in the logbook.
type Jump struct {
	ID                 uint      `gorm:"primaryKey"        json:"id"`
	JumpNumber         int       `gorm:"uniqueIndex;not null" json:"jumpNumber"`
	Date               time.Time `gorm:"not null"          json:"date"`
	Location           string    `json:"location"`
	Dropzone           string    `json:"dropzone"`
	Aircraft           string    `json:"aircraft"`
	ExitAltitude       int       `json:"exitAltitude"`       // feet
	DeploymentAltitude int       `json:"deploymentAltitude"` // feet
	FreefallTime       int       `json:"freefallTime"`       // seconds
	CanopySize         string    `json:"canopySize"`
	JumpType           string    `json:"jumpType"`
	Notes              string    `json:"notes"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
