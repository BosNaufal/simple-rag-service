package infra_services

import (
	"bos_personal_ai/repositories"

	"gorm.io/gorm"
)

type SearchCachesInterface interface {
	FindEmbeddingInCache(searchQuery string) (repositories.SearchCacheEntity, error)
	SaveQueryInCache(searchQuery string, embeddingString string) error
}

type SearchCachesImpl struct {
	repo repositories.SearchCacheRepositoryInterface
}

func NewSearchCacheService(
	repo repositories.SearchCacheRepositoryInterface,
) *SearchCachesImpl {
	return &SearchCachesImpl{
		repo: repo,
	}
}

func (srv *SearchCachesImpl) FindEmbeddingInCache(searchQuery string) (repositories.SearchCacheEntity, error) {
	searchCache, err := srv.repo.Find(searchQuery)
	if err != nil && err != gorm.ErrRecordNotFound {
		return repositories.SearchCacheEntity{}, err
	}

	return searchCache, err
}

func (srv *SearchCachesImpl) SaveQueryInCache(searchQuery string, embeddingString string) error {
	_, err := srv.repo.Add(repositories.SearchCacheEntity{
		Query:     searchQuery,
		Embedding: embeddingString,
	})
	return err
}
