package main

import (
	//"database/sql"
	"log"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/db/types"
	"github.com/go-jet/jet/v2/generator/metadata"
	genpostgres "github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	"github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
)

var schemas []string = []string{
	"arcgis",
	"public",
	"publicreport",
	"stadia",
}

func customTemplate() template.Template {
	return template.Default(postgres.Dialect).UseSchema(func(schema metadata.Schema) template.Schema {
		return template.DefaultSchema(schema).UseModel(template.DefaultModel().UseTable(func(table metadata.Table) template.TableModel {
			return template.DefaultTableModel(table).UseField(func(column metadata.Column) template.TableModelField {
				defaultTableModelField := template.DefaultTableModelField(column)
				//log.Printf("'%s' '%s' '%s'", table.Name, column.Name, column.DataType.Name)
				if column.Name == "extent" && column.DataType.Name == "box2d" {
					defaultTableModelField.Type = template.NewType(types.Box2D{})
				}
				return defaultTableModelField
			})
		}),
		)
	})
}

func main() {
	for _, schema := range schemas {
		err := genpostgres.GenerateDSN(
			"postgresql://?host=/var/run/postgresql&sslmode=disable&dbname=nidus-sync",
			schema,
			"../gen",
			customTemplate(),
		)
		if err != nil {
			log.Printf("Failed: %v", err)
			os.Exit(1)
		}
	}
}
