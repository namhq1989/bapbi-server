package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserSearchRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewUserSearchRepository(db *mongo.Database) UserSearchRepository {
	r := UserSearchRepository{
		db:             db,
		collectionName: database.Tables.LanguageUserSearchHistory,
	}
	r.ensureIndexes()
	return r
}

func (r UserSearchRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "user", Value: 1}, {Key: "term", Value: 1}, {Key: "createdAt", Value: -1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r UserSearchRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r UserSearchRepository) CreateUserSearchHistory(ctx *appcontext.AppContext, history domain.UserSearchHistory) error {
	// convert to mongodb model
	doc, err := model.UserSearchHistory{}.FromDomain(history)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}
