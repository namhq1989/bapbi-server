package infrastructure

import (
	"errors"

	"github.com/namhq1989/bapbi-server/internal/database"
	"github.com/namhq1989/bapbi-server/internal/utils/appcontext"
	apperrors "github.com/namhq1989/bapbi-server/internal/utils/error"
	"github.com/namhq1989/bapbi-server/pkg/user/domain"
	"github.com/namhq1989/bapbi-server/pkg/user/infrastructure/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHub struct {
	db             *mongo.Database
	collectionName string
}

func NewUserHub(db *mongo.Database) UserHub {
	return UserHub{
		db:             db,
		collectionName: database.Tables.User,
	}
}

func (r UserHub) collection() *mongo.Collection {
	return r.db.Collection(r.collectionName)
}

func (r UserHub) FindOneByID(ctx *appcontext.AppContext, id string) (*domain.User, error) {
	oid, err := database.ObjectIDFromString(id)
	if err != nil {
		return nil, apperrors.Common.InvalidID
	}

	// find
	var doc model.User
	if err = r.collection().FindOne(ctx.Context(), bson.M{
		"_id": oid,
	}).Decode(&doc); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		ctx.Logger().Error("failed to find user by id", err, appcontext.Fields{"id": id})
		return nil, err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	// respond
	result := doc.ToDomain()
	return &result, nil
}

func (r UserHub) FindOneByEmail(ctx *appcontext.AppContext, email string) (*domain.User, error) {
	// find
	var doc model.User
	if err := r.collection().FindOne(ctx.Context(), bson.M{
		"email": email,
	}).Decode(&doc); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		ctx.Logger().Error("failed to find user by email", err, appcontext.Fields{"email": email})
		return nil, err
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	// respond
	result := doc.ToDomain()
	return &result, nil
}

func (r UserHub) CreateUser(ctx *appcontext.AppContext, user domain.User) (*domain.User, error) {
	// convert to mongodb model
	doc, err := model.User{}.FromDomain(user)
	if err != nil {
		return nil, err
	}

	_, err = r.collection().InsertOne(ctx.Context(), &doc)
	if err != nil {
		ctx.Logger().Error("failed to insert user", err, appcontext.Fields{"user": user})
		return nil, err
	}

	// respond
	result := doc.ToDomain()
	return &result, nil
}
