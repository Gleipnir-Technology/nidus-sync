package main

import (
	//"database/sql"
	"log"
	"os"

	"github.com/Gleipnir-Technology/jet/generator/metadata"
	genpostgres "github.com/Gleipnir-Technology/jet/generator/postgres"
	"github.com/Gleipnir-Technology/jet/generator/template"
	"github.com/Gleipnir-Technology/jet/postgres"
	_ "github.com/lib/pq"
	"github.com/twpayne/go-geom"
)

var schemas []string = []string{
	"arcgis",
	"comms",
	"public",
	"publicreport",
	"stadia",
}

func customTableSQLBuilderColumn(dialect template.Dialect, column metadata.Column) template.TableSQLBuilderColumn {
	defaultColumn := template.DefaultTableSQLBuilderColumn(dialect, column)
	/*
		if defaultColumn.Name == "Location" {
			log.Printf("current location column: name '%s' type '%s'", defaultColumn.Name, defaultColumn.Type)
			defaultColumn.Import = "github.com/Gleipnir-Technology/nidus-sync/db/column"
			defaultColumn.PackageName = "column"
			defaultColumn.Type = "ColumnGeometry"
			defaultColumn.TypeFactory = "GeometryColumn"
		}
	*/
	return defaultColumn
}
func customTableSQLBuilder(table metadata.Table) template.TableSQLBuilder {
	builder := template.DefaultTableSQLBuilder(table).UseColumn(customTableSQLBuilderColumn)
	log.Printf("table sql builder: path '%s' filename '%s' instancename '%s' typename '%s' defaultalias '%s'", builder.Path, builder.FileName, builder.InstanceName, builder.TypeName, builder.DefaultAlias)
	builder.Imports = []string{"github.com/Gleipnir-Technology/nidus-sync/db/column"}
	return builder
}
func customTableModelField(column metadata.Column) template.TableModelField {
	defaultTableModelField := template.DefaultTableModelField(column)
	//log.Printf("'%s' '%s' '%s'", table.Name, column.Name, column.DataType.Name)
	if column.Name == "extent" && column.DataType.Name == "box2d" {
		defaultTableModelField.Type = template.NewType(geom.Bounds{})
	} else if column.DataType.Name == "geometry" {
		name := "geom.T"
		if column.IsNullable {
			name = "*geom.T"
		}
		geom_type := template.Type{
			ImportPath:            "github.com/twpayne/go-geom",
			AdditionalImportPaths: []string{},
			Name:                  name,
		}
		defaultTableModelField.Type = geom_type
	}
	return defaultTableModelField
}
func customTableModel(table metadata.Table) template.TableModel {
	return template.DefaultTableModel(table).UseField(customTableModelField)
}
func customTemplate() template.Template {
	return template.Default(postgres.Dialect).UseSchema(func(schema metadata.Schema) template.Schema {
		customSchema := template.DefaultSchema(schema)
		customSchema = customSchema.UseModel(template.DefaultModel().UseTable(customTableModel))
		customSchema = customSchema.UseSQLBuilder(template.DefaultSQLBuilder().UseTable(customTableSQLBuilder))
		return customSchema
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
