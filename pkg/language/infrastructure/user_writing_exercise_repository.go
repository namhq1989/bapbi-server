package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure/mapping"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure/model"
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
	doc, err := model.UserWritingExercise{}.FromDomain(exercise)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r UserWritingExerciseRepository) UpdateUserWritingExercise(ctx *appcontext.AppContext, exercise domain.UserWritingExercise) error {
	doc, err := model.UserWritingExercise{}.FromDomain(exercise)
	if err != nil {
		return err
	}

	_, err = r.collection().UpdateByID(ctx.Context(), doc.ID, bson.M{"$set": doc})
	return err
}

func (r UserWritingExerciseRepository) FindByUserIDAndExerciseID(ctx *appcontext.AppContext, userID, exerciseID string) (*domain.UserWritingExercise, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	eid, err := database.ObjectIDFromString(exerciseID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.UserWritingExercise
	if err = r.collection().FindOne(ctx.Context(), bson.M{
		"userId":     uid,
		"exerciseId": eid,
	}).Decode(&doc); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	// respond
	result := doc.ToDomain()
	return &result, nil
}

func (r UserWritingExerciseRepository) IsExerciseCreated(ctx *appcontext.AppContext, userID, exerciseID string) (bool, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return false, apperrors.Common.InvalidID
	}

	eid, err := database.ObjectIDFromString(exerciseID)
	if err != nil {
		return false, apperrors.Common.InvalidID
	}

	total, err := r.collection().CountDocuments(ctx.Context(), bson.M{
		"userId":     uid,
		"exerciseId": eid,
	})
	return total > 0, err
}

func (r UserWritingExerciseRepository) FindUserWritingExercises(ctx *appcontext.AppContext, filter domain.UserWritingExerciseFilter) ([]domain.WritingExerciseDatabaseQuery, error) {
	result := make([]domain.WritingExerciseDatabaseQuery, 0)

	uid, err := database.ObjectIDFromString(filter.UserID)
	if err != nil {
		return result, apperrors.User.InvalidUserID
	}

	matchCondition := bson.D{{Key: "userId", Value: uid}, {Key: "language", Value: filter.Language.String()}}

	if !filter.Time.IsZero() {
		matchCondition = append(matchCondition, bson.E{Key: "createdAt", Value: bson.M{"$lt": filter.Time}})
	}
	if filter.Status != "" {
		matchCondition = append(matchCondition, bson.E{Key: "status", Value: filter.Status})
	}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.M{
			"from": database.Tables.LanguageWritingExercise,
			"let":  bson.M{"exerciseId": "$exerciseId"},
			"pipeline": bson.A{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{
							"$eq": bson.A{"$_id", "$$exerciseId"},
						},
					},
				},
			},
			"as": "writingExercise",
		}},
	}

	projectState := bson.D{{Key: "$project", Value: bson.M{
		"_id":         "$writingExercise._id",
		"language":    "$writingExercise.language",
		"type":        "$writingExercise.type",
		"level":       "$writingExercise.level",
		"topic":       "$writingExercise.topic",
		"question":    "$writingExercise.question",
		"vocabulary":  "$writingExercise.vocabulary",
		"minWords":    "$writingExercise.minWords",
		"data":        "$writingExercise.data",
		"status":      1,
		"createdAt":   1,
		"completedAt": 1,
	}}}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchCondition}},
		lookupStage,
		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$writingExercise",
			"includeArrayIndex":          "0",
			"preserveNullAndEmptyArrays": true,
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
