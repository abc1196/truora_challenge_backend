package data

import (
	"time"
)

// Domain data type
type Domain struct {
	ID               string    `json:"id"`
	SslGrade         string    `json:"sslGrade"`
	PreviousSslGrade string    `json:"previousSslGrade"`
	Logo             string    `json:"logo"`
	Title            string    `json:"title"`
	ServersChanged   bool      `json:"serversChanged"`
	IsDown           bool      `json:"isDown"`
	Servers          []Server  `gorm:"foreignkey:DomainRefer" json:"servers"`
	Created          time.Time `gorm:"type:time" json:"created"`
	Updated          time.Time `gorm:"type:time" json:"updated"`
}
