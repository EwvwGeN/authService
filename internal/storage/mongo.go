package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/EwvwGeN/authService/internal/config"
	"github.com/EwvwGeN/authService/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoProvider struct {
	cfg config.MongoConfig
	db *mongo.Database
}

var (
	ErrDbNotExist = errors.New("database does not exist")
	ErrCollNotExist = errors.New("some of collections does not exist")
)

func NewMongoProvider(ctx context.Context, cfg config.MongoConfig) (*mongoProvider, error) {
	client, err := mongo.Connect(ctx,
		options.Client().
		ApplyURI(
			fmt.Sprintf("%s://%s:%s@%s:%s/?authSource=%s",
			cfg.ConectionFormat,
			cfg.User,
			cfg.Password, 
			cfg.Host,
			cfg.Port,
			cfg.AuthSourse,
			)),
		)
	if err != nil {
		return nil, err
	}
	dbList, err := client.ListDatabases(ctx, bson.D{{Key: "name", Value: cfg.Database}})
	if err != nil {
		return nil, err
	}
	if dbList.TotalSize == 0 {
		return nil, ErrDbNotExist
	}
	db := client.Database(cfg.Database)
	colList, err := db.ListCollectionNames(ctx,
		bson.M{
			"$or": []interface{}{
				bson.D{{Key: "name", Value: cfg.UserCollection}},
				bson.D{{Key: "name", Value: cfg.AppCollection}},
			}})
	if err != nil {
		return nil, err
	}
	if len(colList) != 2 {
		return nil, ErrCollNotExist
	}
	return &mongoProvider{
		cfg: cfg,
		db: db,
	}, nil
}

func (mp *mongoProvider) SaveUser( ctx context.Context,
	email string,
	passHash []byte,
) (string, error) {
	inRes, err := mp.db.Collection(mp.cfg.UserCollection).InsertOne(ctx, models.User{
		Email: email,
		PassHash: passHash,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", ErrUserExist
		}
		return "", err
	}
	return inRes.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (mp *mongoProvider) GetUser(ctx context.Context,
	email string,
) (models.User, error) {
	user := mp.db.Collection(mp.cfg.UserCollection).FindOne(ctx, bson.D{{Key: "email", Value: email}})
	if m, _ := user.Raw(); m != nil {
		var outUser models.User
		user.Decode(&outUser)
		outUser.Id = m.Index(0).Value().ObjectID().Hex()
		return outUser, nil
	}
	return models.User{}, ErrUserNotFound
}

func (mp *mongoProvider) IsAdmin(ctx context.Context,
	userId string,
) (bool, error) {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return false, err
	}
	user := mp.db.Collection(mp.cfg.UserCollection).FindOne(ctx, bson.D{{Key: "_id", Value: objectId}})
	if m, _ := user.Raw(); m == nil {
		return false, ErrUserNotFound
	}
	var decodedUser models.User
	user.Decode(&decodedUser)
	return decodedUser.IsAdmin, nil
}

func (mp *mongoProvider) GetApp(ctx context.Context,
	appId string,
) (models.App, error) {
	objectId, err := primitive.ObjectIDFromHex(appId)
	if err != nil {
		return models.App{}, err
	}
	app := mp.db.Collection(mp.cfg.AppCollection).FindOne(ctx, bson.D{{Key: "_id", Value: objectId}})
	if m, _ := app.Raw(); m != nil {
		var outApp models.App
		app.Decode(&outApp)
		outApp.Id = m.Index(0).Value().ObjectID().Hex()
		return outApp, nil
	}
	return models.App{}, ErrAppNotFound
}