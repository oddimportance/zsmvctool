package lib

import (
	"fmt"
	"log"
	"strings"

	"zsmvctool/persistence"
)

func (d *DbAdapter) Where(field string, value string, wildcard persistence.Wildcard) *DbAdapter {

	d.whereClause = d.setWhere(d.whereClause, field, value, wildcard)

	return d

}

func (d *DbAdapter) WhereOr(field string, value string, wildcard persistence.Wildcard) *DbAdapter {

	d.whereOrClause = d.setWhere(d.whereOrClause, field, value, wildcard)
	return d

}

// Use Where In Clause for multiple
// wheres on same column. IT IS 3x
// faster than multiple ORs on one
// same field
func (d *DbAdapter) WhereIn(field string, values []string) *DbAdapter {
	valueString := ""
	for _, val := range values {
		// remove if any unwanted whitespaces
		val = TrimWhiteSpaces(val)
		// make sure there is not mysql sneak
		// and that it is only a number
		if ValidatorIsNumberOnly(val) {
			valueString = fmt.Sprintf("%s, %s", valueString, val)
		}
	}
	// remove the leading comma and whitespace
	valueString = StringSubString(valueString, 2, len(valueString))

	d.whereClause = append(d.whereClause, fmt.Sprintf("`%s` IN(%s)", field, valueString))
	return d

}

func (d *DbAdapter) setWhere(whereClause []string, field string, value string, wildcard persistence.Wildcard) []string {
	// log.Printf("Key: %s :: Value: %s", field, value)
	// if ContainsString("<", value) || ContainsString(">", value) || ContainsString("(", value) || ContainsString(")", value) || ContainsString("=", value) || ContainsString(" OR ", value) || ContainsString(" AND ", value) || ContainsString("SLEEP", value) || ContainsString("SELECT", value) || ContainsString("UNION", value) {
	// 	d.queryHasPotentialThreat = true
	// 	whereClause = nil
	// 	log.Println(d.queryString)
	// 	d.queryString = ""
	// 	log.Printf("threat detected in key: %s val %s. unsetting where clause\n", field, value)
	// } else {
	// 	whereClause = append(whereClause, d.setWildcard(whereClause, field, StringReplace(EscapeHtml(value), "'", "\\'", -1), wildcard))
	// }
	if d.checkValueForWhereAgainstSqlInjections(value, field) {
		whereClause = append(whereClause, d.setWildcard(whereClause, field, StringReplace(EscapeHtml(value), "'", "\\'", -1), wildcard))
	}
	return whereClause
}

func (d *DbAdapter) WhereMathComparision(field string, value string, mathComparision persistence.MathComparision) *DbAdapter {
	d.whereClause = d.setWhereMathComparision(d.whereClause, field, value, mathComparision)
	return d
}

func (d *DbAdapter) WhereOrMathComparision(field string, value string, mathComparision persistence.MathComparision) *DbAdapter {
	d.whereOrClause = d.setWhereMathComparision(d.whereOrClause, field, value, mathComparision)
	return d
}

func (d *DbAdapter) setWhereMathComparision(whereClause []string, field string, value string, mathComparision persistence.MathComparision) []string {

	if d.checkValueForWhereAgainstSqlInjections(value, field) {
		if mathComparision == d.MathComparision.IsNull() {
			whereClause = append(whereClause, fmt.Sprintf("`%s` %s", field, mathComparision))
		} else {
			whereClause = append(whereClause, fmt.Sprintf("`%s` %s '%s'", field, mathComparision, value))
		}
	}
	return whereClause
}

func (d *DbAdapter) setWildcard(whereClause []string, field string, value string, wildcard persistence.Wildcard) string {

	if wildcard != "" {
		return fmt.Sprintf("`%s` LIKE '%s'", field, strings.Replace(string(wildcard), "#", value, -1))
	} else {
		return fmt.Sprintf("`%s` = '%s'", field, value)
	}

}

func (d *DbAdapter) prepareWhereForSelect() {

	if len(d.whereClause) != 0 {
		d.queryString = fmt.Sprintf("%s WHERE (%s)", d.queryString, d.parseWhere(d.whereClause, " AND "))

		if len(d.whereOrClause) != 0 {
			d.queryString = fmt.Sprintf("%s AND (%s)", d.queryString, d.parseWhere(d.whereOrClause, " OR "))
		}
	}

	if len(d.whereClause) == 0 && len(d.whereOrClause) != 0 {
		d.queryString = fmt.Sprintf("%s WHERE (%s)", d.queryString, d.parseWhere(d.whereOrClause, " OR "))
	}

	// d.queryString = fmt.Sprintf("%s WHERE 'a'='a'", d.queryString)

	// fmt.Println(d.queryString)

}

func (d *DbAdapter) parseWhere(where []string, attribute string) string {
	return strings.Join(where, attribute)
}

func (d *DbAdapter) checkValueForWhereAgainstSqlInjections(valueForWhere, field string) bool {
	if ContainsString("<", valueForWhere) || ContainsString(">", valueForWhere) || ContainsString("(", valueForWhere) || ContainsString(")", valueForWhere) || ContainsString("=", valueForWhere) || ContainsString(" OR ", valueForWhere) || ContainsString(" AND ", valueForWhere) || ContainsString("SLEEP", valueForWhere) || ContainsString("SELECT", valueForWhere) || ContainsString("UNION", valueForWhere) {
		d.queryHasPotentialThreat = true
		d.whereClause = nil
		log.Println(d.queryString)
		d.queryString = ""
		log.Printf("threat detected in key: %s val %s. unsetting where clause\n", field, valueForWhere)
		return false
	}
	return true
}
