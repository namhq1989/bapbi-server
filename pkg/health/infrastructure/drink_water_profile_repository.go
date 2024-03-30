package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/infrastructure/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DrinkWaterProfileRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewDrinkWaterProfileRepository(db *mongo.Database) DrinkWaterProfileRepository {
	r := DrinkWaterProfileRepository{
		db:             db,
		collectionName: database.Tables.HealthDrinkWaterProfile,
	}
	r.ensureIndexes()
	return r
}

func (r DrinkWaterProfileRepository) ensureIndexes() {
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

func (r DrinkWaterProfileRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r DrinkWaterProfileRepository) CreateDrinkWaterProfile(ctx *appcontext.AppContext, profile domain.DrinkWaterProfile) error {
	// convert to mongodb model
	doc, err := model.DrinkWaterProfile{}.FromDomain(profile)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	if err != nil {
		ctx.Logger().Error("failed to insert drink water profile", err, appcontext.Fields{"profile": profile})
	}

	// respond
	return err
}

func (r DrinkWaterProfileRepository) UpdateDrinkWaterProfile(ctx *appcontext.AppContext, profile domain.DrinkWaterProfile) error {
	// convert to mongodb model
	doc, err := model.DrinkWaterProfile{}.FromDomain(profile)
	if err != nil {
		return err
	}

	_, err = r.collection().UpdateByID(ctx.Context(), doc.ID, bson.M{"$set": doc})
	if err != nil {
		ctx.Logger().Error("failed to update drink water profile", err, appcontext.Fields{"profile": profile})
	}

	// respond
	return err
}
