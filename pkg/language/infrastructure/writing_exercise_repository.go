package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure/mapping"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure/model"
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

func (r WritingExerciseRepository) CreateWritingExercise(ctx *appcontext.AppContext, exercise domain.WritingExercise) error {
	doc, err := model.WritingExercise{}.FromDomain(exercise)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r WritingExerciseRepository) FindByID(ctx *appcontext.AppContext, exerciseID string) (*domain.WritingExercise, error) {
	id, err := database.ObjectIDFromString(exerciseID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.WritingExercise
	if err = r.collection().FindOne(ctx.Context(), bson.M{
		"_id": id,
	}).Decode(&doc); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	// respond
	result := doc.ToDomain()
	return &result, nil
}

func (r WritingExerciseRepository) FindWritingExercises(ctx *appcontext.AppContext, filter domain.WritingExerciseFilter) ([]domain.WritingExerciseDatabaseQuery, error) {
	result := make([]domain.WritingExerciseDatabaseQuery, 0)

	uid, err := database.ObjectIDFromString(filter.UserID)
	if err != nil {
		return result, apperrors.User.InvalidUserID
	}

	matchCondition := bson.D{{Key: "language", Value: filter.Language.String()}}

	if !filter.Time.IsZero() {
		matchCondition = append(matchCondition, bson.E{Key: "createdAt", Value: bson.M{"$lt": filter.Time}})
	}
	if filter.Level.IsValid() {
		matchCondition = append(matchCondition, bson.E{Key: "level", Value: filter.Level.String()})
	}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.M{
			"from": database.Tables.LanguageUserWritingExercise,
			"let":  bson.M{"exerciseId": "$_id"},
			"pipeline": bson.A{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": bson.A{
								bson.M{"$eq": bson.A{"$userId", uid}},
								bson.M{"$eq": bson.A{"$exerciseId", "$$exerciseId"}},
							},
						},
					},
				},
				bson.M{"$project": bson.M{
					"_id":    0,
					"status": 1,
				}},
			},
			"as": "userWritingExercise",
		}},
	}

	projectState := bson.D{{Key: "$project", Value: bson.M{
		"_id":        1,
		"language":   1,
		"type":       1,
		"level":      1,
		"topic":      1,
		"question":   1,
		"vocabulary": 1,
		"minWords":   1,
		"data":       1,
		"createdAt":  1,
		"status":     1,
	}}}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchCondition}},
		lookupStage,
		bson.D{{Key: "$addFields", Value: bson.M{
			"status": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$userWritingExercise.status", 0}}, ""}},
		}}},
		projectState,
		bson.D{{Key: "$sort", Value: bson.M{"createdAt": -1}}},
		bson.D{{Key: "$limit", Value: filter.Limit}},
	}

	// find
	cursor, err := r.collection().Aggregate(ctx.Context(), pipeline)
	if err != nil {
		return result, err
	}
	// parse
	defer func() { _ = cursor.Close(ctx.Context()) }()

	// parse
	var docs []mapping.WritingExerciseDatabaseQuery
	if err = cursor.All(ctx.Context(), &docs); err != nil {
		return result, err
	}

	// map data
	for _, doc := range docs {
		result = append(result, doc.ToDomain())
	}
	return result, nil
}
