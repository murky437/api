package model

import "database/sql"

type ProjectBasic struct {
	Id        int
	Title     string
	Slug      string
	SortIndex int
}

type ProjectNotes struct {
	Notes string
}

func CreateProject(db *sql.DB, project ProjectBasic, userId int) (ProjectBasic, error) {
	res, err := db.Exec(`
		INSERT INTO project (title, slug, user_id, sort_index)
		VALUES (?, ?, ?, (SELECT COUNT(*) FROM project WHERE user_id = ?))
	`, project.Title, project.Slug, userId, userId)
	if err != nil {
		return project, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return project, err
	}
	project.Id = int(id)

	return project, nil
}

func UpdateProject(db *sql.DB, currentSlug string, project ProjectBasic) error {
	_, err := db.Exec(`
		UPDATE project 
		SET title = :title, slug = :slug, updated_at = STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')
		WHERE slug = :currentSlug
	`,
		sql.Named("title", project.Title),
		sql.Named("slug", project.Slug),
		sql.Named("currentSlug", currentSlug),
	)
	return err
}

func UpdateProjectNotes(db *sql.DB, currentSlug string, project ProjectNotes) error {
	_, err := db.Exec(`
		UPDATE project 
		SET notes = :notes, updated_at = STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')
		WHERE slug = :currentSlug
	`,
		sql.Named("notes", project.Notes),
		sql.Named("currentSlug", currentSlug),
	)
	return err
}

type SlugCheckOptions struct {
	CurrentSlug string
}

func IsProjectSlugUnique(db *sql.DB, slug string, opts SlugCheckOptions) (bool, error) {
	var (
		query string
		args  []any
	)

	if opts.CurrentSlug != "" {
		query = "SELECT COUNT(1) FROM project WHERE slug = ? AND slug != ?"
		args = []any{slug, opts.CurrentSlug}
	} else {
		query = "SELECT COUNT(1) FROM project WHERE slug = ?"
		args = []any{slug}
	}

	var count int
	err := db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func GetProjectsByUserId(db *sql.DB, userId int) ([]ProjectBasic, error) {
	rows, err := db.Query(`
		SELECT p.id, p.title, p.slug, p.sort_index
		FROM project p
		WHERE p.user_id = ?
		ORDER BY p.sort_index
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []ProjectBasic
	for rows.Next() {
		var p ProjectBasic
		if err := rows.Scan(&p.Id, &p.Title, &p.Slug, &p.SortIndex); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func GetProjectByUserIdAndSlug(db *sql.DB, userId int, slug string) (ProjectBasic, error) {
	var p ProjectBasic
	err := db.QueryRow(`
		SELECT p.id, p.title, p.slug, p.sort_index
		FROM project p
		WHERE p.user_id = ? AND p.slug = ?
	`, userId, slug).Scan(&p.Id, &p.Title, &p.Slug, &p.SortIndex)

	return p, err
}

func GetProjectNotesByUserIdAndSlug(db *sql.DB, userId int, slug string) (ProjectNotes, error) {
	var p ProjectNotes
	err := db.QueryRow(`
		SELECT p.notes
		FROM project p
		WHERE p.user_id = ? AND p.slug = ?
	`, userId, slug).Scan(&p.Notes)

	return p, err
}

func DeleteProjectByUserIdAndSlug(db *sql.DB, userId int, slug string) (int64, error) {
	var delIndex int
	err := db.QueryRow(`
		SELECT
		    sort_index
		FROM
		    project
		WHERE
		    user_id = ?
		  	AND slug = ?
	`, userId, slug).Scan(&delIndex)
	if err != nil {
		return 0, err
	}

	res, err := db.Exec(`
		DELETE FROM project
		WHERE 
		    user_id = ?
		  	AND slug = ?;
	`, userId, slug)
	if err != nil {
		return 0, err
	}

	_, err = db.Exec(`
		UPDATE
		    project
		SET
		    sort_index = sort_index - 1,
		    updated_at = STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')
		WHERE 
		    user_id = ?
			AND sort_index > ?;
	`, userId, delIndex)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func UpdateProjectSortIndex(db *sql.DB, userId int, slug string, newIndex int) error {
	var oldIndex int
	err := db.QueryRow(`
		SELECT
		    sort_index
		FROM
		    project
		WHERE
		    user_id = ?
		    AND slug = ?
	`, userId, slug).Scan(&oldIndex)
	if err != nil {
		return err
	}

	if newIndex > oldIndex {
		_, err = db.Exec(`
			UPDATE
			    project
			SET
			    sort_index = sort_index - 1,
			    updated_at = STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')
			WHERE
			    user_id = :userId
			    AND sort_index > :oldIndex
			  	AND sort_index <= :newIndex;
		`,
			sql.Named("userId", userId),
			sql.Named("oldIndex", oldIndex),
			sql.Named("newIndex", newIndex),
		)
		if err != nil {
			return err
		}
	}

	if newIndex < oldIndex {
		_, err = db.Exec(`
			UPDATE
			    project
			SET
			    sort_index = sort_index + 1,
			    updated_at = STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')
			WHERE
			    user_id = :userId
			    AND sort_index >= :newIndex
			  	AND sort_index < :oldIndex;
		`,
			sql.Named("userId", userId),
			sql.Named("oldIndex", oldIndex),
			sql.Named("newIndex", newIndex),
		)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec(`
		UPDATE
			project
		SET
			sort_index = :newIndex,
			updated_at = STRFTIME('%Y-%m-%dT%H:%M:%fZ','NOW')
		WHERE
			user_id = :userId
			AND slug = :slug;
	`,
		sql.Named("newIndex", newIndex),
		sql.Named("userId", userId),
		sql.Named("slug", slug),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetProjectCount(db *sql.DB, userId int) (int, error) {
	var count int
	err := db.QueryRow(`
		SELECT
		    COUNT(*)
		FROM
		    project
		WHERE
		    user_id = ?
	`, userId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
