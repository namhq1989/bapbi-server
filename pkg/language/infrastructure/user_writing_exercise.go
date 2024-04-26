package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"

	"github.com/namhq1989/bapbi-server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserWritingExerciseRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewUserWritingExerciseRepository(db *mongo.Database) UserWritingExerciseRepository {
	r := UserWritingExerciseRepository{
		db:             db,
		collectionName: database.Tables.LanguageUserWritingExercise,
	}
	r.ensureIndexes()
	return r
}

func (r UserWritingExerciseRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "userId", Value: 1}, {Key: "exerciseId", Value: 1}, {Key: "status", Value: 1}, {Key: "createdAt", Value: -1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r UserWritingExerciseRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r UserWritingExerciseRepository) CreateUserWritingExercise(ctx *appcontext.AppContext, exercise domain.UserWritingExercise) error {
	return nil
}
