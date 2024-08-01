package app_services

import (
	infra_services "bos_personal_ai/services/infra"
	"bos_personal_ai/thirdparties"
	"regexp"
	"strings"
)

type RAGOutput struct {
	RawAnswer     string   `json:"raw_answer"`
	Answer        string   `json:"answer"`
	RawReferences string   `json:"raw_references"`
	References    []string `json:"references"`
}

type RAGInterface interface {
	AskQuestion(modelProvider string, question string) (RAGOutput, error)
}

type RAGImpl struct {
	embeddingThirdparty      thirdparties.EmbeddingThirdPartyInterface
	aiChatService            infra_services.AIChatServiceInterface
	embeddedKnowledgeService EmbeddedKnowledgeServiceInterface
}

func NewRAG(
	embeddingThirdparty thirdparties.EmbeddingThirdPartyInterface,
	aiChatService infra_services.AIChatServiceInterface,
	embeddedKnowledgeService EmbeddedKnowledgeServiceInterface,
) RAGInterface {
	return &RAGImpl{
		embeddingThirdparty:      embeddingThirdparty,
		aiChatService:            aiChatService,
		embeddedKnowledgeService: embeddedKnowledgeService,
	}
}

func (srv *RAGImpl) AskQuestion(modelProvider string, question string) (RAGOutput, error) {

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

	Provide answer with format as same as the example format. BEGIN!
		`

	knowledgeList, err := srv.embeddedKnowledgeService.RetriveKnowledgeBySearchQuery(question)
	if err != nil {
		return RAGOutput{}, err
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

	rawAnswer, err := srv.aiChatService.Prompt(modelProvider, systemPrompt, userPrompt, 0.2, 1000)
	if err != nil {
		return RAGOutput{}, err
	}

	// references: https://stackoverflow.com/a/39102969 (how to get string on submatch regex golang)
	re := regexp.MustCompile(`(?mi)(ANSWER|):\s([^+]+)\s---\s\*\*REFERENCES\*\*\s([^+]+)`)
	matches := re.FindAllStringSubmatch(rawAnswer, -1)

	if len(matches) == 0 {
		return RAGOutput{
			RawAnswer: rawAnswer,
		}, nil
	}

	cleanAnswer := strings.TrimSpace(matches[0][2])
	rawReferences := strings.TrimSpace(matches[0][3])

	var references []string
	if len(rawReferences) > 0 {
		splitLinks := strings.Split(rawReferences, "\n")
		for _, link := range splitLinks {
			cleanLink := strings.TrimSpace(link)
			cleanLink = strings.TrimPrefix(cleanLink, "- ")
			references = append(references, cleanLink)
		}
	}

	output := RAGOutput{
		RawAnswer:     rawAnswer,
		Answer:        cleanAnswer,
		RawReferences: rawReferences,
		References:    references,
	}

	return output, nil
}
