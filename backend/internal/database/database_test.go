package database

import (
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/wayne-stays/backend/internal/models"
	"github.com/wayne-stays/backend/internal/testutil"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}

	stop, err := testutil.StartPostgres()
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	stop()
	os.Exit(code)
}

func count(t *testing.T, db *gorm.DB, model interface{}) int64 {
	t.Helper()
	var n int64
	if err := db.Model(model).Count(&n).Error; err != nil {
		t.Fatalf("count failed: %v", err)
	}
	return n
}

func TestMigrateAndSeedFreshDatabase(t *testing.T) {
	db := testutil.NewDatabase(t)

	// Running twice must be a no-op the second time, as it is on every restart
	for i := 0; i < 2; i++ {
		if err := Migrate(db); err != nil {
			t.Fatalf("Migrate run %d: %v", i+1, err)
		}
		if err := Seed(db); err != nil {
			t.Fatalf("Seed run %d: %v", i+1, err)
		}
	}

	if n := count(t, db, &models.Accommodation{}); n != 5 {
		t.Errorf("accommodations = %d, want 5", n)
	}
	if n := count(t, db, &models.Merchandise{}); n != 10 {
		t.Errorf("merchandise = %d, want 10", n)
	}
}

func TestSeedFillsEachTableIndependently(t *testing.T) {
	db := testutil.NewDatabase(t)
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}
	if err := Seed(db); err != nil {
		t.Fatal(err)
	}

	// Simulates a database seeded before merchandise existed
	if err := db.Exec("DELETE FROM merchandise").Error; err != nil {
		t.Fatal(err)
	}
	if err := Seed(db); err != nil {
		t.Fatal(err)
	}

	if n := count(t, db, &models.Merchandise{}); n != 10 {
		t.Errorf("merchandise = %d, want 10 after reseeding an empty table", n)
	}
	if n := count(t, db, &models.Accommodation{}); n != 5 {
		t.Errorf("accommodations = %d, want 5 (must not be duplicated)", n)
	}
}

// Merchandise mirrors the model as it was before versioned migrations. Declared
// here without a TableName method, GORM names its table "merchandises".
type Merchandise struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Price       int64  `gorm:"not null"`
	Category    string `gorm:"index"`
	ImageURL    string
	IsNew       bool `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func TestMigrateAdoptsAutoMigratedDatabase(t *testing.T) {
	db := testutil.NewDatabase(t)

	// Recreate the schema the old AutoMigrate-based startup produced
	if err := db.AutoMigrate(&models.Accommodation{}, &Merchandise{}); err != nil {
		t.Fatal(err)
	}
	legacy := Merchandise{Name: "Gotham Skyline Silk Tie", Price: 18500, Category: "Apparel"}
	if err := db.Create(&legacy).Error; err != nil {
		t.Fatal(err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate on legacy schema: %v", err)
	}

	var exists bool
	db.Raw("SELECT to_regclass('public.merchandises') IS NOT NULL").Scan(&exists)
	if exists {
		t.Error("legacy table merchandises still exists after migration")
	}

	var item models.Merchandise
	if err := db.First(&item, "id = ?", legacy.ID).Error; err != nil {
		t.Fatalf("legacy row not found in merchandise: %v", err)
	}
	if item.Category != "apparel" {
		t.Errorf("category = %q, want %q", item.Category, "apparel")
	}
}
