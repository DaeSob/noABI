package query

import "cia/api/db/collection"

func MakeLogEventFindQuery(
	_filter interface{},
	_limit int64,
) TQueryOPRequest {
	return TQueryOPRequest{
		OP: OP_FIND.String(),
		Query: TDocdbQuery{
			Collection: collection.COLLECTION_LOG_EVENT.String(),
			Filter:     _filter,
			Options: TDocdbQueryOption{
				Limit: _limit,
			},
		},
	}
}

func MakeLogEventAggregateQuery(
	_pipeline []interface{},
) TQueryOPRequest {
	return TQueryOPRequest{
		OP: OP_AGGREGATE.String(),
		Query: TDocdbQuery{
			Collection: collection.COLLECTION_LOG_EVENT.String(),
			Pipeline:   _pipeline,
		},
	}
}

func MakeFindAllQuery(
	_collectionName string,
) TQueryOPRequest {
	return TQueryOPRequest{
		OP: OP_FIND.String(),
		Query: TDocdbQuery{
			Collection: _collectionName,
		},
	}
}
