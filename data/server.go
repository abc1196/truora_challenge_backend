package data

// Server data type
type Server struct {
	Address     string `gorm:"primary_key" json:"address"`
	SslGrade    string `json:"sslGrade"`
	Country     string `json:"country"`
	Owner       string `json:"owner"`
	DomainRefer string `json:"domain_id"`
}
