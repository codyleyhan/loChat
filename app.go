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
	defer db.Close()

	db.LogMode(true)

	// Needed for dealing with location
	db.Exec("CREATE EXTENSION IF NOT EXISTS cube;")
	db.Exec("CREATE EXTENSION IF NOT EXISTS earthdistance;")

	jwtNumDays := conf.GetInt("jwt.days_valid_for")
	jwtSecret := conf.Get("jwt.secret")

	jwt := utils.JWTGen{
		NumDaysValid: int(jwtNumDays),
		Secret:       []byte(jwtSecret),
	}

	app := echo.New()

	utils.AddMiddleware(app, logo)

	userStore := stores.CreateUserStore(db)
	roomStore := stores.CreateRoomStore(db)

	services.RegisterServices(app, userStore, roomStore, &jwt)

	port := ":" + conf.Get("port")
	fmt.Println("Starting server on port:", conf.Get("port"))
	app.Logger.Fatal(app.Start(port))
}
