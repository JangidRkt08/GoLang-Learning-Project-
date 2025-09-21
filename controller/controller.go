package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jangidRkt08/mongoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://ravi1007:ravi1007@cluster0.x2js4d6.mongodb.net/?retryWrites=true&w=majority&authSource=admin"
// const connectionString ="mongodb+srv://rktbruce01:ravi1007@cluster0.x2js4d6.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const dbName = "Netflix"

const collectionName = "Watchlist"


// add collection

var collection *mongo.Collection

// connect with mongoDb

func init(){
	// init function run only for one time and only once

	// client opotion
	clientOption := options.Client().ApplyURI(connectionString)

	// connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	

	fmt.Println("Connected to MongoDB!!")

	collection = client.Database(dbName).Collection(collectionName)

	// collection instance 
	fmt.Println("Collection instance is ready")


}

// ------> MONGO HELPERS - file  <------


// insert one Record
func insertOneMovie(movie model.Netflix){
	inserted, err := collection.InsertOne(context.Background(),movie) 
	if err != nil {
		log.Fatal(err)
	}
		
	fmt.Println("Inserted movie in db with id: ", inserted.InsertedID)
	
}

// Update One Record
func updateOneMovie(movieId string){
	id, _ := primitive.ObjectIDFromHex(movieId)   //string to object id 

	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"watched" : true}}

	result, _ :=collection.UpdateOne(context.Background(),filter,update)

	fmt.Println("modified count: ", result.ModifiedCount)
}

// Delete One Record
func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)   //string to object id

	filter := bson.M{"_id":id}
	deleteCount, err := collection.DeleteOne(context.Background(),filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Movie got deleted with id: ", movieId)
	fmt.Println("Deleted Count: ", deleteCount.DeletedCount)
}


// Delete all records
func deleteAllMovies() int64{
	deleteResult, err := collection.DeleteMany(context.Background(),bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Number of movies deleted: ", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

// get All Movies
func getAllMovies() []primitive.M{
	curr, err :=collection.Find(context.Background(),bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var movies []primitive.M
	for curr.Next(context.Background()){
		var movie bson.M
		err := curr.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	defer curr.Close(context.Background())

	return movies
}



//-------> Actual Controllers - file <------

func GetMyAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/x-www-form-urlencoded")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)

}

func CreateMovie(w http.ResponseWriter,r *http.Request){
	w.Header().Set("content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-control-Allow-Methods","POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)

}

func MarkAsWatched(w http.ResponseWriter,r *http.Request){
	w.Header().Set("content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-control-Allow-Methods","PUT")

	params := mux.Vars(r)

	updateOneMovie(params["id"])

	json.NewEncoder(w).Encode(params["id"])

}


func DeleteOneMovie(w http.ResponseWriter,r *http.Request){
	w.Header().Set("content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-control-Allow-Methods","DELETE")

	params := mux.Vars(r)

	deleteOneMovie(params["id"])

	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAllMovies(w http.ResponseWriter,r *http.Request){
	w.Header().Set("content-type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-control-Allow-Methods","DELETE")

	count := deleteAllMovies()

	json.NewEncoder(w).Encode(count)

}

