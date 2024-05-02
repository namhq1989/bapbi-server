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
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserTermExerciseRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewUserTermExerciseRepository(db *mongo.Database) UserTermExerciseRepository {
	r := UserTermExerciseRepository{
		db:             db,
		collectionName: database.Tables.LanguageUserTermExercise,
	}
	r.ensureIndexes()
	return r
}

func (r UserTermExerciseRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "userId", Value: 1}, {Key: "termId", Value: 1}, {Key: "status", Value: 1}, {Key: "createdAt", Value: -1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r UserTermExerciseRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r UserTermExerciseRepository) CreateUserTermExercise(ctx *appcontext.AppContext, exercise domain.UserTermExercise) error {
	doc, err := model.UserTermExercise{}.FromDomain(exercise)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r UserTermExerciseRepository) UpdateUserTermExercise(ctx *appcontext.AppContext, exercise domain.UserTermExercise) error {
	doc, err := model.UserTermExercise{}.FromDomain(exercise)
	if err != nil {
		return err
	}

	_, err = r.collection().UpdateByID(ctx.Context(), doc.ID, bson.M{"$set": doc})
	return err
}

func (r UserTermExerciseRepository) FindByExerciseID(ctx *appcontext.AppContext, exerciseID string) (*domain.UserTermExercise, error) {
	id, err := database.ObjectIDFromString(exerciseID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.UserTermExercise
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

func (r UserTermExerciseRepository) FindByUserIDAndTermID(ctx *appcontext.AppContext, userID, termID string) (*domain.UserTermExercise, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	tid, err := database.ObjectIDFromString(termID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.UserTermExercise
	if err = r.collection().FindOne(ctx.Context(), bson.M{
		"userId": uid,
		"termId": tid,
	}).Decode(&doc); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	// respond
	result := doc.ToDomain()
	return &result, nil
}

func (r UserTermExerciseRepository) IsExerciseCreated(ctx *appcontext.AppContext, userID, termID string) (bool, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return false, apperrors.Common.InvalidID
	}

	tid, err := database.ObjectIDFromString(termID)
	if err != nil {
		return false, apperrors.Common.InvalidID
	}

	total, err := r.collection().CountDocuments(ctx.Context(), bson.M{
		"userId": uid,
		"termId": tid,
	})
	return total > 0, err
}
