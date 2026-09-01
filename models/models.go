package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex" json:"username"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Tokens  []Token  `gorm:"foreignKey:UserID" json:"-"`
	Reviews []Review `gorm:"foreignKey:UserID" json:"-"`
}

type Token struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Token     string    `gorm:"uniqueIndex;type:text" json:"token"`
	TokenType string    `json:"token_type"` // BEARER
	Expired   bool      `json:"expired"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Movies []Movie `gorm:"many2many:movie_categories" json:"-"`
}

type Movie struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	ReleaseDate time.Time      `json:"release_date"`
	Director    string         `json:"director"`
	Duration    int            `json:"duration"` // in minutes
	Poster      string         `json:"poster"`   // URL
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Categories []Category `gorm:"many2many:movie_categories" json:"categories"`
	Reviews    []Review   `gorm:"foreignKey:MovieID" json:"reviews,omitempty"`
}

type Review struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MovieID   uint           `json:"movie_id"`
	UserID    uint           `json:"user_id"`
	Rating    int            `json:"rating"` // 1-10
	Comment   string         `gorm:"type:text" json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Movie Movie `gorm:"foreignKey:MovieID" json:"-"`
	User  User  `gorm:"foreignKey:UserID" json:"-"`
}
