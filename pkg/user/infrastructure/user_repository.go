package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewUserRepository(db *mongo.Database) UserRepository {
	r := UserRepository{
		db:             db,
		collectionName: database.Tables.User,
	}
	r.ensureIndexes()
	return r
}

func (r UserRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "email", Value: 1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r UserRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}
