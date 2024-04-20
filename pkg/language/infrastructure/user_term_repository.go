package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	"github.com/namhq1989/bapbi-server/pkg/language/domain"
	"github.com/namhq1989/bapbi-server/pkg/language/infrastructure/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserTermRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewUserTermRepository(db *mongo.Database) UserTermRepository {
	r := UserTermRepository{
		db:             db,
		collectionName: database.Tables.LanguageUserTerm,
	}
	r.ensureIndexes()
	return r
}

func (r UserTermRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "userId", Value: 1}, {Key: "term", Value: 1}, {Key: "createdAt", Value: -1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r UserTermRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r UserTermRepository) FindUserTermByID(ctx *appcontext.AppContext, termID string) (*domain.UserTerm, error) {
	id, err := database.ObjectIDFromString(termID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.UserTerm
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

func (r UserTermRepository) AddUserTerm(ctx *appcontext.AppContext, term domain.UserTerm) error {
	doc, err := model.UserTerm{}.FromDomain(term)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r UserTermRepository) UpdateUserTerm(ctx *appcontext.AppContext, term domain.UserTerm) error {
	doc, err := model.UserTerm{}.FromDomain(term)
	if err != nil {
		return err
	}

	_, err = r.collection().UpdateByID(ctx.Context(), doc.ID, bson.M{"$set": doc})
	return err
}

func (r UserTermRepository) IsUserTermAdded(ctx *appcontext.AppContext, userID, term string) (bool, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return false, apperrors.Common.InvalidID
	}

	total, err := r.collection().CountDocuments(ctx.Context(), bson.M{
		"userId": uid,
		"term":   term,
	})
	return total > 0, err
}

func (r UserTermRepository) CountTotalTermAddedByTimeRange(ctx *appcontext.AppContext, userID string, start, end time.Time) (int64, error) {
	uid, err := database.ObjectIDFromString(userID)
	if err != nil {
		return 0, apperrors.Common.InvalidID
	}

	total, err := r.collection().CountDocuments(ctx.Context(), bson.M{
		"userId": uid,
		"createdAt": bson.M{
			"$gte": start,
			"$lte": end,
		},
	})
	return total, err
}

func (r UserTermRepository) FindUserTerms(ctx *appcontext.AppContext, filter domain.UserTermFilter) ([]domain.UserTerm, error) {
	uid, err := database.ObjectIDFromString(filter.UserID)
	if err != nil {
		return nil, apperrors.User.InvalidUserID
	}

	var (
		condition = bson.M{"userId": uid}
		result    = make([]domain.UserTerm, 0)
	)

	if !filter.Time.IsZero() {
		condition["createdAt"] = bson.M{"$lte": filter.Time}
	}
	if filter.IsFavourite != nil {
		condition["isFavourite"] = *filter.IsFavourite
	}
	if filter.Keyword != "" {
		condition["term"] = bson.M{"$regex": filter.Keyword, "$options": "i"}
	}

	// find
	cursor, err := r.collection().Find(ctx.Context(), condition, &options.FindOptions{
		Sort: bson.M{"createdAt": -1},
	}, &options.FindOptions{
		Limit: &filter.Limit,
	})
	if err != nil {
		return result, err
	}
	// parse
	defer func() { _ = cursor.Close(ctx.Context()) }()

	// parse
	var docs []model.UserTerm
	if err = cursor.All(ctx.Context(), &docs); err != nil {
		return result, err
	}

	// map data
	for _, doc := range docs {
		result = append(result, doc.ToDomain())
	}
	return result, nil
}
