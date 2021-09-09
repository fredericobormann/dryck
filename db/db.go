package db

import (
	"github.com/fredericobormann/dryck/models"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// DB stores database connection
type DB struct {
	Datastore *gorm.DB
}

// New creates a new database datastore
func New(dbType string, dbConnInfo string) *DB {
	db, err := gorm.Open(dbType, dbConnInfo)
	if err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Drink{})
	db.AutoMigrate(&models.Purchase{})
	db.AutoMigrate(&models.Payment{})
	db.Model(&models.Purchase{}).AddForeignKey("customer_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Purchase{}).AddForeignKey("product_id", "drinks(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Payment{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	db.Exec("UPDATE purchases " +
		"SET price = drinks.price " +
		"FROM drinks " +
		"WHERE purchases.product_id = drinks.id " +
		"AND (purchases.price = 0 " +
		"OR purchases.price IS NULL);")

	dryckdb := DB{
		db,
	}

	return &dryckdb
}

// Close closes current database connection
func (db *DB) Close() {
	err := db.Datastore.Close()
	if err != nil {
		log.Print(err)
	}
}

// GetAllUsers returns all users currently in the database
func (db *DB) GetAllUsers() []models.User {
	var allUsers []models.User
	db.Datastore.Find(&allUsers)
	return allUsers
}

// GetUsername returns the username for a given user id
func (db *DB) GetUsername(userID uint) string {
	var user models.User
	db.Datastore.Where("id = ?", userID).First(&user)
	return user.Name
}

// CreateNewUser creates a new user with given username
func (db *DB) CreateNewUser(username string) {
	newUser := models.User{Name: username}
	db.Datastore.Create(&newUser)
}

// GetAllDrinks returns all drinks in the database
func (db *DB) GetAllDrinks() []models.Drink {
	var allDrinks []models.Drink
	db.Datastore.Find(&allDrinks)
	return allDrinks
}

// GetPurchasesOfUser returns all purchases of one user specified by user id with drink information preloaded
func (db *DB) GetPurchasesOfUser(userID uint) []models.Purchase {
	var purchases []models.Purchase
	db.Datastore.Preload("Product").Where("customer_id = ?", userID).Order("purchase_time desc").Find(&purchases)
	return purchases
}

// GetPaginatedPurchasesOfUser returns purchases for a given user on a specific page (beginning with 1) with given numberOfItems on each page
func (db *DB) GetPaginatedPurchasesOfUser(userID uint, itemsPerPage int, page int) []models.Purchase {
	var purchases []models.Purchase
	offset := (page - 1) * itemsPerPage
	db.Datastore.Limit(itemsPerPage).Offset(offset).Preload("Product").Where("customer_id = ?", userID).Order("purchase_time desc").Find(&purchases)
	return purchases
}

// GetNumberOfPurchasesOfUser returns the number of purchases by a given user
func (db *DB) GetNumberOfPurchasesOfUser(userID uint) int {
	var count int64
	db.Datastore.Model(&models.Purchase{}).Preload("Product").Where("customer_id = ?", userID).Order("purchase_time desc").Count(&count)
	return int(count)
}

// GetTotalDebtOfUser returns the total debt of one user specified by their user id
func (db *DB) GetTotalDebtOfUser(userID uint) int {
	var totalDebt int
	db.Datastore.Table("purchases").Where("customer_id = ?", userID).Select("sum(price)").Row().Scan(&totalDebt)
	var totalPayed int
	db.Datastore.Table("payments").Where("user_id = ?", userID).Select("sum(amount)").Row().Scan(&totalPayed)
	return totalDebt - totalPayed
}

// PurchaseDrink adds a purchase for one user specified by their id and a drink also specified by id
func (db *DB) PurchaseDrink(userID uint, drinkID uint) {
	var drink models.Drink
	db.Datastore.Where("id = ?", drinkID).First(&drink)
	purchase := models.Purchase{CustomerID: userID, ProductID: drinkID, PurchaseTime: time.Now(), Price: drink.Price}
	db.Datastore.Create(&purchase)
}

// DeletePurchase deletes a purchase from database
func (db *DB) DeletePurchase(purchaseID uint) {
	db.Datastore.Where("id = ?", purchaseID).Unscoped().Delete(&models.Purchase{})
}

// AddPayment adds a payment for a user specified by id with a certain amount
func (db *DB) AddPayment(userID uint, amount int) {
	payment := models.Payment{UserID: userID, Amount: amount, PaymentTime: time.Now()}
	db.Datastore.Create(&payment)
}

// DeletePayment deletes a payment specified by its id
func (db *DB) DeletePayment(paymentID uint) {
	db.Datastore.Where("id = ?", paymentID).Unscoped().Delete(&models.Payment{})
}

// GetAllPaymentsOfUser returns all payments of one user specified by id
func (db *DB) GetAllPaymentsOfUser(userID uint) []models.Payment {
	var payments []models.Payment
	db.Datastore.Where("user_id = ?", userID).Find(&payments)
	return payments
}
