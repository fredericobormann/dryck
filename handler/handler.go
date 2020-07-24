package handler

import (
	"github.com/fredericobormann/dryck/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Handler stores necessary information to handle requests
type Handler struct {
	Datastore *db.DB
}

// New creates a new Postgres Datastore
func New(datastore *db.DB) *Handler {
	return &Handler{
		datastore,
	}
}

// HandleIndex handles requests to the index page
func (h *Handler) HandleIndex(c *gin.Context) {
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title": "Home Page",
			"users": h.Datastore.GetAllUsers(),
		},
	)
}

// HandleNewUser handles creation of a new user
func (h *Handler) HandleNewUser(c *gin.Context) {
	newUserName := c.PostForm("new-user-name")
	h.Datastore.CreateNewUser(newUserName)

	c.Redirect(http.StatusMovedPermanently, "/")
}

// HandleUserPage handles user page with purchase history and the option to buy new drinks
func (h *Handler) HandleUserPage(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	username := h.Datastore.GetUsername(uint(userID))
	totalDebt := h.Datastore.GetTotalDebtOfUser(uint(userID))
	purchases := h.Datastore.GetPurchasesOfUser(uint(userID))
	payments := h.Datastore.GetAllPaymentsOfUser(uint(userID))
	drinks := h.Datastore.GetAllDrinks()

	c.HTML(
		http.StatusOK,
		"user.html",
		gin.H{
			"title":     username,
			"username":  username,
			"userID":    userID,
			"totalDebt": totalDebt,
			"drinks":    drinks,
			"purchases": purchases,
			"payments":  payments,
		},
	)
}

// HandlePurchase handles a new purchase
func (h *Handler) HandlePurchase(c *gin.Context) {
	drinkID, _ := strconv.ParseUint(c.PostForm("drink"), 10, 64)
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	h.Datastore.PurchaseDrink(uint(userID), uint(drinkID))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// HandleDeletePurchase handles the deletion of a purchase
func (h *Handler) HandleDeletePurchase(c *gin.Context) {
	purchaseID, _ := strconv.ParseUint(c.PostForm("delete_purchase"), 10, 64)
	h.Datastore.DeletePurchase(uint(purchaseID))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// HandlePayment handles a new payment
func (h *Handler) HandlePayment(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	paymentAmount, _ := strconv.ParseInt(c.PostForm("payment_amount"), 10, 64)

	h.Datastore.AddPayment(uint(userID), int(paymentAmount))
	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// HandleDeletePayment handles the deletion of a payment
func (h *Handler) HandleDeletePayment(c *gin.Context) {
	paymentID, _ := strconv.ParseUint(c.PostForm("delete_payment"), 10, 64)
	h.Datastore.DeletePayment(uint(paymentID))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}
