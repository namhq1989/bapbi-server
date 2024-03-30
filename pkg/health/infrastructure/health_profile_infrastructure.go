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

type HealthProfileRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewHealthProfileRepository(db *mongo.Database) HealthProfileRepository {
	r := HealthProfileRepository{
		db:             db,
		collectionName: database.Tables.HealthProfile,
	}
	r.ensureIndexes()
	return r
}

func (r HealthProfileRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "userId", Value: 1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r HealthProfileRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r HealthProfileRepository) CreateHealthProfile(ctx *appcontext.AppContext, profile domain.HealthProfile) error {
	// convert to mongodb model
	doc, err := model.HealthProfile{}.FromDomain(profile)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	if err != nil {
		ctx.Logger().Error("failed to insert health profile", err, appcontext.Fields{"profile": profile})
	}

	// respond
	return err
}
