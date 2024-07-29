package app_services

import (
	"bos_personal_ai/thirdparties"
)

type RAGInterface interface {
	AskQuestion(question string) (string, error)
}

type RAGImpl struct {
	embeddingThirdparty      thirdparties.EmbeddingThirdPartyInterface
	aiChatThirdparty         thirdparties.AIChatInterface
	embeddedKnowledgeService EmbeddedKnowledgeServiceInterface
}

func NewRAG(
	embeddingThirdparty thirdparties.EmbeddingThirdPartyInterface,
	aiChatThirdparty thirdparties.AIChatInterface,
	embeddedKnowledgeService EmbeddedKnowledgeServiceInterface,
) *RAGImpl {
	return &RAGImpl{
		embeddingThirdparty:      embeddingThirdparty,
		aiChatThirdparty:         aiChatThirdparty,
		embeddedKnowledgeService: embeddedKnowledgeService,
	}
}

func (srv *RAGImpl) AskQuestion(question string) (string, error) {

	systemPrompt := `
	you're helpful personal assistant.
	Answer given QUESTION.
	provide ANSWER with user's preferred language based on QUESTION.
	answer with young voice tone.
	Provide answer only based on the given CONTEXT.
	the answer should be relevant with the CONTEXT.
	if the QUESTION is far from the given CONTEXT, please answer that you don't have enough information to answer the QUESTION.
	Provide concise, clear, and calming tone answer.
	Provide the REFERENCES too with the same format as the examples below.

	EXAMPLE 1:
	CONTEXT:

	Manusia diciptakan untuk beribadah
	--source: https://a.com/

	Salah satu bentuk ibadah adalah Sholat
	--source: https://yyy.com/

	Allah mencintai orang-orang yang beribadah kepadanya dengan ikhlas
	--source: https:/hhh.com

	QUESTION:
	Untuk apa manusia diciptakan dan bagaimana cara melaksanakan tujuan tersebut?

	ANSWER:
	Manusia diciptakan untuk beribadah kepada Allah [1].  Salah satu bentuk ibadah adalah Sholat. Manusia bisa melaksanakan sholat untuk memenuhi tujuan hidup agar Allah cinta kepadanya [2][3].

	---
	**REFERENCES**
	- https://a.com/
	- https://yyy.com/
	- https://hhh.com/

	EXAMPLE 2:
	CONTEXT:

	Manusia diciptakan untuk beribadah
	--source: https://a.com/

	QUESTION:
	Untuk apa setan diciptakan?

	ANSWER:
	Maaf saya belum mempunyai informasi yang cukup untuk menjawab pertanyaan tersebut.

	---
	**REFERENCES**
	- none

	Provide answer as same as the example format. BEGIN!
		`

	knowledgeList, err := srv.embeddedKnowledgeService.RetriveKnowledgeBySearchQuery(question)
	if err != nil {
		return "", err
	}

	var allKnowledgeTextString string
	for i := 0; i < len(knowledgeList); i++ {
		knowledgeText := knowledgeList[i].Title + "\n" + knowledgeList[i].Content + "\n\n"
		allKnowledgeTextString = allKnowledgeTextString + knowledgeText
	}

	userPrompt := `
		CONTEXT:
		` + allKnowledgeTextString + `
		QUESTION:
		` + question + `
		`

	result, err := srv.aiChatThirdparty.Prompt(systemPrompt, userPrompt, 0.2, 1000)
	if err != nil {
		return "", err
	}

	// re := regexp.MustCompile(`(\d+)-(\d+)`)
	// submatches := re.FindStringSubmatch("123-456")
	// fmt.Println(submatches) // [123-456 123 456]

	return result, nil
}
