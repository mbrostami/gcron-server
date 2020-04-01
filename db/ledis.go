package db

import (
	"encoding/json"
	"time"
	"unsafe"

	pb "github.com/mbrostami/gcron/grpc"
	"github.com/rs/xid"
	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
	log "github.com/sirupsen/logrus"
)

// LedisDB database
type LedisDB struct {
	ledis *ledis.Ledis
	db    *ledis.DB
}

// NewLedis create ledisdb instance
func NewLedis() *LedisDB {
	cfg := config.NewConfigDefault()
	ledis, err := ledis.Open(cfg)
	if err != nil {
		log.Fatalf("DB Connect error! %v", err)
	}
	db, _ := ledis.Select(0)
	return &LedisDB{db: db}
}

// Store data in db
func (l LedisDB) Store(task *pb.Task) (string, error) {
	key := task.GetUID()
	byteKeys := (*[4]byte)(unsafe.Pointer(&key))[:] // 32 bit id (4 byte)

	guid, _ := xid.FromString(task.GetGUID())

	jsonByte, err := json.Marshal(&task)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	score1 := ledis.ScorePair{
		Score:  guid.Time().Unix(),
		Member: jsonByte,
	}
	number, err := l.db.ZAdd(byteKeys, score1)
	if err != nil {
		log.Fatalf("DB Store error! %v", err)
	}
	return string(number), nil
}

// Get members of a key
func (l LedisDB) Get(uid uint32, start int, stop int) {
	byteKeys := (*[4]byte)(unsafe.Pointer(&uid))[:] // 32 bit id (4 byte)
	scorePairs, _ := l.db.ZRange(byteKeys, start, stop)
	for _, scorePair := range scorePairs {
		score := scorePair.Score
		member := scorePair.Member
		unixTimeUTC := time.Unix(score, 0)
		log.Printf("Score: %v", unixTimeUTC.Format(time.RFC3339))
		task := &pb.Task{}
		json.Unmarshal(member, &task)
		log.Printf("Member: %+v", task)
	}
	log.Fatal("Scores")
}

// Close members of a key
func (l LedisDB) Close() {
	l.ledis.Close()
}
