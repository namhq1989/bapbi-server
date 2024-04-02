package infrastructure

import (
	"context"
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

type WaterIntakeLogRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewWaterIntakeLogRepository(db *mongo.Database) WaterIntakeLogRepository {
	r := WaterIntakeLogRepository{
		db:             db,
		collectionName: database.Tables.HealthWaterIntakeLog,
	}
	r.ensureIndexes()
	return r
}

func (r WaterIntakeLogRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "userId", Value: 1}, {Key: "intakeAt", Value: -1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r WaterIntakeLogRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r WaterIntakeLogRepository) CreateWaterIntakeLog(ctx *appcontext.AppContext, log domain.WaterIntakeLog) error {
	// convert to mongodb model
	doc, err := model.WaterIntakeLog{}.FromDomain(log)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r WaterIntakeLogRepository) FindWaterIntakeLogsByUserID(ctx *appcontext.AppContext, userID string, filter domain.WaterIntakeLogFilter) ([]domain.WaterIntakeLog, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	var (
		condition = bson.M{"userId": uid}
		result    = make([]domain.WaterIntakeLog, 0)
	)

	if !filter.From.IsZero() {
		condition["intakeAt"] = bson.M{"$gte": filter.From}
	}
	if !filter.To.IsZero() {
		condition["intakeAt"] = bson.M{"$lte": filter.To}
	}

	// find
	cursor, err := r.collection().Find(ctx.Context(), condition, &options.FindOptions{
		Sort: bson.M{"createdAt": -1},
	})
	if err != nil {
		return result, err
	}
	// parse
	defer func() { _ = cursor.Close(ctx.Context()) }()

	// parse
	var docs []model.WaterIntakeLog
	if err = cursor.All(ctx.Context(), &docs); err != nil {
		return result, err
	}

	// map data
	for _, doc := range docs {
		result = append(result, doc.ToDomain())
	}
	return result, nil
}
