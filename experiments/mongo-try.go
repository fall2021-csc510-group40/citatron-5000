package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"time"

	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Stack Overflow coming in clutch with a helper method
func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func main() {
	//Driver example code to connect the client to a local DB with user guest
	clientOptions := options.Client().
		ApplyURI("mongodb://guest:guest@localhost/sandbox")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//Sanity Check
	fmt.Print("Connected to DB")

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	//Knowing the contents of a collection (SQL Table), search for contents
	coll := client.Database("sandbox").Collection("user")
	title := "Marty McFly"
	fmt.Print(title)
	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"name", title}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
		return
	}
	if err != nil {
		panic(err)
	}

	//Convert the contents from the Mongo BSON format to JSON for output
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

	//-------------------------------------------------------------------------------------------
	//Read in a CSV of test data to use for testing of data insertion
	test := readCsvFile("C:\\Users\\gage\\Downloads\\testData.csv")

	//Sanity check
	print(len(test))

	//Commented out since the test DB already has the data

	//Loop through each term in the collection, and create an object to insert in
	//the DB
	/*for i := 1; i < len(test); i++ {
		title := test[i][0]
		authors := test[i][1]
		venue := test[i][2]
		year := test[i][3]
		volume := test[i][4]
		issue := test[i][5]
		startPg := test[i][6]
		endPg := test[i][7]
		isbn := test[i][8]
		doi := test[i][9]
		keywords := test[i][10]
		publisher := test[i][11]
		count, err := coll.InsertOne(ctx, bson.D{
			{Key: "title", Value: title},
			{Key: "authors", Value: authors},
			{Key: "venue", Value: venue},
			{Key: "year", Value: year},
			{Key: "volume", Value: volume},
			{Key: "issue", Value: issue},
			{Key: "startPg", Value: startPg},
			{Key: "endPg", Value: endPg},
			{Key: "isbn", Value: isbn},
			{Key: "doi", Value: doi},
			{Key: "keywords", Value: keywords},
			{Key: "publisher", Value: publisher},
		},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("Successfull add " + string(i))
		fmt.Print(count)
	}*/
}
