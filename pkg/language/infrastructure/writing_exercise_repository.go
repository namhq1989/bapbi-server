package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure/model"

	"github.com/namhq1989/bapbi-server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WritingExerciseRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewWritingExerciseRepository(db *mongo.Database) WritingExerciseRepository {
	r := WritingExerciseRepository{
		db:             db,
		collectionName: database.Tables.LanguageWritingExercise,
	}
	r.ensureIndexes()
	return r
}

func (r WritingExerciseRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "language", Value: 1}, {Key: "type", Value: 1}, {Key: "level", Value: 1}, {Key: "createdAt", Value: -1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r WritingExerciseRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r WritingExerciseRepository) FindWritingExercises(ctx *appcontext.AppContext, filter domain.WritingExerciseFilter) ([]domain.WritingExercise, error) {
	var (
		condition = bson.M{"language": filter.Language.String()}
		result    = make([]domain.WritingExercise, 0)
	)

	if !filter.Time.IsZero() {
		condition["createdAt"] = bson.M{"$lt": filter.Time}
	}
	if filter.Level.IsValid() {
		condition["level"] = filter.Level.String()
	}

	// TODO: add pipeline to query data

	// find
	cursor, err := r.collection().Find(ctx.Context(), condition, &options.FindOptions{
		Sort: bson.M{"createdAt": -1},
	}, &options.FindOptions{
		Limit: &filter.Limit,
	})
	if err != nil {
		return result, err
	}
	// parse
	defer func() { _ = cursor.Close(ctx.Context()) }()

	// parse
	var docs []model.WritingExercise
	if err = cursor.All(ctx.Context(), &docs); err != nil {
		return result, err
	}

	// map data
	for _, doc := range docs {
		result = append(result, doc.ToDomain())
	}
	return result, nil
}
