package repository

import "time"

type Restriction struct {
	Crawler   string `gorm:"type:varchar(128);uniqueIndex:crawler_job_created"`
	Job       string `gorm:"type:varchar(128);uniqueIndex:crawler_job_created"`
	StartDate time.Time
	EndDate   time.Time
	HourFrom  int
	HourTo    int
	CreatedAt time.Time `gorm:"autoCreateTime;uniqueIndex:crawler_job_created"`
}

func (Restriction) TableName() string {
	return "restriction"
}

type Event struct {
	Crawler   string    `gorm:"type:varchar(128)"`
	Job       string    `gorm:"type:varchar(128)"`
	UserName  string    `gorm:"type:varchar(128)"`
	ID        string    `gorm:"type:varchar(128);uniqueIndex:id_name"`
	Name      string    `gorm:"type:varchar(128);uniqueIndex:id_name"`
	Message   string    `gorm:"type:varchar(1024)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Event) TableName() string {
	return "event"
}

type User struct {
	ID   string `gorm:"type:varchar(128);primary_key"`
	Name string `gorm:"type:varchar(128);index"`
}

func (User) TableName() string {
	return "user"
}
