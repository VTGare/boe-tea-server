package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

//Connect connects to a MongoDB database.
func Connect(uri, db string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &DB{client, client.Database(db)}, nil
}

func (db *DB) Close() error {
	return db.Client.Disconnect(context.Background())
}

/*func (db *mongodb) GetGuild(id string) (*models.Guild, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := db.Database.Collection(guilds).FindOne(ctx, bson.D{
		{Key: "guild_id", Value: id},
	})

	guild := &models.Guild{}

	err := res.Decode(guild)
	if err != nil {
		return nil, err
	}

	return guild, nil
}

func (db *mongodb) GetArtwork(id int) (*models.Artwork, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := db.Database.Collection(artworks).FindOne(ctx, bson.D{
		{Key: "artwork_id", Value: id},
	})

	artwork := &models.Artwork{}

	err := res.Decode(artwork)
	if err != nil {
		return nil, err
	}

	return artwork, nil
}

func (db *mongodb) AddArtwork(art *models.Artwork) (*models.Artwork, error) {
	id, err := db.artworkID()
	if err != nil {
		return nil, err
	}

	art.ID = id
	art.Favourites = 0
	art.NSFW = 0
	art.CreatedAt = time.Now()
	art.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = db.Database.Collection(artworks).InsertOne(ctx, art)
	return art, nil
}

func (db *mongodb) artworkID() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := db.Database.Collection(counters).FindOneAndUpdate(
		ctx,
		bson.D{{Key: "_id", Value: "artworks"}},
		bson.D{{Key: "$inc", Value: bson.D{{Key: "counter", Value: 1}}}},
		options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true),
	)

	counter := &struct {
		ID      string `bson:"_id"`
		Counter int    `bson:"counter"`
	}{}

	err := res.Decode(counter)
	if err != nil {
		return 0, err
	}

	return counter.Counter, nil
}*/
