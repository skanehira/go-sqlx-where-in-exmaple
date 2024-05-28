package main

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Output struct {
	Name   string `db:"name"`
	Age    int    `db:"age"`
	Skills string `db:"skills"`
	Likes  string `db:"likes"`
}

type Condition struct {
	Name   *string
	Age    *int
	Skills []string
	Likes  []string
}

func ptr[T any](v T) *T {
	return &v
}

func main() {
	db := sqlx.MustOpen("sqlite3", ":memory:")

	db.MustExec(`
	CREATE TABLE IF NOT EXISTS user (
		name TEXT,
		age INTEGER,
		skills TEXT,
		likes TEXT
	);
	`)

	db.MustExec(`
INSERT INTO
  user (name, age, skills, likes)
VALUES
  ('John', 30, 'Go', 'Apple'),
  ('Jane', 25, 'Python', 'Orange'),
  ('Doe', 35, 'Rust', 'Banana');
	`)

	cond := Condition{
		Skills: []string{
			"Go",
			"Rust",
		},
		Likes: []string{
			"Banana",
			"Apple",
		},
	}

	baseQuery := `SELECT * FROM user WHERE`

	conds := []string{}
	params := map[string]any{}

	if cond.Name != nil {
		conds = append(conds, "name = :name")
		params["name"] = *cond.Name
	}

	if cond.Age != nil {
		conds = append(conds, "age = :age")
		params["age"] = *cond.Age
	}

	if len(cond.Skills) > 0 {
		conds = append(conds, "skills IN (:skills)")
		params["skills"] = cond.Skills
	}

	if len(cond.Likes) > 0 {
		conds = append(conds, "likes IN (:likes)")
		params["likes"] = cond.Likes
	}

	query, args, err := sqlx.Named(baseQuery+" "+strings.Join(conds, " AND "), params)
	if err != nil {
		panic(err)
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		panic(err)
	}

	query = db.Rebind(query)

	fmt.Println(query, args)

	var out []Output
	if err := db.Select(&out, query, args...); err != nil {
		panic(err)
	}

	fmt.Printf("%#+v\n", out)
}
