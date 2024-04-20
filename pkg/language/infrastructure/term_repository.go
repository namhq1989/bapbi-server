package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

type TermRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewTermRepository(db *mongo.Database) TermRepository {
	r := TermRepository{
		db:             db,
		collectionName: database.Tables.LanguageTerm,
	}
	r.ensureIndexes()
	return r
}

func (r TermRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "term", Value: 1}, {Key: "from.language", Value: 1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r TermRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r TermRepository) FindByID(ctx *appcontext.AppContext, termID string) (*domain.Term, error) {
	id, err := database.ObjectIDFromString(termID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.Term
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

func (r TermRepository) FindByTerm(ctx *appcontext.AppContext, term, fromLanguage string) (*domain.Term, error) {
	// find
	var doc model.Term
	if err := r.collection().FindOne(ctx.Context(), bson.M{
		"term":          strings.ToLower(term),
		"from.language": fromLanguage,
	}).Decode(&doc); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	// respond
	result := doc.ToDomain()
	return &result, nil
}

func (r TermRepository) CreateTerm(ctx *appcontext.AppContext, term domain.Term) error {
	// convert to mongodb model
	doc, err := model.Term{}.FromDomain(term)
	if err != nil {
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	return err
}

func (r TermRepository) UpdateTerm(ctx *appcontext.AppContext, term domain.Term) error {
	// convert to mongodb model
	doc, err := model.Term{}.FromDomain(term)
	if err != nil {
		return err
	}

	_, err = r.collection().UpdateByID(ctx.Context(), doc.ID, bson.M{"$set": doc})
	return err
}
