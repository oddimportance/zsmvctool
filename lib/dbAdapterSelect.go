package lib

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Executes the sql query
// @ return map[int]map[string]string
func (d *DbAdapter) ExecSelect() map[int]map[string]string {

	placeHolder := map[int]map[string]string{}

	if !d.isConnectedToServer {
		return placeHolder
	}

	// prepare the joins
	d.prepareJoin(d.rightJoinClause)
	d.prepareJoin(d.leftJoinClause)

	// fmt.Println(d.queryString)
	d.prepareWhereForSelect()

	d.prepareQueryWithGroupBy()

	d.prepareQueryWithOrderBy()

	d.prepareQueryWithLimit(true)

	//	log.Println("Calling Select")

	// prepare additional select as fields
	// ex.: COUNT(zsartcmt_comments) AS article_comments
	d.concatAdditionalSelectAsFields()

	// fmt.Printf("query: %s\n", d.queryString)
	if d.queryHasPotentialThreat {
		log.Println("there seems to be a threat in the query")
		log.Println(d.queryString)
		//unset the query
		d.unsetGlobalQueryVars()
		return nil
	} else {

		d.setLastExectuedQuery(d.queryString)

		rows, dbQueryError := d._db.Query(d.queryString)

		//unset the query
		d.unsetGlobalQueryVars()

		//	fmt.Println(reflect.TypeOf(rows))

		// To avoid SQL sneaks, look if there
		// were any mysql query errors
		if nil == rows || dbQueryError != nil {

			// Debug
			//		fmt.Println(dbQueryError)
			//		fmt.Println(d.GetLastExecutedQuery())

			log.Println(dbQueryError)
			log.Println(d.GetLastExecutedQuery())

			d._handleRedirectAndPanic.TriggerPanic("0003")
			d._handleRedirectAndPanic.RedirectToPanicUrl()
			return placeHolder
		}

		// free memory
		defer rows.Close()

		placeHolder = d.dbResultToReadable(rows)

		return placeHolder
		//fmt.Printf("Res %v", readableDbResult)
	}
}

func (d *DbAdapter) unsetGlobalQueryVars() {

	d.queryString = ""
	d.whereClause = nil
	d.whereOrClause = nil
	d.rightJoinClause = nil
	d.leftJoinClause = nil
	d.additionalSelectAsFields = nil
	d.groupByField = nil
	d.orderByClause = nil

}

func (d *DbAdapter) concatAdditionalSelectAsFields() {
	if len(d.additionalSelectAsFields) != 0 {
		d.queryString = StringReplace(d.queryString,
			"*",
			fmt.Sprintf("*, %s ", strings.Join(d.additionalSelectAsFields, ", ")),
			1)
	}
}

func (d *DbAdapter) dbResultToReadable(dbResult *sql.Rows) map[int]map[string]string {

	final_result := map[int]map[string]string{}

	columns, dbFetchError := dbResult.Columns()

	if dbFetchError != nil {
		d._handleRedirectAndPanic.TriggerPanic("0004")
		d._handleRedirectAndPanic.RedirectToPanicUrl()
		return final_result
	}

	count := len(columns)

	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	result_id := 0
	for dbResult.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		dbResult.Scan(valuePtrs...)

		tmp_struct := map[string]string{}

		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			tmp_struct[col] = fmt.Sprintf("%s", v)
		}

		final_result[result_id] = tmp_struct
		result_id++
	}

	return final_result
}

func (d *DbAdapter) SelectAll() *DbAdapter {

	d.queryHasPotentialThreat = false
	d.prepareSelectStatement(make([]string, 0))
	return d

}

func (d *DbAdapter) SelectOne() *DbAdapter {
	d.queryHasPotentialThreat = false
	d.SetLimit(0, 1)
	d.prepareSelectStatement(make([]string, 0))
	return d
}

func (d *DbAdapter) SelectFields(fields []string) *DbAdapter {
	d.queryHasPotentialThreat = false
	d.prepareSelectStatement(fields)
	return d
}

// Select fields with AS ex. SELECT name AS test
// @ param fields map[string]string
// @ return *DbAdapter
func (d *DbAdapter) SelectFieldsAs(fields map[string]string) *DbAdapter {
	d.queryHasPotentialThreat = false

	fieldsToString := ""

	for field, fieldAs := range fields {
		fieldsToString = fmt.Sprintf("%s %s AS `%s`,", fieldsToString, field, fieldAs)
	}

	fieldsToString = strings.TrimRight(strings.TrimSpace(fieldsToString), ",")

	d.queryString = fmt.Sprintf("SELECT %s FROM `%s`", fieldsToString, d.dbTable)

	return d
}

// Select count
// @ param fields []string can be nil
// @ return *DbAdapter
func (d *DbAdapter) SelectCount(fields []string) *DbAdapter {
	d.queryHasPotentialThreat = false

	fieldsToString := ""

	if nil != fields {
		for _, field := range fields {
			fieldsToString = fmt.Sprintf("`%s`, %s", field, fieldsToString)
		}
		fieldsToString = fmt.Sprintf(", %s", fieldsToString)
	}

	fieldsToString = strings.TrimRight(strings.TrimSpace(fieldsToString), ",")

	d.queryString = fmt.Sprintf("SELECT COUNT(*) AS count %s FROM `%s`", fieldsToString, d.dbTable)

	return d

}

func (d *DbAdapter) prepareSelectStatement(fields []string) {
	if len(fields) == 0 {
		d.queryString = fmt.Sprintf("SELECT * FROM `%s`", d.dbTable)
	} else {
		d.queryString = fmt.Sprintf("SELECT %s FROM `%s`", d.prepareFieldsForQueryStatement(fields, "`", ",", true), d.dbTable)
	}
}

func (d *DbAdapter) SetAdditionalSelectAsField(additionalSelectAsField string) {
	d.additionalSelectAsFields = append(d.additionalSelectAsFields, additionalSelectAsField)
}

func (d *DbAdapter) RightJoin(table string, fieldOn string, fieldTo string) *DbAdapter {
	d.rightJoinClause = append(d.rightJoinClause, fmt.Sprintf("RIGHT JOIN %s ON %s = %s", table, fieldOn, fieldTo))
	return d
}

func (d *DbAdapter) LeftJoin(table string, fieldOn string, fieldTo string) *DbAdapter {
	d.leftJoinClause = append(d.leftJoinClause, fmt.Sprintf("LEFT JOIN %s ON %s = %s", table, fieldOn, fieldTo))
	return d
}

func (d *DbAdapter) prepareJoin(joins []string) {
	for _, join := range joins {
		d.queryString = fmt.Sprintf("%s %s", d.queryString, join)
	}
}

func (d *DbAdapter) MakeAsField(field, fieldAs string) string {
	// quote only the field as and not the field as field may
	// contain mysql functions
	return fmt.Sprintf("%s AS `%s`", field, fieldAs)
}

func (d *DbAdapter) MakeMysqlFunction(mysql_function, field string) string {
	d.initMysqlFunctionKey()
	// fmt.Printf("%s%s(%s)\n", d.mysqlFunctionKey, StringToUpper(mysql_function), field)
	return fmt.Sprintf("%s%s(%s)", d.mysqlFunctionKey, StringToUpper(mysql_function), field)
}

func (d *DbAdapter) initMysqlFunctionKey() {
	if d.mysqlFunctionKey == "" {
		d.mysqlFunctionKey = "MYSQL_FUNC_"
	}
}
