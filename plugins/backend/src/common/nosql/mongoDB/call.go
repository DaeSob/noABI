package mongoDB

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
type TCredential struct {
	Uri            string
	DataBase       string
	Username       string
	Password       string
	ReplicaSetName string
}

type TOptions struct {
	MaxPoolSize     uint64
	MinPoolSize     uint64
	ConnectTimeout  time.Duration
	MaxConnIdleTime time.Duration
}

type TMongo struct {
	ctx        context.Context
	Session    *mongo.Client
	DataBase   *mongo.Database
	Collection string
}

type TModel struct {
	Id        string    `bson:"_id,omitempty" json:"_id"`
	Source    bson.M    `bson:"_source" json:"_source"`
	Timestamp time.Time `bson:"_timestamp" json:"_timestamp"`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
func New(_cred TCredential, _opt TOptions) (tMongo *TMongo, err error) {
	defer func() {
		if e := recover(); e != nil {
			tMongo = nil
		}
		return
	}()

	tMongo = new(TMongo)
	tMongo.ctx = context.Background()
	ctx := tMongo.ctx

	fmt.Println(_cred.Username)
	fmt.Println(_cred.Password)
	fmt.Println(_cred.ReplicaSetName)

	clientOptions := options.Client().ApplyURI(_cred.Uri).SetAuth(options.Credential{
		//AuthSource: "",
		AuthMechanism: "SCRAM-SHA-256", //MONGODB-CR
		AuthSource:    "admin",
		Username:      _cred.Username,
		Password:      _cred.Password,
	}).SetReplicaSet(_cred.ReplicaSetName)

	if _opt.MaxPoolSize > 0 {
		clientOptions.SetMaxPoolSize(_opt.MaxPoolSize)
	}
	if _opt.MinPoolSize > 0 {
		clientOptions.SetMinPoolSize(_opt.MinPoolSize)
	}
	if _opt.MaxConnIdleTime > 0 {
		clientOptions.SetMaxConnIdleTime(_opt.MaxConnIdleTime)
	}
	if _opt.ConnectTimeout > 0 {
		clientOptions.SetConnectTimeout(_opt.ConnectTimeout)
	}
	// clientOptions.SetServerSelectionTimeout(5 * time.Second)

	tMongo.Session, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
		//return nil, err
	}

	err = tMongo.Session.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	tMongo.DataBase = tMongo.Session.Database(_cred.DataBase)

	return
}

func (qm *TMongo) Disconnect() error {
	if qm.Session == nil {
		return nil
	}
	return qm.Session.Disconnect(qm.ctx)
}

func (qm *TMongo) getCollection(_collection string) *mongo.Collection {
	return qm.DataBase.Collection(_collection)
}

func (qm *TMongo) Insert(_timeout time.Duration, _collection string, _id string, _source bson.M) (*mongo.InsertOneResult, error) {
	ctx := qm.ctx
	var cancel context.CancelFunc
	if _timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, _timeout)
		defer cancel()
	}

	c := qm.getCollection(_collection)
	insert := TModel{Id: _id, Source: _source, Timestamp: time.Now()}
	return c.InsertOne(ctx, insert, nil)
}

func (qm *TMongo) Upsert(_timeout time.Duration, _collection string, _id string, _source bson.M) (*mongo.UpdateResult, error) {
	ctx := qm.ctx
	var cancel context.CancelFunc
	if _timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, _timeout)
		defer cancel()
	}

	c := qm.getCollection(_collection)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": TModel{Source: _source, Timestamp: time.Now()}}
	return c.UpdateOne(ctx, filter, update, opts)
}

func (qm *TMongo) Find(_timeout time.Duration, _collection string, filter interface{}, opts ...*options.FindOptions) (interface{}, error) {
	ctx := qm.ctx
	var cancel context.CancelFunc
	if _timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, _timeout)
		defer cancel()
	}

	// var datas []bson.M
	var datas []TModel

	c := qm.getCollection(_collection)
	cursor, err := c.Find(ctx, filter, opts...)

	defer func() {
		if cursor != nil {
			cursor.Close(ctx)
		}
		datas = nil
	}()

	if err != nil {
		return nil, err
	}

	//결과를 Josn으로 Build 하기
	err = cursor.All(ctx, &datas)
	return datas, err
}

func (qm *TMongo) FindOne(_timeout time.Duration, _collection string, filter interface{}, opts ...*options.FindOneOptions) (interface{}, error) {
	ctx := qm.ctx
	var cancel context.CancelFunc
	if _timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, _timeout)
		defer cancel()
	}

	// var data bson.M
	// defer func() {
	// 	data = nil
	// }()
	var data TModel

	c := qm.getCollection(_collection)
	err := c.FindOne(ctx, filter, opts...).Decode(&data)
	return data, err
}

/*
KST, _ := time.LoadLocation("Asia/Seoul")
now, _ := time.Now()
startTime := time.Date(now.Year(), now.Month(), now.Date(), 0,0,0,0, KST).UTC() // 자정
endTime := time.Date(now.Year(), now.Month(), now.Date(), 23,59,59,0, KST).UTC() // 자정 직전
filter := bson.M{"intime": bson.M{"$gte": startTime, "$lte": endTime}}

cursor, err := collection.Find(ctx, filter)
*/
