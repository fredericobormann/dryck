package main

import (
	"fmt"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/fredericobormann/dryck/db"
	"github.com/fredericobormann/dryck/format"
	"github.com/fredericobormann/dryck/handler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/utrack/gin-csrf"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func main() {
	databaseHost := os.Getenv("POSTGRES_HOST")
	databaseUser := os.Getenv("POSTGRES_USER")
	databasePassword := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DATABASE")
	dryckdb := db.New("postgres", "host="+databaseHost+" user="+databaseUser+" dbname="+databaseName+" password="+databasePassword+" sslmode=disable")
	_, httpBasicAuthActive := os.LookupEnv("HTTP_PASSWORD")
	jwtSecret, jwtSecretDefined := os.LookupEnv("JWT_SECRET")

	defer dryckdb.Close()

	dryckhandler := handler.New(dryckdb)

	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsPrice": format.AsPrice,
		"formatAsTime":  format.AsTime,
	})
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	store := cookie.NewStore([]byte(jwtSecret))
	router.Use(sessions.Sessions("drycksession", store))
	router.Use(csrf.Middleware(csrf.Options{
		Secret: jwtSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	var authorized *gin.RouterGroup

	if httpBasicAuthActive {
		if !jwtSecretDefined {
			log.Fatal("JWT Secret is not defined.")
		}
		jwtMiddleware, jwtErr := createJWTMiddleware(jwtSecret)
		if jwtErr != nil {
			log.Fatal(jwtErr)
		}
		router.GET("/login", dryckhandler.HandleLoginPage)
		router.POST("/login", jwtMiddleware.LoginHandler)

		authorized = router.Group("/", jwtMiddleware.MiddlewareFunc(), jwtMiddleware.RefreshHandler)
	} else {
		authorized = router.Group("/")
	}

	authorized.GET("/", dryckhandler.HandleIndex)
	authorized.POST("/new-user", dryckhandler.HandleNewUser)
	authorized.GET("/user/:user_id", dryckhandler.HandleUserPage)
	authorized.POST("/purchase/:user_id", dryckhandler.HandlePurchase)
	authorized.POST("/delete-purchase/:user_id", dryckhandler.HandleDeletePurchase)

	authorized.POST("/new-payment/:user_id", dryckhandler.HandlePayment)
	authorized.POST("/delete-payment/:user_id", dryckhandler.HandleDeletePayment)

	router.Run()
}

func createJWTMiddleware(jwtSecret string) (authMiddleware *jwt.GinJWTMiddleware, err error) {
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "authenticated zone",
		Key:         []byte(jwtSecret),
		Timeout:     30 * 24 * time.Hour,
		MaxRefresh:  30 * 24 * time.Hour,
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if _, ok := data.(bool); ok {
				return jwt.MapClaims{
					"id": "0",
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims["id"].(string) == "0"
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			httpPassword, _ := os.LookupEnv("HTTP_PASSWORD")

			if loginVals.Username == "dryck" && loginVals.Password == httpPassword {
				return true, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(bool); ok && v {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			if message == "cookie token is empty" {
				c.Redirect(http.StatusFound, "/login")
			} else {
				c.Redirect(http.StatusFound, fmt.Sprintf("/login?message=%s", message))
			}
		},
		LoginResponse: func(c *gin.Context, i int, s string, t time.Time) {
			c.Redirect(http.StatusFound, "/")
		},
		RefreshResponse: func(c *gin.Context, code int, message string, i time.Time) {
			c.Next()
		},
		TokenLookup:    "header: Authorization, query: token, cookie: jwt",
		CookieSameSite: http.SameSiteStrictMode,
		SendCookie:     true,
		SecureCookie: func() bool {
			ginMode, ok := os.LookupEnv("GIN_MODE")
			return ok && ginMode == "release"
		}(),
		CookieHTTPOnly: true, // JS can't modify

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	return
}
