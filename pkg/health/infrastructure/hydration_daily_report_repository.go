package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/health/domain"
	"github.com/namhq1989/bapbi-server/pkg/health/infrastructure/model"

	"github.com/namhq1989/bapbi-server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HydrationDailyReportRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewHydrationDailyReportRepository(db *mongo.Database) HydrationDailyReportRepository {
	r := HydrationDailyReportRepository{
		db:             db,
		collectionName: database.Tables.HealthHydrationDailyReport,
	}
	r.ensureIndexes()
	return r
}

func (r HydrationDailyReportRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "userId", Value: 1}, {Key: "date", Value: -1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r HydrationDailyReportRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r HydrationDailyReportRepository) FindHydrationDailyReportByUserID(ctx *appcontext.AppContext, userID string, date time.Time) (*domain.HydrationDailyReport, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.HydrationDailyReport
	if err = r.collection().FindOne(ctx.Context(), bson.M{
		"userId": uid,
		"date":   date,
	}).Decode(&doc); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	// respond
	result := doc.ToDomain()
	return &result, nil
}

func (r HydrationDailyReportRepository) CreateHydrationDailyReport(ctx *appcontext.AppContext, report domain.HydrationDailyReport) error {
	// convert to mongodb model
	doc, err := model.HydrationDailyReport{}.FromDomain(report)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r HydrationDailyReportRepository) UpdateHydrationDailyReport(ctx *appcontext.AppContext, report domain.HydrationDailyReport) error {
	// convert to mongodb model
	doc, err := model.HydrationDailyReport{}.FromDomain(report)
	if err != nil {
		return err
	}

	_, err = r.collection().UpdateByID(ctx.Context(), doc.ID, bson.M{"$set": doc})
	return err
}
