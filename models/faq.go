package models

import (
	"bankapi/config"
	"database/sql"

	"github.com/google/uuid"
)

type FaqCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Faqs []Faq     `json:"faqs"`
}

type Faq struct {
	ID       uuid.UUID `json:"id"`
	Question string    `json:"question"`
	Answer   string    `json:"answer"`
	VideoUrl string    `json:"video_url"`
}

// get all category and its data
func GetAllCategoriesWithFaqs(appOnly, webOnly bool) ([]FaqCategory, error) {
	db := config.GetDB()

	query := "SELECT c.id, c.name, f.id, f.question, f.answer, f.video_url FROM faq_categories c LEFT JOIN faqs f ON c.id = f.category_id"
	var args []interface{}

	if appOnly {
		query += " AND f.app_only = true"
	} else if webOnly {
		query += " AND f.website_only = true"
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make(map[uuid.UUID]*FaqCategory)
	for rows.Next() {
		var categoryID, faqID uuid.UUID
		var categoryName string
		var question, answer, videoUrl sql.NullString
		if err := rows.Scan(&categoryID, &categoryName, &faqID, &question, &answer, &videoUrl); err != nil {
			return nil, err
		}

		if _, exists := categories[categoryID]; !exists {
			categories[categoryID] = &FaqCategory{ID: categoryID, Name: categoryName, Faqs: []Faq{}}
		}

		if faqID != uuid.Nil {
			categories[categoryID].Faqs = append(categories[categoryID].Faqs, Faq{
				ID:       faqID,
				Question: question.String,
				Answer:   answer.String,
				VideoUrl: videoUrl.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]FaqCategory, 0, len(categories))
	for _, category := range categories {
		result = append(result, *category)
	}

	return result, nil
}
