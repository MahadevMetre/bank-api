package models

import (
	"database/sql"
	"strings"
	"time"

	"bankapi/constants"
)

type Relations struct {
	Id              int64     `json:"id"`
	NomineeCode     string    `json:"nominee_code"`
	NomineeRelation string    `json:"nominee_relation"`
	ShouldShow      bool      `json:"should_show"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewRelations() *Relations {
	return &Relations{}
}

func FindAllActiveRelations(db *sql.DB) ([]Relations, error) {
	relations := make([]Relations, 0)
	rows, err := db.Query("SELECT id, nominee_code::integer, nominee_relation, should_show, created_at, updated_at FROM nominee_master WHERE should_show = true ORDER BY nominee_code::integer ASC")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoRelationsFound
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		rel := NewRelations()

		if err := rows.Scan(
			&rel.Id,
			&rel.NomineeCode,
			&rel.NomineeRelation,
			&rel.ShouldShow,
			&rel.CreatedAt,
			&rel.UpdatedAt,
		); err != nil {
			return nil, err
		}
		relations = append(relations, *rel)
	}

	return relations, nil
}

func CustomFindQuery(db *sql.DB, query string) ([]Relations, error) {
	rows, err := db.Query(query)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoRelationsFound
		}
		return nil, err
	}

	defer rows.Close()

	var relations []Relations

	for rows.Next() {
		rel := NewRelations()

		if err := rows.Scan(
			&rel.Id,
			&rel.NomineeCode,
			&rel.NomineeRelation,
			&rel.ShouldShow,
			&rel.CreatedAt,
			&rel.UpdatedAt,
		); err != nil {
			return nil, err
		}

		relations = append(relations, *rel)
	}

	return relations, nil
}
func FindRelationsClause(db *sql.DB, clause string, params ...interface{}) ([]Relations, error) {
	rows, err := db.Query("SELECT id, nominee_code, nominee_relation, should_show, created_at, updated_at FROM nominee_master WHERE "+clause, params...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoRelationsFound
		}
		return nil, err
	}

	defer rows.Close()

	var relations []Relations

	for rows.Next() {
		rel := NewRelations()

		if err := rows.Scan(
			&rel.Id,
			&rel.NomineeCode,
			&rel.NomineeRelation,
			&rel.ShouldShow,
			&rel.CreatedAt,
			&rel.UpdatedAt,
		); err != nil {
			return nil, err
		}

		relations = append(relations, *rel)
	}

	return relations, nil

}

func FindOneRelationClause(db *sql.DB, clause string, params ...interface{}) (*Relations, error) {
	rel := NewRelations()

	row := db.QueryRow("SELECT id, nominee_code, nominee_relation, should_show, created_at, updated_at FROM nominee_master WHERE "+clause, params...)

	if err := row.Scan(
		&rel.Id,
		&rel.NomineeCode,
		&rel.NomineeRelation,
		&rel.ShouldShow,
		&rel.CreatedAt,
		&rel.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrRelationNotFound
		}
		return nil, err
	}

	return rel, nil
}

func FindNomineeByNomineeRelation(db *sql.DB, relation string) (*Relations, error) {

	relationCondition := strings.ToUpper(relation)
	rel := NewRelations()
	row := db.QueryRow("SELECT id, nominee_code, nominee_relation, should_show, created_at, updated_at FROM nominee_master WHERE nominee_relation = $1 AND should_show = true", relationCondition)

	if err := row.Scan(
		&rel.Id,
		&rel.NomineeCode,
		&rel.NomineeRelation,
		&rel.ShouldShow,
		&rel.CreatedAt,
		&rel.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrRelationNotFound
		}
		return nil, err
	}

	return rel, nil
}

func FindNomineeByCode(db *sql.DB, code string) (*Relations, error) {
	rel := NewRelations()
	row := db.QueryRow("SELECT id, nominee_code, nominee_relation, should_show, created_at, updated_at FROM nominee_master WHERE nominee_code = $1 AND should_show = true", code)

	if err := row.Scan(
		&rel.Id,
		&rel.NomineeCode,
		&rel.NomineeRelation,
		&rel.ShouldShow,
		&rel.CreatedAt,
		&rel.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrRelationNotFound
		}
		return nil, err
	}

	return rel, nil
}
