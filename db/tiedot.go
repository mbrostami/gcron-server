package db

import (
	"encoding/json"

	"github.com/HouzuoGuo/tiedot/db"

	pb "github.com/mbrostami/gcron/grpc"
	log "github.com/sirupsen/logrus"
)

// BleveDB index
type TiedotDB struct {
	db   *db.DB
	logs *db.Col
}

// NewTiedot creates TiedotDB
func NewTiedot() (TiedotDB, error) {
	myDBDir := "/tmp/Gcron"
	myDB, err := db.OpenDB(myDBDir)
	if err != nil {
		log.Fatal("DB Open error")
	}
	var logCol *db.Col
	if !myDB.ColExists("Logs") {
		if err := myDB.Create("Logs"); err != nil {
			log.Fatal("Collection coulnd't be created!")
		}
		logCol = myDB.Use("Logs")
		if err := logCol.Index([]string{"Success"}); err != nil {
			log.Fatal("Couldn't create Success index!")
		}
		if err := logCol.Index([]string{"Output"}); err != nil {
			log.Fatal("Couldn't create Output index!")
		}
	} else {
		logCol = myDB.Use("Logs")
	}
	return TiedotDB{db: myDB, logs: logCol}, err
}

// Store task data into db
func (b TiedotDB) Store(task *pb.Task) (string, error) {
	// Cast task struct to map
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(task)
	json.Unmarshal(inrec, &inInterface)
	docID, err := b.logs.Insert(inInterface)

	if err != nil {
		log.Fatalf("Store failed %v", err)
	}
	log.Infof("Stored into database, ID: %v", docID)
	return string(docID), err
}

// Search search index
func (b TiedotDB) Search(text string, limit int) {
	// Native Array
	query := map[string]interface{}{
		"eq":    text, // FIXME output is []byte
		"in":    []interface{}{"Output"},
		"limit": limit,
	}
	queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

	if err := db.EvalQuery(query, b.logs, &queryResult); err != nil {
		panic(err)
	}
	// Query result are document IDs
	for id := range queryResult {
		// To get query result document, simply read it
		readBack, err := b.GetByID(id)
		if err != nil {
			panic(err)
		}
		log.Printf("Doc found query: %+v", readBack)
	}
	log.Fatalf("Query returned document: %+v", queryResult)
}

// GetByID get doc by id
func (b TiedotDB) GetByID(id int) (*pb.Task, error) {
	readBack, err := b.logs.Read(id)
	if err != nil {
		log.Fatal(err)
	}
	task := pb.Task{}

	inrec, _ := json.Marshal(readBack)
	json.Unmarshal(inrec, &task)

	return &task, err
}

// Close db
func (b TiedotDB) Close() {
	// Gracefully close database
	if err := b.db.Close(); err != nil {
		log.Fatalf("DB Close failed: %v", err)
	}
}
