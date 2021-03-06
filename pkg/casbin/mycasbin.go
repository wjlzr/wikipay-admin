package mycasbin

import (
	"fmt"
	"wikipay-admin/database"
	"wikipay-admin/tools/config"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/go-kit/kit/endpoint"
	_ "github.com/go-sql-driver/mysql"
)

var Em endpoint.Middleware

func Casbin() (*casbin.Enforcer, error) {
	conn := database.GetMysqlConnect()
	if config.DatabaseConfig.Dbtype == "sqlite3" {
		conn = config.DatabaseConfig.Host
	}
	Apter, err := gormadapter.NewAdapter(config.DatabaseConfig.Dbtype, conn, true)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer("config/rbac_model.conf", Apter)
	if err != nil {
		return nil, err
	}
	if err := e.LoadPolicy(); err == nil {
		return e, err
	} else {
		fmt.Print("casbin rbac_model or policy init error, message: %v", err)
		return nil, err
	}
}
