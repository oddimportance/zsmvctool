package lib

import (
	"fmt"
)

// returns int64, string
func (d *DbAdapter) Insert(fields []string, values []string) (int64, string) {

	if !d.isConnectedToServer {
		fmt.Println("Error with SQL Connection")
		return 0, "Has Ping Error"
	}

	var insertId = d.execInsert(d.prepareInsertStatement(fields, values))
	return insertId, Int64ToStr(insertId)
}

func (d *DbAdapter) prepareInsertStatement(fields []string, values []string) string {
	d.queryString = d.prepareFieldsForQueryStatement(fields, "`", ",", true)
	d.valueString = d.prepareFieldsForQueryStatement(values, "'", ",", true)
	// fmt.Println(fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", d.dbTable, d.queryString, d.valueString))
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", d.dbTable, d.queryString, d.valueString)
}

func (d *DbAdapter) execInsert(insertString string) int64 {

	var insert, insertError = d._db.Prepare(insertString)

	if insertError != nil {
		fmt.Println(insertError)
		d._handleRedirectAndPanic.TriggerPanic("00006")
		return 0
	}

	var insertRes, insertResError = insert.Exec()
	if insertResError != nil {
		fmt.Println(insertResError)
		d._handleRedirectAndPanic.TriggerPanic("00007")
		return 0
	}

	//fmt.Println(insertRes.LastInsertId())

	var lastInssertedId, error = insertRes.LastInsertId()
	if error == nil {
		return lastInssertedId
	} else {
		return 0
	}

	return 0
}
