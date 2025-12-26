package query

import (
	"cia/api/db/query"
	"cia/common/utils"
)

func makeFindAllQuery(_collectionName string) string {
	return utils.InterfaceToJsonString(
		query.MakeFindAllQuery(
			_collectionName,
		),
		false,
	)
}

func FindAllQuery(_collectionName string) TFindQueryResponse {
	query := makeFindAllQuery(_collectionName)
	return findQueryOnDB(query)
}
