package models

// equals
type SearchCache struct {
	ID        uint `gorm:"primaryKey"`
	Query     string
	Embedding string
}

// TableName overrides the table name used by User to `profiles`
func (SearchCache) TableName() string {
	return "search_caches"
}
