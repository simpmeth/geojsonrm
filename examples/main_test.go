package examples

import (
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/simpmeth/geojsonrm/examples/testutil"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	conn, closer := testutil.InitTempDB()
	db = conn

	code := m.Run()

	closer()

	os.Exit(code)
}
