package query

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type pendingSlot struct {
	Id        string `bson:"_id" json:"_id,omitempty"`
	Offset    uint8  `bson:"offset" json:"offset,omitempty"`
	Timestamp uint64 `bson:"timestamp" json:"timestamp,omitempty"`
}

func FindOneNonceStatus(fromTimeStamp int64, nonce, vpaOwner string) (res pendingSlot, err error) {

	payload := QueryPayload{
		Op: "findone",
		Query: DocdbQuery{
			Collection: "message_slot",
			Filter: bson.M{
				"offset":                    bson.M{"$lt": 5},
				"timestamp":                 bson.M{"$gte": fromTimeStamp},
				"slots.vpaOwner":            strings.ToLower(vpaOwner),
				"slots.userOperation.nonce": nonce,
			},
		},
	}
	payload.Query.Options.ProjectEasy = []string{"_id", "offset", "timestamp"}

	retrieved, err := ExecQuery[pendingSlot](payload, true)

	return retrieved, err
}
