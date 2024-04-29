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

type UserVocabularyExerciseRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewUserVocabularyExerciseRepository(db *mongo.Database) UserVocabularyExerciseRepository {
	r := UserVocabularyExerciseRepository{
		db:             db,
		collectionName: database.Tables.LanguageUserVocabularyExercise,
	}
	r.ensureIndexes()
	return r
}

func (r UserVocabularyExerciseRepository) ensureIndexes() {
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

func (r UserVocabularyExerciseRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r UserVocabularyExerciseRepository) CreateUserVocabularyExercise(ctx *appcontext.AppContext, exercise domain.UserVocabularyExercise) error {
	doc, err := model.UserVocabularyExercise{}.FromDomain(exercise)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r UserVocabularyExerciseRepository) UpdateUserVocabularyExercise(ctx *appcontext.AppContext, exercise domain.UserVocabularyExercise) error {
	doc, err := model.UserVocabularyExercise{}.FromDomain(exercise)
	if err != nil {
		return err
	}

	_, err = r.collection().UpdateByID(ctx.Context(), doc.ID, bson.M{"$set": doc})
	return err
}

func (r UserVocabularyExerciseRepository) FindByUserIDAndTermID(ctx *appcontext.AppContext, userID, termID string) (*domain.UserVocabularyExercise, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	tid, err := database.ObjectIDFromString(termID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.UserVocabularyExercise
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

func (r UserVocabularyExerciseRepository) IsExerciseCreated(ctx *appcontext.AppContext, userID, termID string) (bool, error) {
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
