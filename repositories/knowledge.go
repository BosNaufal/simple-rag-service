package repositories

import (
	"bos_personal_ai/models"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type KnowledgeRepositoryInterface interface {
	Find(embeddingString string, limit int) ([]KnowledgeEntity, error)
	Add(knowledgeInput KnowledgeEntity) (KnowledgeEntity, error)
	Update(knowledgeInput KnowledgeEntity) (KnowledgeEntity, error)
	UpdateEmbedding(ID uint, embeddingString string) (KnowledgeEntity, error)
	Delete(id uint) error
}

type KnowledgeRepositoryImpl struct {
	db *gorm.DB
}

type KnowledgeEntity struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Embedding string    `json:"embedding,omitempty"`
}

func NewKnowledgeRepository(db *gorm.DB) *KnowledgeRepositoryImpl {
	return &KnowledgeRepositoryImpl{
		db: db,
	}
}

func (repo *KnowledgeRepositoryImpl) convertToKnowledgeEntity(k models.Knowledge) KnowledgeEntity {
	var result KnowledgeEntity
	result.ID = k.ID
	result.Title = k.Title
	result.Content = k.Content
	result.UpdatedAt = k.UpdatedAt
	result.CreatedAt = k.CreatedAt
	return result
}

func (repo *KnowledgeRepositoryImpl) Find(embeddingString string, limit int) ([]KnowledgeEntity, error) {
	var knowledgeList []models.Knowledge

	result := repo.db.Raw(
		"SELECT id, title, content, created_at, updated_at, embedding <-> ? AS distance FROM note_chunks ORDER BY distance LIMIT ?;",
		embeddingString,
		limit).
		Scan(&knowledgeList)

	if result.Error != nil {
		return []KnowledgeEntity{}, result.Error
	} // returns error

	modefiedResult := make([]KnowledgeEntity, len(knowledgeList))
	for i, v := range knowledgeList {
		modefiedResult[i] = repo.convertToKnowledgeEntity(v)
	}

	return modefiedResult, nil
}

func (repo *KnowledgeRepositoryImpl) Add(knowledgeInput KnowledgeEntity) (KnowledgeEntity, error) {
	var knowledge models.Knowledge

	knowledge.Title = knowledgeInput.Title
	knowledge.Content = knowledgeInput.Content
	knowledge.Embedding = sql.NullString{
		String: knowledgeInput.Embedding,
		Valid:  true,
	}

	result := repo.db.Create(&knowledge) // pass pointer of data to Create

	if result.Error != nil {
		return KnowledgeEntity{}, result.Error
	} // returns error

	return repo.convertToKnowledgeEntity(knowledge), nil
}

func (repo *KnowledgeRepositoryImpl) Update(knowledgeInput KnowledgeEntity) (KnowledgeEntity, error) {
	var knowledge models.Knowledge

	result := repo.db.First(&knowledge, knowledgeInput.ID)
	if result.Error != nil {
		return KnowledgeEntity{}, result.Error
	}

	knowledge.Title = knowledgeInput.Title
	knowledge.Content = knowledgeInput.Content

	result = repo.db.Save(&knowledge) // pass pointer of data to Create

	if result.Error != nil {
		return KnowledgeEntity{}, result.Error
	} // returns error

	return repo.convertToKnowledgeEntity(knowledge), nil
}

func (repo *KnowledgeRepositoryImpl) UpdateEmbedding(ID uint, embeddingString string) (KnowledgeEntity, error) {
	var knowledge models.Knowledge

	result := repo.db.First(&knowledge, ID)
	if result.Error != nil {
		return KnowledgeEntity{}, result.Error
	}

	knowledge.Embedding = sql.NullString{
		String: embeddingString,
		Valid:  true,
	}

	result = repo.db.Save(&knowledge) // pass pointer of data to Create

	if result.Error != nil {
		return KnowledgeEntity{}, result.Error
	} // returns error

	return repo.convertToKnowledgeEntity(knowledge), nil
}

func (repo *KnowledgeRepositoryImpl) Delete(id uint) error {
	var knowledge models.Knowledge

	result := repo.db.First(&knowledge, id)
	if result.Error != nil {
		return result.Error
	}

	result = repo.db.Delete(&knowledge) // pass pointer of data to Create

	if result.Error != nil {
		return result.Error
	} // returns error

	return nil
}
