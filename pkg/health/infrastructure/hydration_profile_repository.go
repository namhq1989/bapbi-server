package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/infrastructure/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HydrationProfileRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewHydrationProfileRepository(db *mongo.Database) HydrationProfileRepository {
	r := HydrationProfileRepository{
		db:             db,
		collectionName: database.Tables.HealthHydrationProfile,
	}
	r.ensureIndexes()
	return r
}

func (r HydrationProfileRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "userId", Value: 1}, {Key: "isEnabled", Value: 1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r HydrationProfileRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r HydrationProfileRepository) FindHydrationProfileByUserID(ctx *appcontext.AppContext, userID string) (*domain.HydrationProfile, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.HydrationProfile
	if err = r.collection().FindOne(ctx.Context(), bson.M{
		"userId": uid,
	}).Decode(&doc); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	// respond
	result := doc.ToDomain()
	return &result, nil
}

func (r HydrationProfileRepository) CreateHydrationProfile(ctx *appcontext.AppContext, profile domain.HydrationProfile) error {
	// convert to mongodb model
	doc, err := model.HydrationProfile{}.FromDomain(profile)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r HydrationProfileRepository) UpdateHydrationProfile(ctx *appcontext.AppContext, profile domain.HydrationProfile) error {
	// convert to mongodb model
	doc, err := model.HydrationProfile{}.FromDomain(profile)
	if err != nil {
		return err
	}

	_, err = r.collection().UpdateByID(ctx.Context(), doc.ID, bson.M{"$set": doc})
	return err
}
