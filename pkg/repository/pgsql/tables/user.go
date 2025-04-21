package tables

type User struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string `gorm:"unique;not null" json:"username"`
	Password  string `gorm:"not null" json:"password"`
	Email     string `gorm:"unique;not null" json:"email"`
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt string `gorm:"autoDeleteTime" json:"deleted_at"`
}
