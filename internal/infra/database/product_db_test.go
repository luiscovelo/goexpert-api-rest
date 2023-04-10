package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/luiscovelo/goexpert-api-rest/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateNewProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10)
	assert.Nil(t, err)
	assert.NotNil(t, product)

	productDB := NewProduct(db)

	err = productDB.Create(product)
	assert.Nil(t, err)
}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&entity.Product{})

	for i := 1; i <= 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.Nil(t, err)

		err = db.Create(product).Error
		assert.Nil(t, err)
	}

	productDB := NewProduct(db)

	products, err := productDB.FindAll(1, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDB.FindAll(3, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 4)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 24", products[3].Name)
}

func TestFindByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, product)

	err = db.Create(product).Error
	assert.Nil(t, err)

	productDB := NewProduct(db)

	productFound, err := productDB.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.NotNil(t, productFound)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, "Product 1", productFound.Name)
	assert.Equal(t, 10.0, productFound.Price)
}

func TestUpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, product)

	err = db.Create(product).Error
	assert.Nil(t, err)

	product.Name = "Product 2"
	product.Price = 20.0

	productDB := NewProduct(db)

	err = productDB.Update(product)
	assert.Nil(t, err)

	productFound, err := productDB.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, "Product 2", productFound.Name)
	assert.Equal(t, 20.0, productFound.Price)
}

func TestDeleteProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, product)

	err = db.Create(product).Error
	assert.Nil(t, err)

	productDB := NewProduct(db)

	err = productDB.Delete(product.ID.String())
	assert.Nil(t, err)

	productFound, err := productDB.FindByID(product.ID.String())
	assert.NotNil(t, err)
	assert.Nil(t, productFound)
}
