package migrations

import (
	"testing"

	"github.com/go-xorm/xorm"

	. "github.com/smartystreets/goconvey/convey"
)

// func cleanDB(x *xorm.Engine) {
// 	tables, _ := x.DBMetas()
// 	sess := x.NewSession()
// 	defer sess.Close()
//
// 	for _, table := range tables {
// 		if _, err := sess.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
// 			panic("Failed to disable foreign key checks")
// 		}
// 		if _, err := sess.Exec("DROP TABLE " + table.Name); err != nil {
// 			panic(fmt.Sprintf("Failed to delete table: %v, err: %v", table.Name, err))
// 		}
// 		if _, err := sess.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
// 			panic("Failed to disable foreign key checks")
// 		}
// 	}
// }
//
// var indexTypes = []string{"Unknown", "", "UNIQUE"}
//

func TestMigrator(t *testing.T) {

	Convey("Migrator", t, func() {
		x, err := xorm.NewEngine(SQLITE, ":memory:")
		So(err, ShouldBeNil)

		mg := NewMigrator(x)

		Convey("Given one migration", func() {
			mg.AddMigration("test migration", new(RawSqlMigration).
				Sqlite(`
			    CREATE TABLE account (
		       	id INTEGER PRIMARY KEY AUTOINCREMENT
				  )`).
				Mysql(`
			   	CREATE TABLE account (
						id BIGINT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id)
					)`))

			err := mg.Start()
			So(err, ShouldBeNil)

			log, err := mg.GetMigrationLog()
			So(err, ShouldBeNil)
			So(len(log), ShouldEqual, 1)
		})

		// So(err, ShouldBeNil)
		//
		// So(len(tables), ShouldEqual, 2)
		// fmt.Printf("\nDB Schema after migration: table count: %v\n", len(tables))
		//
		// for _, table := range tables {
		// 	fmt.Printf("\nTable: %v \n", table.Name)
		// 	for _, column := range table.Columns() {
		// 		fmt.Printf("\t %v \n", column.String(x.Dialect()))
		// 	}
		//
		// 	if len(table.Indexes) > 0 {
		// 		fmt.Printf("\n\tIndexes:\n")
		// 		for _, index := range table.Indexes {
		// 			fmt.Printf("\t %v (%v) %v \n", index.Name, strings.Join(index.Cols, ","), indexTypes[index.Type])
		// 		}
		// 	}
		// }
	})
}
