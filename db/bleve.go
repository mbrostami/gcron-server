package db

import (
	"encoding/json"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
	pb "github.com/mbrostami/gcron/grpc"
	log "github.com/sirupsen/logrus"
)

// BleveDB index
type BleveDB struct {
	db bleve.Index
}

// NewBleve creates BleveDB
func NewBleve() (BleveDB, error) {
	dbName := "example.bleve"
	index, err := bleve.Open(dbName)
	if err != nil {
		mapping := mappings()
		index, err = bleve.New(dbName, mapping)
	}
	return BleveDB{db: index}, err
}

// Store task data into db
func (b BleveDB) Store(task *pb.Task) error {
	log.Info("Storing into database")
	return b.db.Index(task.GUID, task)
}

// Close task data into db
func (b BleveDB) Close() {
	b.db.Close()
}

// Search search index
func (b BleveDB) Search(text string, limit int) {
	query := bleve.NewMatchQuery(text)
	search := bleve.NewSearchRequest(query)
	searchResults, err := b.db.Search(search)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	log.Fatalf("Docs: %+v", getOriginalDocsFromSearchResults(searchResults, b.db))
	// dc, err := b.db.Document("bq0cgc2uof2v4eo9btd0")
	// if err != nil {
	// 	log.Fatalf("Search failed: %v", err)
	// }
	// log.Fatalf("Doc: %+v", dc.GoString())
	// log.Fatalf("%+v", searchResults)
}

func getOriginalDocsFromSearchResults(
	results *bleve.SearchResult,
	index bleve.Index,
) [][]byte {
	docs := make([][]byte, 0)

	for _, val := range results.Hits {
		id := val.ID
		doc, err := index.Document(id)
		log.Fatalf("ODoc:%+v", doc.GoString())
		if err != nil {
			log.Fatal("Trouble getting internal doc:", err)
		}
		rv := struct {
			ID     string                 `json:"id"`
			Fields map[string]interface{} `json:"fields"`
		}{
			ID:     id,
			Fields: map[string]interface{}{},
		}
		for _, field := range doc.Fields {
			var newval interface{}
			switch field := field.(type) {
			case *document.BooleanField:
				newval = true
				if string(field.Value()) == "F" {
					newval = false
				}
			case *document.TextField:
				newval = string(field.Value())
				log.Infof("Text Field: %+v", field.Name())
				if field.Name() == "Output" {
					log.Fatalf("Output: %+v", string(field.Value()))
				}
			case *document.NumericField:
				n, err := field.Number()
				if err == nil {
					newval = n
				}
			case *document.DateTimeField:
				d, err := field.DateTime()
				if err == nil {
					newval = d.Format(time.RFC3339Nano)
				}
			}
			rv.Fields[field.Name()] = newval
		}
		log.Fatalf("Doc: %v", rv)
		j2, _ := json.MarshalIndent(rv, "", "    ")
		docs = append(docs, j2)
	}
	return docs
}

func mappings() *mapping.IndexMappingImpl {
	mapping := bleve.NewIndexMapping()

	logsMapping := bleve.NewDocumentStaticMapping()
	boolNotIndexed := bleve.NewBooleanFieldMapping()
	boolNotIndexed.Index = false
	boolIndexed := bleve.NewBooleanFieldMapping()

	textNotIndexed := bleve.NewTextFieldMapping()
	textNotIndexed.Index = false
	textIndexed := bleve.NewTextFieldMapping()

	numericNotIndexed := bleve.NewNumericFieldMapping()
	numericNotIndexed.Index = false
	numericIndexed := bleve.NewNumericFieldMapping()

	dateTimeIndexed := bleve.NewDateTimeFieldMapping()

	logsMapping.AddFieldMappingsAt("FLock", boolNotIndexed)
	logsMapping.AddFieldMappingsAt("Success", boolIndexed)
	logsMapping.AddFieldMappingsAt("FLockName", textNotIndexed)
	logsMapping.AddFieldMappingsAt("FOverride", textNotIndexed)
	logsMapping.AddFieldMappingsAt("GUID", textNotIndexed) // already using as doc id
	logsMapping.AddFieldMappingsAt("Parent", textIndexed)
	logsMapping.AddFieldMappingsAt("Hostname", textIndexed)
	logsMapping.AddFieldMappingsAt("Username", textIndexed)
	logsMapping.AddFieldMappingsAt("Command", textIndexed)
	logsMapping.AddFieldMappingsAt("Output", textIndexed)
	logsMapping.AddFieldMappingsAt("FDelay", numericNotIndexed)
	logsMapping.AddFieldMappingsAt("Pid", numericNotIndexed)
	logsMapping.AddFieldMappingsAt("UID", numericIndexed)
	logsMapping.AddFieldMappingsAt("ExitCode", numericIndexed)
	logsMapping.AddFieldMappingsAt("StartTime", dateTimeIndexed)
	logsMapping.AddFieldMappingsAt("EndTime", dateTimeIndexed)
	logsMapping.AddFieldMappingsAt("SystemTime", dateTimeIndexed)
	logsMapping.AddFieldMappingsAt("UserTime", dateTimeIndexed)
	mapping.DefaultMapping = logsMapping
	mapping.DefaultType = "_doc"
	mapping.AddDocumentMapping("_doc", logsMapping)
	return mapping
}
