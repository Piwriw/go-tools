package main

import (
	"github.com/piwriw/gorm/model"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Querier interface {
	// SELECT * FROM @@table WHERE employee_id=@id
	GetByID(id int) (gen.T, error) // GetByID query data by id and return it as *struct*

	// SELECT * FROM @@table WHERE first_name=@name
	GetByName(name string) (gen.T, error)

	//  SELECT * FROM @@table WHERE salary<@salary
	GetEmps(salary int) ([]gen.T, error)

	// SELECT e.employee_id, e.first_name, e.last_name, e.department_id, s.salary_amount
	//FROM @@table e
	//LEFT JOIN Salary s ON e.employee_id = s.employee_id;
	GetSalary() ([]gen.M, error)

	// SELECT * FROM Employees WHERE first_name REGEXP '^J'
	FitEmps() ([]gen.T, error)
}

func GenSql(db *gorm.DB) {

	g := gen.NewGenerator(gen.Config{
		OutPath: "./dao",
	})
	g.UseDB(db)
	//list, err := db.Migrator().GetTables()
	//if err != nil {
	//	panic(err)
	//}
	// export all tables Models
	//for _, s := range list {
	//	model := g.GenerateModel(s)
	//	g.ApplyInterface(func(Querier) {}, model)
	//}
	g.ApplyInterface(func(Querier) {}, model.HestiaInstance{})
	g.Execute()
}
func main() {
	//dsn := "root:123456@tcp(10.0.0.197:3306)/joohwan_dev?parseTime=true"
	//
	//mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	pgDb, err := gorm.Open(postgres.Open("postgres://yunqu:YunquTech01*@@10.0.0.195:15432/devops"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	GenSql(pgDb)
}
