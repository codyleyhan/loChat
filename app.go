package loChat

import (
	"fmt"

	"github.com/codyleyhan/loChat/services"
	"github.com/codyleyhan/loChat/stores"
	"github.com/codyleyhan/loChat/utils"

	"github.com/codyleyhan/config-loader"
	"github.com/labstack/echo"
)

const logo = `
    _/              _/_/_/  _/                    _/      
   _/    _/_/    _/        _/_/_/      _/_/_/  _/_/_/_/   
  _/  _/    _/  _/        _/    _/  _/    _/    _/        
 _/  _/    _/  _/        _/    _/  _/    _/    _/         
_/    _/_/      _/_/_/  _/    _/    _/_/_/      _/_/                                                     
`

//StartApp call this to start the application
func StartApp() {
	conf := config.LoadConfig("./config", "config")

	dbOpts := utils.DbOptions{
		User:     conf.Get("db.user"),
		Password: conf.Get("db.password"),
		Name:     conf.Get("db.name"),
		Port:     conf.Get("db.port"),
	}
	db := utils.ConnectToDB(dbOpts)
	db.LogMode(true)
	defer db.Close()

	jwtNumDays := conf.GetInt("jwt.days_valid_for")
	jwtSecret := conf.Get("jwt.secret")

	jwt := utils.JWTGen{
		NumDaysValid: int(jwtNumDays),
		Secret:       []byte(jwtSecret),
	}

	app := echo.New()

	utils.AddMiddleware(app, logo)

	userStore := stores.CreateUserStore(db)

	services.RegisterServices(app, userStore, &jwt)

	// Example of a secure route
	app.GET("/secure", func(c echo.Context) error {
		user := c.Get("user")

		return c.JSON(200, map[string]interface{}{
			"Hello": "World",
			"user":  user,
		})
	}, utils.CheckToken(&jwt))

	// Example of unsecured
	app.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"Hello": "World",
		})
	})

	port := ":" + conf.Get("port")
	fmt.Println("Starting server on port:", conf.Get("port"))
	app.Logger.Fatal(app.Start(port))
}
