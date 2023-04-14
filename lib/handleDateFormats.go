package lib

import (
	"fmt"
	//	"github.com/asaskevich/govalidator"
	"errors"
	"time"
)

const (
	// year, month, day
	// 2006 Jan 02
	layoutISO = "2006-01-02"
)

// returns a mysql formated date on
// given language
func StandardizeDateToMysqlFormat(language, date string) (string, error) {
	_date, err := makeParseableDate(language, date)
	if err != nil {
		//		fmt.Println(err)
		return "", err
	}
	//	fmt.Println(_date.Format(layoutISO))
	return _date.Format(layoutISO), nil

}

func getDelimiterFromLanguage(language string) string {

	switch language {
	case "de-de":
		// germans date 31.12.2019
		return "."
	case "en-uk", "fr-fr", "en-us":
		// wierd ammis 12/31/2019
		// europeans 12/31/2019
		return "/"
	case "es-es":
		// the spanish 12-31-2019
		return "-"
	default:
		// default the rest of the world 12-31-2019
		return "-"
	}

}

func splitDate(date, delimiter string) []string {
	return SplitString(date, delimiter)
}

func makeParseableDate(language, date string) (time.Time, error) {
	// split the date to break lang barrier
	var splitedDate = splitDate(date, getDelimiterFromLanguage(language))

	// if the date has no 3 parts, its a fake
	if date == "" || len(splitedDate) != 3 {
		return time.Now(), errors.New("Invalid date")
	}

	// the f**%$$%* ammis have a wierd date format
	if "en-us" == language {
		// put month first
		return time.Parse(layoutISO, fmt.Sprintf("%s-%s-%s", splitedDate[2], splitedDate[0], splitedDate[1]))

	}

	// for the rest of the world
	return time.Parse(layoutISO, fmt.Sprintf("%s-%s-%s", splitedDate[2], splitedDate[1], splitedDate[0]))

}
