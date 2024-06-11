package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type User struct {
	Id        int        `gorm:"type:int;primary_key;autoIncrement" json:"id,omitempty"`
	Role      string     `gorm:"type:varchar(10)" json:"role"`
	Name      string     `gorm:"type:varchar(255)" json:"name"`
	Email     string     `gorm:"type:varchar(50)" json:"email"`
	Password  string     `gorm:"type:varchar(255)" json:"password,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Tasks     []Task     `gorm:"constraint:OnDelete:CASCADE" json:"tasks,omitempty"` // has many
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("userId = ?", u.Id).Delete(&Task{})
	return
}
