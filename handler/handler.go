package handler

import (
	"fmt"
	"github.com/fredericobormann/dryck/db"
	"github.com/fredericobormann/dryck/format"
	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
	"math"
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
			"csrftoken": csrf.GetToken(c),
			"title":     "Home Page",
			"users":     h.Datastore.GetAllUsers(),
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
	itemsPerPurchasePage := 10
	var purchasePaginator []Pagebutton
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)

	errorMessage := ""
	if errorType := c.Query("error"); errorType != "" {
		if errorType == "wrong_format" {
			errorMessage = "Der Betrag wurde nicht im richtigen Format eingegeben. Bitte verwende \"1,50\" um eine Zahlung von 1,50€ einzutragen."
		}
	}

	numberOfPurchases := h.Datastore.GetNumberOfPurchasesOfUser(uint(userID))
	numberOfPurchasePages := 1
	if numberOfPurchases > 0 {
		numberOfPurchasePages = int(math.Ceil(float64(numberOfPurchases) / float64(itemsPerPurchasePage)))
	}

	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || int(page) > numberOfPurchasePages || int(page) < 1 {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if numberOfPurchases > itemsPerPurchasePage {
		purchasePaginator = newPagebuttons(int(page), numberOfPurchasePages)
	}

	username := h.Datastore.GetUsername(uint(userID))
	totalDebt := h.Datastore.GetTotalDebtOfUser(uint(userID))
	purchases := h.Datastore.GetPaginatedPurchasesOfUser(uint(userID), itemsPerPurchasePage, int(page))
	payments := h.Datastore.GetAllPaymentsOfUser(uint(userID))
	drinks := h.Datastore.GetAllDrinks()

	c.HTML(
		http.StatusOK,
		"user.html",
		gin.H{
			"csrftoken":         csrf.GetToken(c),
			"title":             username,
			"username":          username,
			"userID":            userID,
			"totalDebt":         totalDebt,
			"drinks":            drinks,
			"purchases":         purchases,
			"purchasePaginator": purchasePaginator,
			"payments":          payments,
			"errorMessage":      errorMessage,
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
	paymentAmount, err := format.FromPrice(c.PostForm("payment_amount"))

	if err != nil {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/user/%d?error=wrong_format", userID))
		return
	}

	h.Datastore.AddPayment(uint(userID), int(paymentAmount))
	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// HandleDeletePayment handles the deletion of a payment
func (h *Handler) HandleDeletePayment(c *gin.Context) {
	paymentID, _ := strconv.ParseUint(c.PostForm("delete_payment"), 10, 64)
	h.Datastore.DeletePayment(uint(paymentID))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// HandleLoginPage handles the login page
func (h *Handler) HandleLoginPage(c *gin.Context) {
	message := c.Query("message")
	c.HTML(
		http.StatusOK,
		"login.html",
		gin.H{
			"csrftoken":  csrf.GetToken(c),
			"hasMessage": message != "",
			"message":    message,
		},
	)
}
