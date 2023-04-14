package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/oddimportance/zsmvctool/persistence"
)

type DbAdapter struct {
	_db                      *sql.DB
	isConnectedToServer      bool
	dbTable                  string
	dbTableFieldPrefix       string
	queryString              string
	lastExecutedQuery        string
	valueString              string
	additionalSelectAsFields []string
	whereClause              []string
	whereOrClause            []string
	groupByField             []string
	rightJoinClause          []string
	leftJoinClause           []string
	join                     string
	OrderBy                  persistence.OrderBy
	orderByClause            []string
	limitFrom                int
	limitTo                  int
	Wildcard                 persistence.Wildcard
	MathComparision          persistence.MathComparision
	_handleRedirectAndPanic  *HandleRedirectAndPanic
	dbCredentials            *Credentials
	mysqlFunctionKey         string
	queryHasPotentialThreat  bool
}

type Credentials struct {
	host     string
	port     string
	user     string
	password string
	database string
}

func (d *DbAdapter) Connect(dbTable string, dbTableFieldPrefix string, _handleRedirectAndPanic *HandleRedirectAndPanic) {

	d.SetHandleRedirectAndPanic(_handleRedirectAndPanic)

	d.dbCredentials = d.getMysqlConfig()
	// log.Println(d.dbCredentials)

	// Initialize connection string.
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?allowNativePasswords=true", d.dbCredentials.user, d.dbCredentials.password, d.dbCredentials.host, d.dbCredentials.port, d.dbCredentials.database)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Println("Sql Open error")
		db.Close()
		fmt.Println(err)
		d._handleRedirectAndPanic.TriggerPanic("0001")
		d.isConnectedToServer = false
		return
	}

	// fmt.Println(db)

	if err = db.Ping(); err != nil {
		log.Println("DB Ping error")
		db.Close()
		fmt.Println(err)
		d._handleRedirectAndPanic.TriggerPanic("0001")
		d.isConnectedToServer = false
		return
	}

	// fmt.Printf("%s:%s@%s/%s?charset=utf8\n", dbCredentials.user, dbCredentials.password, dbCredentials.host, dbCredentials.database)

	d._db = db

	d._db.SetMaxOpenConns(1)
	d._db.SetMaxIdleConns(0)
	d._db.SetConnMaxLifetime(time.Nanosecond)

	d.setDbTable(dbTable)

	d.setDbTableFieldPrefix(dbTableFieldPrefix)

	d.isConnectedToServer = true

	// fmt.Printf("DB type %v", reflect.TypeOf(d._db))
}

func (d *DbAdapter) PrintDBDetails() {
	fmt.Println("+++++++++++++ Database Details +++++++++++++")
	fmt.Println("============================================")
	fmt.Printf("Host: %s\n", d.dbCredentials.host)
	fmt.Printf("Port: %s\n", d.dbCredentials.port)
	fmt.Printf("User Name: %s\n", d.dbCredentials.user)
	fmt.Printf("Database Name: %s\n", d.dbCredentials.database)
	fmt.Println("============================================")
	fmt.Println()
}

func (d *DbAdapter) GetSqlConnection() *sql.DB {
	return d._db
}

func (d *DbAdapter) SetSqlConnection(db *sql.DB) {
	// Make sure connected to sql server
	if db == nil {
		d.isConnectedToServer = false
	} else {
		d.isConnectedToServer = true
	}
	d._db = db

}

func (d *DbAdapter) SetHandleRedirectAndPanic(_handleRedirectAndPanic *HandleRedirectAndPanic) {
	d._handleRedirectAndPanic = _handleRedirectAndPanic
}

// Set the table and its preifx, if you
// have to handle additional tables other
// than the main table.
// @ param dbTable string Name of the table
// @ param tableFieldPrefix string Table prefix
// @ return void
func (d *DbAdapter) SetTableAndPrefix(dbTable string, tableFieldPrefix string) {
	d.setDbTable(dbTable)
	d.setDbTableFieldPrefix(tableFieldPrefix)
}

func (d *DbAdapter) setDbTable(dbTable string) {
	d.dbTable = dbTable
}

func (d *DbAdapter) setDbTableFieldPrefix(dbTableFieldPrefix string) {
	d.dbTableFieldPrefix = dbTableFieldPrefix
}

// returns table name
func (d *DbAdapter) GetTableName() string {
	return d.dbTable
}

func (d *DbAdapter) GetDbTableFieldPrefix() string {
	return d.dbTableFieldPrefix
}

func (d *DbAdapter) getMysqlConfig() *Credentials {

	port := d.getEnvConf("MP_API_DATABASE_PORT")
	if port == "" {
		port = "3306"
	}

	var dbCredentials = new(Credentials)
	dbCredentials.host = d.getEnvConf("MP_API_DATABASE_HOST")
	dbCredentials.port = port
	dbCredentials.database = d.getEnvConf("MP_API_DATABASE_DB")
	dbCredentials.user = d.getEnvConf("MP_API_DATABASE_USERNAME")
	dbCredentials.password = d.getEnvConf("MP_API_DATABASE_PASSWORD")

	return dbCredentials

}

func (d *DbAdapter) getEnvConf(key string) string {

	envConf, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("DB Host Env Conf Var %s was not found\n", key)
		panic("DB Credentials Error")
	}
	return envConf
}

func (d *DbAdapter) setPanicError(panicCode string) {
	d._handleRedirectAndPanic.PanicCode = panicCode
}

func (d *DbAdapter) prepareFieldsForQueryStatement(fields []string, decorator string, seperator string, trimRight bool) string {
	d.initMysqlFunctionKey()
	var fieldsToString string
	for _, val := range fields {
		// trim any white space
		val = TrimWhiteSpaces(val)
		// fmt.Println(key, val)
		// decorator prefix value decorator seperator, ex,: 'zsusr_val',
		// Howerver, if the field is a COUNT(DISTINCT zsartvw_id) AS article_views
		// then escape the decorator
		if strings.Contains(val, d.mysqlFunctionKey) {
			val = StringReplace(val, d.mysqlFunctionKey, "", -1)
			decorator = ""
		} else if (strings.Contains(val, "as") || strings.Contains(val, "AS") || strings.Contains(val, "COUNT") || strings.Contains(val, "MONTH") || strings.Contains(val, "YEAR")) && decorator == "`" {
			decorator = ""
		}
		fieldsToString = fmt.Sprintf("%s %s%s%s%s",
			fieldsToString,
			decorator,
			EscapeHtml(StringReplace(val, "'", "\\'", -1)),
			decorator,
			seperator)
		//fmt.Println(val)
	}
	if trimRight {
		return strings.TrimRight(strings.TrimSpace(fieldsToString), seperator)
	} else {
		return strings.TrimSpace(fieldsToString)
	}
}

// Set a particular field as Distinct by assigining
// it to GROUP BY DELEMETRI
// @ param field string
func (d *DbAdapter) SetDistinctGroupByField(field string) {
	d.groupByField = append(d.groupByField, field)
}

func (d *DbAdapter) prepareQueryWithGroupBy() {
	if len(d.groupByField) != 0 {
		d.queryString = fmt.Sprintf("%s GROUP BY %s", d.queryString, strings.Join(d.groupByField, ", "))
	}
}

// Set a particular field as Sorting field to
// order by
// @ param field string
func (d *DbAdapter) SetOrderByField(field string, order persistence.OrderBy) {
	d.orderByClause = append(d.orderByClause, fmt.Sprintf("%s %s", field, order))
}

func (d *DbAdapter) prepareQueryWithOrderBy() {
	if len(d.orderByClause) != 0 {
		d.queryString = fmt.Sprintf("%s ORDER BY %s", d.queryString, strings.Join(d.orderByClause, ", "))
	}
}

// Sets limit from and to
func (d *DbAdapter) SetLimit(limitFrom, limitTo int) *DbAdapter {
	d.limitFrom = limitFrom
	d.limitTo = limitTo
	return d
}

func (d *DbAdapter) prepareQueryWithLimit(needsLimitBegin bool) {
	if needsLimitBegin && d.limitTo != 0 {
		d.queryString = fmt.Sprintf("%s LIMIT %d, %d", d.queryString, d.limitFrom, d.limitTo)
	} else if !needsLimitBegin && d.limitTo != 0 {
		d.queryString = fmt.Sprintf("%s LIMIT %d", d.queryString, d.limitTo)
	}
	// unset limit
	d.limitFrom = 0
	d.limitTo = 0
}

func (d *DbAdapter) setLastExectuedQuery(executedQuery string) {
	d.lastExecutedQuery = executedQuery
}

func (d *DbAdapter) GetLastExecutedQuery() string {
	return d.lastExecutedQuery
}

func (d *DbAdapter) ExecCustomQuery(queryToExecute string) map[int]map[string]string {

	placeHolder := map[int]map[string]string{}

	if !d.isConnectedToServer {
		return placeHolder
	}

	d.queryString = queryToExecute

	d.setLastExectuedQuery(d.queryString)

	//	fmt.Printf("query: %s\n", d.queryString)

	rows, dbQueryError := d._db.Query(d.queryString)

	//	fmt.Println(reflect.TypeOf(rows))

	// To avoid SQL sneaks, look if there
	// were any mysql query errors
	if nil == rows || dbQueryError != nil {

		// Debug
		fmt.Println(dbQueryError)
		fmt.Println(d.GetLastExecutedQuery())

		d._handleRedirectAndPanic.TriggerPanic("0003")
		d._handleRedirectAndPanic.RedirectToPanicUrl()
		return placeHolder
	}

	//unset the query
	defer d.unsetGlobalQueryVars()

	// free memory
	defer rows.Close()

	return d.dbResultToReadable(rows)
	//fmt.Printf("Res %v", readableDbResult)

}
