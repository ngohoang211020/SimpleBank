package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/ngohoang211020/simplebank/util"
	"log"
	"os"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	err := util.Config.LoadConfig("../../")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), util.Config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
