package repository

import "github.com/huandu/go-sqlbuilder"

const sqlBuilderFlavor = sqlbuilder.PostgreSQL

func newSelectBuilder() *sqlbuilder.SelectBuilder {
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlBuilderFlavor)

	return sb
}

func newInsertBuilder() *sqlbuilder.InsertBuilder {
	ib := sqlbuilder.NewInsertBuilder()
	ib.SetFlavor(sqlBuilderFlavor)

	return ib
}

func newUpdateBuilder() *sqlbuilder.UpdateBuilder {
	ub := sqlbuilder.NewUpdateBuilder()
	ub.SetFlavor(sqlBuilderFlavor)

	return ub
}

func newDeleteBuilder() *sqlbuilder.DeleteBuilder {
	db := sqlbuilder.NewDeleteBuilder()
	db.SetFlavor(sqlBuilderFlavor)

	return db
}
