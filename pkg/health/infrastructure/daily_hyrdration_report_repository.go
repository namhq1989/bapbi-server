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

type DailyHydrationReportRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewDailyHydrationReportRepository(db *mongo.Database) DailyHydrationReportRepository {
	r := DailyHydrationReportRepository{
		db:             db,
		collectionName: database.Tables.HealthDailyHydrationReport,
	}
	r.ensureIndexes()
	return r
}

func (r DailyHydrationReportRepository) ensureIndexes() {
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

func (r DailyHydrationReportRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r DailyHydrationReportRepository) FindDailyHydrationReportByUserID(ctx *appcontext.AppContext, userID string, date time.Time) (*domain.DailyHydrationReport, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.DailyHydrationReport
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

func (r DailyHydrationReportRepository) CreateDailyHydrationReport(ctx *appcontext.AppContext, report domain.DailyHydrationReport) error {
	// convert to mongodb model
	doc, err := model.DailyHydrationReport{}.FromDomain(report)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r DailyHydrationReportRepository) UpdateDailyHydrationReport(ctx *appcontext.AppContext, report domain.DailyHydrationReport) error {
	// convert to mongodb model
	doc, err := model.DailyHydrationReport{}.FromDomain(report)
	if err != nil {
		return err
	}

	_, err = r.collection().UpdateByID(ctx.Context(), doc.ID, bson.M{"$set": doc})
	return err
}
