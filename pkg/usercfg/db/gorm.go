package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SourceSystemType uint

const (
	SourceConfluence = iota
	SourceHTTP
)

type User struct {
	gorm.Model
	Name        string `gorm:"unique,index"`
	APIKey      string
	Collections []Collection
}

type Collection struct {
	gorm.Model
	UserID         uint
	DisplayName    string `gorm:"index"`        // DisplayName is the name displayed to the user
	CollectionName string `gorm:"unique,index"` // CollectionName is the internal unique name of the collection
	Sources        []SourceSystem
}

type SourceSystem struct {
	gorm.Model
	CollectionID uint
	Name         string
	Type         SourceSystemType
	URL          string
	Key          string
	Parts        string
}

func (u *User) SetName(n string) *User {
	u.Name = n
	return u
}

func (u *User) SetAPIKey(k string) *User {
	u.APIKey = k
	return u
}

func InitGorm() error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	u:=User{}
	u.SetAPIKey("jhhh").SetName("ddd")

	// Migrate the schema
	if err := db.AutoMigrate(&User{}, &Collection{}, &SourceSystem{}); err != nil {
		return fmt.Errorf("automigration DB: %w", err)
	}

	// ctx := context.Background()
	// // Create
	// err = gorm.G[Product](db).Create(ctx, &Product{Code: "D42", Price: 100})

	// // Read
	// product, err := gorm.G[Product](db).Where("id = ?", 1).First(ctx)       // find product with integer primary key
	// products, err := gorm.G[Product](db).Where("code = ?", "D42").Find(ctx) // find product with code D42

	// // Update - update product's price to 200
	// err = gorm.G[Product](db).Where("id = ?", product.ID).Update(ctx, "Price", 200)
	// // Update - update multiple fields
	// err = gorm.G[Product](db).Where("id = ?", product.ID).Updates(ctx, map[string]interface{}{"Price": 200, "Code": "F42"})

	// // Delete - delete product
	// err = gorm.G[Product](db).Where("id = ?", product.ID).Delete(ctx)
	return nil
}
