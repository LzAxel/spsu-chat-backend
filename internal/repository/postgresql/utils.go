package postgresql

func FormatINClause(field string, count int) string {
	query := field + " IN ("

	for i := 0; i < count; i++ {
		query += "?"
		if i < count-1 {
			query += ","
		}
	}
	query += ")"

	return query
}
