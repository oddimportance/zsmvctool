package lib

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func (d *DbAdapter) Update(fields []string, values []string) *DbAdapter {
	d.prepareUpdateQuery(fields, values)
	//	fmt.Println(d.queryString)
	return d
}

func (d *DbAdapter) prepareUpdateQuery(fields []string, values []string) {
	setString := ""
	for i, field := range fields {
		setString = fmt.Sprintf("%s `%s` = '%s', ", setString, field, StringReplace(EscapeHtml(values[i]), "'", "\\'", -1))
	}
	setString = strings.TrimRight(strings.TrimSpace(setString), ",")

	d.queryString = fmt.Sprintf("UPDATE %s SET %s", d.dbTable, setString)
}

// Executes the sql query
// @ return map[int]map[string]string
func (d *DbAdapter) ExecUpdate() int {

	if !d.isConnectedToServer {
		return 0
	}

	var affectedRows int

	// prepare the joins
	d.prepareJoin(d.rightJoinClause)

	//	fmt.Println(d.queryString)
	d.prepareWhereForSelect()

	d.prepareQueryWithGroupBy()

	d.prepareQueryWithOrderBy()

	d.prepareQueryWithLimit(false)

	//	fmt.Printf("query: %s\n", d.queryString)
	if d.queryHasPotentialThreat {
		log.Println("there seems to be a threat in the query")
		log.Println(d.queryString)
		//unset the query
		d.unsetGlobalQueryVars()
		return -1
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
			d._handleRedirectAndPanic.RedirectToPanicUrl()
		}

		//	fmt.Printf("Res %v", rows)

		return affectedRows
	}

}
