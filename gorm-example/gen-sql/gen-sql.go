package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Querier interface {
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id int) (gen.T, error) // GetByID query data by id and return it as *struct*

	// SELECT * FROM @@table WHERE name=@name
	GetByName(name string) (gen.T, error)

	//  SELECT * FROM @@table WHERE age<@age
	GetEmps(age int) ([]*gen.T, error)
}

func GenSql(db *gorm.DB) {

	g := gen.NewGenerator(gen.Config{
		OutPath: "./dao",
	})
	g.UseDB(db)
	// Apply the interface to existing `User` and generated `Employee`
	g.ApplyInterface(func(Querier) {}, g.GenerateModel("emp"))

	g.Execute()
}
func main() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:abc123@192.168.113.1:55433/postgres"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	GenSql(db)
}
