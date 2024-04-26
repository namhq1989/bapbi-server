package infrastructure

import (
	"context"
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
