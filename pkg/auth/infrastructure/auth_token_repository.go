package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/auth/domain"
	"github.com/namhq1989/bapbi-server/pkg/auth/infrastructure/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthTokenRepository struct {
	db             *mongo.Database
	collectionName string
}

func NewAuthTokenRepository(db *mongo.Database) AuthTokenRepository {
	r := AuthTokenRepository{
		db:             db,
		collectionName: database.Tables.AuthToken,
	}
	r.ensureIndexes()
	return r
}

func (r AuthTokenRepository) ensureIndexes() {
	var (
		ctx     = context.Background()
		opts    = options.CreateIndexes().SetMaxTime(time.Minute * 30)
		indexes = []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "user", Value: 1}, {Key: "refreshToken", Value: 1}},
			},
		}
	)

	if _, err := r.collection().Indexes().CreateMany(ctx, indexes, opts); err != nil {
		fmt.Printf("index collection %s err: %v \n", r.collectionName, err)
	}
}

func (r AuthTokenRepository) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r AuthTokenRepository) CreateAuthToken(ctx *appcontext.AppContext, token domain.RefreshToken) error {
	// convert to mongodb model
	doc, err := model.AuthToken{}.FromDomain(token)
	if err != nil {
		ctx.Logger().Error("failed to convert auth token", err, appcontext.Fields{"token": token})
		return err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	if err != nil {
		ctx.Logger().Error("failed to insert refresh token", err, appcontext.Fields{"token": token})
	}

	return err
}

func (r AuthTokenRepository) DeleteAuthToken(ctx *appcontext.AppContext, token domain.RefreshToken) error {
	// convert to mongodb model
	doc, err := model.AuthToken{}.FromDomain(token)
	if err != nil {
		ctx.Logger().Error("failed to convert auth token", err, appcontext.Fields{"token": token})
		return err
	}

	_, err = r.collection().DeleteOne(ctx.Context(), bson.M{
		"user":         doc.User,
		"refreshToken": doc.RefreshToken,
	})
	if err != nil {
		ctx.Logger().Error("failed to delete refresh token", err, appcontext.Fields{"token": token})
	}

	return err
}

func (r AuthTokenRepository) FindAuthToken(ctx *appcontext.AppContext, userID, refreshToken string) (*domain.RefreshToken, error) {
	uID, err := database.ObjectIDFromString(userID)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	var doc *domain.RefreshToken
	err = r.collection().FindOne(ctx.Context(), bson.M{
		"user":         uID,
		"refreshToken": refreshToken,
	}).Decode(&doc)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		ctx.Logger().Error("failed to find refresh token", err, appcontext.Fields{"user": userID, "refreshToken": refreshToken})
		return nil, err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperrors.Auth.InvalidAuthToken
	}

	return doc, nil
}
