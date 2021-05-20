package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var client *mongo.Client

func Init() {
	mongoHost := os.Getenv("MONGO_HOST")
	mongoUser := os.Getenv("MONGO_USERNAME")
	mongoPass := os.Getenv("MONGO_PASSWORD")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", mongoUser, mongoPass, mongoHost, mongoPort, mongoDatabase)
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success connection mongo: ", url)
}

func MongoClient() *mongo.Client {
	return client
}

func FindMany(database string, collection string, querys map[string]interface{}) ([]bson.M, error) {
	return FindManyWithLimit(database, collection, querys, 10000)
}

func FindManyWithLimit(database string, collection string, conditions map[string]interface{}, limit int64) ([]bson.M,
	error) {
	var result []bson.M
	filter := bson.M{}
	findOpt := options.FindOptions{}
	findOpt.SetLimit(limit)
	for k, v := range conditions {
		filter[k] = v
	}
	coll := client.Database(database).Collection(collection)
	cur, err := coll.Find(context.TODO(), filter, &findOpt)
	if err != nil {
		return result, err
	}

	for cur.Next(context.TODO()) {
		var elem bson.M
		err := cur.Decode(&elem)
		if err != nil {
			return result, err
		}
		result = append(result, elem)
	}
	return result, nil
}

func FindOne(db string, coll string, q map[string]interface{}) (bson.M, error) {
	var result bson.M
	var err error
	query := bson.M{}
	for k, v := range q {
		query[k] = v
	}
	collection := client.Database(db).Collection(coll)
	res := collection.FindOne(context.TODO(), query)
	err = res.Err()
	if err != nil {
		return result, err
	}
	err = res.Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func UpdateOne(db string, coll string, q map[string]interface{}, value interface{}) (*mongo.UpdateResult, error) {
	var err error
	query := bson.M{}
	for k, v := range q {
		query[k] = v
	}
	collection := client.Database(db).Collection(coll)
	res, err := collection.UpdateOne(context.TODO(), query, value)
	if err != nil {
		return res, err
	}
	return nil, err

}

func DeleteMany(db string, coll string, condition map[string]interface{}) (*mongo.DeleteResult, error){
	var err error
	cond := bson.M{}
	for k,v := range condition {
		cond[k] = v
	}
	collection := client.Database(db).Collection(coll)
	res, err := collection.DeleteMany(context.TODO(), cond, nil)
	if err != nil {
		return res, err
	}
	return nil, err
}