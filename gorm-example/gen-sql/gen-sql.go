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
	list, err := db.Migrator().GetTables()
	if err != nil {
		panic(err)
	}
	// export all tables Models
	for _, s := range list {
		model := g.GenerateModel(s)
		g.ApplyInterface(func(Querier) {}, model)
	}

	g.Execute()
}
func main() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:abc123@xxx.xxx.xxx.xx:55433/postgres"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	GenSql(db)
}
