package lib

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (d *DbAdapter) Delete() *DbAdapter {
	d.prepareDeleteQuery()
	//	fmt.Println(d.queryString)
	return d
}

func (d *DbAdapter) prepareDeleteQuery() {
	d.queryString = fmt.Sprintf("DELETE FROM %s", d.dbTable)
}

// Executes the sql query
// @ return map[int]map[string]string
func (d *DbAdapter) ExecDelete() int {

	if !d.isConnectedToServer {
		return 0
	}

	var affectedRows int = -1

	// prepare the joins
	d.prepareJoin(d.rightJoinClause)

	//	fmt.Println(d.queryString)
	d.prepareWhereForSelect()

	d.prepareQueryWithGroupBy()

	d.prepareQueryWithOrderBy()

	d.prepareQueryWithLimit(false)

	//	fmt.Printf("query: %s\n", d.queryString)

	// fmt.Printf("query: %s\n", d.queryString)
	if d.queryHasPotentialThreat {
		log.Println("there seems to be a threat in the query")
		log.Println(d.queryString)
		//unset the query
		d.unsetGlobalQueryVars()
	} else {

		d.setLastExectuedQuery(d.queryString)

		rows, dbQueryError := d._db.Query(d.queryString)

		//fmt.Println(reflect.TypeOf(rows))

		//unset the query
		d.unsetGlobalQueryVars()

		// prepare affected rows to return
		rows.Scan(&affectedRows)

		// free memory
		defer rows.Close()

		// To avoid SQL sneaks, look if there
		// were any mysql query errors
		if nil == rows || dbQueryError != nil {
			d._handleRedirectAndPanic.TriggerPanic("0003")
		}
	}
	//	fmt.Printf("Res %v", rows)

	return affectedRows

}
