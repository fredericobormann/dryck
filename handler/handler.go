package handler

import (
	"github.com/fredericobormann/dryck/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	Datastore *db.DB
}

//Creates a new Handler
func New(datastore *db.DB) *Handler {
	return &Handler{
		datastore,
	}
}

// Handles requests to the index page
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

// Handles creation of a new user
func (h *Handler) HandleNewUser(c *gin.Context) {
	newUserName := c.PostForm("new-user-name")
	h.Datastore.CreateNewUser(newUserName)

	c.Redirect(http.StatusMovedPermanently, "/")
}

// Handles user page with purchase history and the option to buy new drinks
func (h *Handler) HandleUserPage(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	username := h.Datastore.GetUsername(uint(userId))
	totalDebt := h.Datastore.GetTotalDebtOfUser(uint(userId))
	purchases := h.Datastore.GetPurchasesOfUser(uint(userId))
	payments := h.Datastore.GetAllPaymentsOfUser(uint(userId))
	drinks := h.Datastore.GetAllDrinks()

	c.HTML(
		http.StatusOK,
		"user.html",
		gin.H{
			"title":     username,
			"username":  username,
			"userId":    userId,
			"totalDebt": totalDebt,
			"drinks":    drinks,
			"purchases": purchases,
			"payments":  payments,
		},
	)
}

// Handles a new purchase
func (h *Handler) HandlePurchase(c *gin.Context) {
	drinkId, _ := strconv.ParseUint(c.PostForm("drink"), 10, 64)
	userId, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	h.Datastore.PurchaseDrink(uint(userId), uint(drinkId))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// Handles the deletion of a purchase
func (h *Handler) HandleDeletePurchase(c *gin.Context) {
	purchaseId, _ := strconv.ParseUint(c.PostForm("delete_purchase"), 10, 64)
	h.Datastore.DeletePurchase(uint(purchaseId))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// Handles a new payment
func (h *Handler) HandlePayment(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	paymentAmount, _ := strconv.ParseInt(c.PostForm("payment_amount"), 10, 64)

	h.Datastore.AddPayment(uint(userId), int(paymentAmount))
	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// Handles the deletion of a payment
func (h *Handler) HandleDeletePayment(c *gin.Context) {
	paymentId, _ := strconv.ParseUint(c.PostForm("delete_payment"), 10, 64)
	h.Datastore.DeletePayment(uint(paymentId))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}
