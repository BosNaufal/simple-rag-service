package domain

import (
	"bos_personal_ai/repositories"
)

type SearchCacheDomainInterface interface {
	Find(searchQuery string) (repositories.SearchCacheEntity, error)
	Add(searchQuery string, embeddingString string) (repositories.SearchCacheEntity, error)
}

type SearchCacheDomainImpl struct {
	SearchCacheRepo repositories.SearchCacheRepositoryInterface
}

func NewSearchCacheDomain(searchrepo repositories.SearchCacheRepositoryInterface) *SearchCacheDomainImpl {
	return &SearchCacheDomainImpl{
		SearchCacheRepo: searchrepo,
	}
}

func (domain *SearchCacheDomainImpl) Find(searchQuery string) (repositories.SearchCacheEntity, error) {
	result, err := domain.SearchCacheRepo.Find(searchQuery)
	return result, err
}

func (domain *SearchCacheDomainImpl) Add(searchQuery string, embeddingString string) (repositories.SearchCacheEntity, error) {
	SearchCache := repositories.SearchCacheEntity{
		Query:     searchQuery,
		Embedding: embeddingString,
	}

	return domain.SearchCacheRepo.Add(SearchCache)
}
