package db

import (
	"github.com/fredericobormann/dryck/models"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type DB struct {
	Datastore *gorm.DB
}

func New(dbType string, dbConnInfo string) *DB {
	db, err := gorm.Open(dbType, dbConnInfo)
	if err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Drink{})
	db.AutoMigrate(&models.Purchase{})
	db.AutoMigrate(&models.Payment{})

	dryckdb := DB{
		db,
	}

	return &dryckdb
}

func (db *DB) Close() {
	err := db.Datastore.Close()
	if err != nil {
		log.Print(err)
	}
}

// Returns all users currently in the database
func (db *DB) GetAllUsers() []models.User {
	var allUsers []models.User
	db.Datastore.Find(&allUsers)
	return allUsers
}

// Returns the username for a given user id
func (db *DB) GetUsername(userId uint) string {
	var user models.User
	db.Datastore.Where("id = ?", userId).First(&user)
	return user.Name
}

// Creates a new user with given username
func (db *DB) CreateNewUser(username string) {
	newUser := models.User{Name: username}
	db.Datastore.Create(&newUser)
}

// Returns all drinks in the database
func (db *DB) GetAllDrinks() []models.Drink {
	var allDrinks []models.Drink
	db.Datastore.Find(&allDrinks)
	return allDrinks
}

// Return all purchases of one user specified by user id with drink information preloaded
func (db *DB) GetPurchasesOfUser(userId uint) []models.Purchase {
	var purchases []models.Purchase
	db.Datastore.Preload("Product").Where("customer_id = ?", userId).Order("purchase_time desc").Find(&purchases)
	return purchases
}

// Returns the total debt of one user specified by their user id
func (db *DB) GetTotalDebtOfUser(userId uint) int {
	var totalDebt int
	db.Datastore.Table("purchases").Where("customer_id = ?", userId).Joins("inner join drinks on purchases.product_id = drinks.id").Select("sum(drinks.price)").Row().Scan(&totalDebt)
	var totalPayed int
	db.Datastore.Table("payments").Where("user_id = ?", userId).Select("sum(amount)").Row().Scan(&totalPayed)
	return totalDebt - totalPayed
}

// Adds a purchase for one user specified by their id and a drink also specified by id
func (db *DB) PurchaseDrink(userId uint, drinkId uint) {
	purchase := models.Purchase{CustomerID: userId, ProductID: drinkId, PurchaseTime: time.Now()}
	db.Datastore.Create(&purchase)
}

// Delete a purchase from database
func (db *DB) DeletePurchase(purchaseId uint) {
	db.Datastore.Where("id = ?", purchaseId).Unscoped().Delete(&models.Purchase{})
}

// Adds a payment for a user specified by id with a certain amount
func (db *DB) AddPayment(userId uint, amount int) {
	payment := models.Payment{UserID: userId, Amount: amount, PaymentTime: time.Now()}
	db.Datastore.Create(&payment)
}

func (db *DB) DeletePayment(paymentId uint) {
	db.Datastore.Where("id = ?", paymentId).Unscoped().Delete(&models.Payment{})
}

// Returns all payments of one user specified by id
func (db *DB) GetAllPaymentsOfUser(userId uint) []models.Payment {
	var payments []models.Payment
	db.Datastore.Where("user_id = ?", userId).Find(&payments)
	return payments
}
