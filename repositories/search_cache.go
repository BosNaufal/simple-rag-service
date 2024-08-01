package repositories

import (
	"bos_personal_ai/models"

	"gorm.io/gorm"
)

type SearchCacheRepositoryInterface interface {
	Find(embeddingString string) (SearchCacheEntity, error)
	Add(searchCacheInput SearchCacheEntity) (SearchCacheEntity, error)
}

type SearchCacheRepositoryImpl struct {
	db *gorm.DB
}

type SearchCacheEntity struct {
	ID        uint   `json:"id"`
	Query     string `json:"query"`
	Embedding string `json:"embedding"`
}

func NewSearchCacheRepository(db *gorm.DB) SearchCacheRepositoryInterface {
	return &SearchCacheRepositoryImpl{
		db: db,
	}
}

func (repo *SearchCacheRepositoryImpl) convertToSearchCacheEntity(k models.SearchCache) SearchCacheEntity {
	var result SearchCacheEntity
	result.ID = k.ID
	result.Query = k.Query
	result.Embedding = k.Embedding
	return result
}

func (repo *SearchCacheRepositoryImpl) Find(searchQuery string) (SearchCacheEntity, error) {
	var searchCache models.SearchCache

	result := repo.db.Where("query = ?", searchQuery).First(&searchCache)

	if result.Error != nil {
		return SearchCacheEntity{}, result.Error
	} // returns error

	return repo.convertToSearchCacheEntity(searchCache), nil
}

func (repo *SearchCacheRepositoryImpl) Add(searchCacheInput SearchCacheEntity) (SearchCacheEntity, error) {
	var searchCache models.SearchCache

	searchCache.Query = searchCacheInput.Query
	searchCache.Embedding = searchCacheInput.Embedding

	result := repo.db.Create(&searchCache) // pass pointer of data to Create

	if result.Error != nil {
		return SearchCacheEntity{}, result.Error
	} // returns error

	return repo.convertToSearchCacheEntity(searchCache), nil
}
