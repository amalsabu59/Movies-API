package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mongoapi/model"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://admin:kdecWBbJKp16XZN5@cluster0.a1kbgi2.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const dbName = "netflix"
const colName = "Watchlist"

//most imnportant

var collection *mongo.Collection

// connnect with mongoDB

func init() {
	// client options
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongdb
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongo db connection sucess")

	collection = client.Database(dbName).Collection(colName)

	//collection insteance
	fmt.Println("Collection instance is ready")
}

//MONGO DB helpers  - file

// insert 1 record

func insertOneMovie(movie model.Netfilx){
  inserted, err :=	collection.InsertOne(context.Background(), movie)

  if err != nil {
	log.Fatal(err)
  }

  fmt.Println("inserted 1 movie in db with: ", inserted)
}

// update 1 record

func updateOneMovie(movieId string){
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"watched":true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified", result)
}


// delete one record

func deleteOneMovie(movieId string) {
	_id,_ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id":_id}
	deleteCount, err := collection.DeleteOne(context.Background(),filter)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Movie got delete with delete count: ", deleteCount)
}

	//delete all records form mongo db

	func deleteAllMovie() int64{
		deleteResult, err := collection.DeleteMany(context.Background(),bson.D{{}},nil);
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Number of movies deleted:",deleteResult.DeletedCount)
		return deleteResult.DeletedCount
	}

	// get all movies for database

	func getAllMovies() []primitive.M {
		cur,err := collection.Find(context.Background(),bson.D{{}})
		if err != nil {
			log.Fatal(err)
		}

		var movies []primitive.M

		for cur.Next(context.Background()){
			var movie bson.M
			err := cur.Decode(&movie)
			if err != nil {
				log.Fatal(err)
			}
			movies = append(movies, movie)
		}

		defer cur.Close(context.Background())

		return movies
	}

	// Actual controller - file 
	func GetMyAllMovies(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		allMovies := getAllMovies()

		json.NewEncoder(w).Encode(allMovies)

	}

	func CreateMovie(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Allow-Control-Allow-Methods","POST")

		var movie model.Netfilx
		_ = json.NewDecoder(r.Body).Decode(&movie)
		insertOneMovie(movie)
		json.NewEncoder(w).Encode(movie)

	}

	func MarkAsWatched(w http.ResponseWriter, r *http.Request){
				w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Allow-Control-Allow-Methods","PUT")
		params := mux.Vars(r)

		updateOneMovie(params["id"])
		json.NewEncoder(w).Encode(params["id"])

	}


	func DeleteAMovie(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Allow-Control-Allow-Methods","DELETE")

		params := mux.Vars(r)
		deleteOneMovie(params["id"])
		json.NewEncoder(w).Encode(params["id"])
	}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods","DELETE")
		
	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)

}