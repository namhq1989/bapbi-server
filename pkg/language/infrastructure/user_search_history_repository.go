package infrastructure

import (
	"context"
	"fmt"
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserSearchHistoryRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewUserSearchHistoryRepository(db *mongo.Database) UserSearchHistoryRepository {
	r := UserSearchHistoryRepository{
		db:             db,
		collectionName: database.Tables.LanguageUserSearchHistory,
	}
	r.ensureIndexes()
	return r
}

func (r UserSearchHistoryRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "userId", Value: 1}, {Key: "term", Value: 1}, {Key: "createdAt", Value: -1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r UserSearchHistoryRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r UserSearchHistoryRepository) CreateUserSearchHistory(ctx *appcontext.AppContext, history domain.UserSearchHistory) error {
	doc, err := model.UserSearchHistory{}.FromDomain(history)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r UserSearchHistoryRepository) CountTotalSearchedByTimeRange(ctx *appcontext.AppContext, userID string, start, end time.Time) (int64, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return 0, apperrors.Common.InvalidID
	}

	total, err := r.collection().CountDocuments(ctx.Context(), bson.M{
		"userId": uid,
		"createdAt": bson.M{
			"$gte": start,
			"$lte": end,
		},
	})
	return total, err
}
