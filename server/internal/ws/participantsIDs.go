package ws

import (
	"sync"

	snapws "github.com/Atheer-Ganayem/SnapWS"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SafeIDs struct {
	IDs []bson.ObjectID
	Mu  sync.RWMutex
}

func NewParticipantIDs(ids []bson.ObjectID) *SafeIDs {
	return &SafeIDs{
		IDs: ids,
	}
}

func (sIDs *SafeIDs) Append(id bson.ObjectID) {
	sIDs.Mu.Lock()
	defer sIDs.Mu.Unlock()

	sIDs.IDs = append(sIDs.IDs, id)
}

func AppendParticipant(conn *snapws.ManagedConn[bson.ObjectID], id bson.ObjectID) {
	val, ok := conn.MetaData.Load("participantsIDs")
	if !ok {
		return
	}

	sIDs, ok := val.(*SafeIDs)
	if !ok {
		return
	}

	sIDs.Append(id)
}
