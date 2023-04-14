package lib

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html"
	"image"
	"image/jpeg"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"

	// "zsmvctool-api/persistence"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/kevinburke/twilio-go"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// constants to encrypt/decrypt
// ids for and from urls
const encryption_adder = 97
const encryption_multiplier = 13

var messagePrinter = message.NewPrinter(language.German)

// ---------- //

// matchInArray returns true if the given string value is in the array.
func InArray(arr []string, value string) bool {
	for _, v := range arr {
		if strings.TrimSpace(v) == value {
			return true
		}
	}
	return false
}

func ValidatorIsAphabetOnly(stringToValidate string) bool {
	return govalidator.IsAlpha(stringToValidate)
}

func ValidatorIsNumberOnly(stringToValidate string) bool {
	return govalidator.IsInt(stringToValidate)
}

// Generates a random number between range
// @param max int
// @param min int
// @return int
func GenerateRandomNumberInRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// Creates a random alphabetic string of given length
func GenerateRandomAlphaString(length int) string {
	var pattern = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	return createRandomString(length, pattern)
}

// Creates a random alphanumeric string of given length
func GenerateRandomAlphaNumericString(length int) string {
	var pattern = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	return createRandomString(length, pattern)
}

func createRandomString(length int, pattern []rune) string {

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, length)
	for i := range b {
		b[i] = pattern[rand.Intn(len(pattern))]
	}
	return string(b)
}

func HashStringToSHA512(stringToHash string, saltToHash string) string {
	var hashSHA512 = sha512.New()
	//io.WriteString(hashSHA512, fmt.Sprintf("%s%s", passwordToHash, saltToHash))
	hashSHA512.Write([]byte(fmt.Sprintf("%s%s", stringToHash, saltToHash)))
	//return string(hashSHA512.Sum(nil))
	return hex.EncodeToString(hashSHA512.Sum(nil))
}

func HashToMD5(stringToHash string) string {
	var hashMD5 = md5.New()
	//io.WriteString(hashMD5, string(uidToHash))
	hashMD5.Write([]byte(stringToHash))
	//return string(hashMD5.Sum(nil))
	return hex.EncodeToString(hashMD5.Sum(nil))
}

func Implode(arr []string, seperator string) string {
	return strings.Join(arr, seperator)
}

func IntToStr(intToConvert int) string {
	return strconv.Itoa(intToConvert)
}

func Int64ToStr(int64ToConvert int64) string {
	return strconv.FormatInt(int64ToConvert, 10)
}

func Float64ToStr(float64ToConvert float64) string {
	return strconv.FormatFloat(float64ToConvert, 'f', -1, 64)
}

func StringToInt(stringToConvert string) int {
	var intToReturn, _ = strconv.Atoi(stringToConvert)
	return intToReturn
}

func StringToInt64(stringToConvert string) int64 {
	var intToReturn, _ = strconv.ParseInt(stringToConvert, 10, 64)
	return intToReturn
}

func StringToFloat(stringToConvert string, bitSize int) float64 {
	var floatToReturn, _ = strconv.ParseFloat(NormalizeGermanFormat(stringToConvert), bitSize)
	return floatToReturn
}

func StringToUpper(s string) string {
	return strings.ToUpper(s)
}

func StringToLower(s string) string {
	return strings.ToLower(s)
}

func StringToBool(s string) bool {
	switch s {
	case "true":
		return true
	case "1":
		return true
	default:
		return false
	}
}

func ContainsString(match, s string) bool {
	return strings.Contains(s, match)
}

// Get preferred outbound ip of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func GetUserRealIP(r *http.Request) string {

	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		//		fmt.Println("IP Proxy:", ipProxy)
		// at times the ip address contained a port with a semicolne
		// hence pre check the length
		if len(ipProxy) > 15 {
			return StringSubString(ipProxy, 0, 14)
		}
		return ipProxy
	}

	/*
	 ****** Output of the following code has always been the host ip
	 ****** and not the client ip. Which is why, I've, for now decided
	 ****** to disable it and wait for a few hundered db entries and
	 ****** validate them further
	 */
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	//	fmt.Println(getOutboundIP())
	if ip != getOutboundIP() {
		//	fmt.Println("IP:", ip)
		return ip
	}
	return "0.0.0.0"

}
func Round(floatToRound float64) float64 {
	return math.Round(floatToRound)
}

func RoundToEvenNumber(floatToRound float64) int {
	return int(math.RoundToEven(floatToRound))
}

func RoundToGreater(floatToRound float64) int {
	var intToReturn = Round(floatToRound)
	if intToReturn < floatToRound {
		return int(intToReturn + 1)
	}
	return int(intToReturn)
}

func SplitString(stringToSplit string, seperator string) []string {
	return strings.Split(stringToSplit, seperator)
}

func TrimWhiteSpaces(s string) string {
	return strings.TrimSpace(s)
}

func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

func FindReplace(stringToSearch string, find string, replace string) string {
	return strings.Replace(stringToSearch, find, replace, -1)
}

func makeDateTimeFormString(format, dateTime, layout string) string {

	if format == "" || dateTime == "" {
		return dateTime
	}
	// layout eg. 2006-01-02 15:04:05
	dt, err := time.Parse(layout, dateTime)
	if err != nil {
		fmt.Println(err)
	}
	return string(dt.Format(format))

}

func LocalizeDate(format string, date string) string {
	// to avoid an error when date is datetime format
	// substring to date format by deleting time
	return makeDateTimeFormString(format, StringSubString(date, 0, 10), "2006-01-02")
}

func LocalizeDateTime(format string, dateTime string) string {
	return makeDateTimeFormString(format, dateTime, "2006-01-02 15:04:05")
}

// replace given stirng with new one
// repeat should be -1 for infinite
func StringReplace(s, replace, replaceWith string, repeat int) string {
	return strings.Replace(s, replace, replaceWith, repeat)
}

func VerifyLanguage(languageToVerify, languagesAvailable string) bool {
	var availableLanguages = strings.Split(languagesAvailable, ",")
	//	fmt.Println(languagesAvailable, languageToVerify, InArray(availableLanguages, languageToVerify))
	return InArray(availableLanguages, languageToVerify)

}

func StripTableFieldPrefixFromDataForView(tableFieldPrefix, replacer string, data map[string]string) map[string]string {
	// init slice to return
	var dataToReturn = map[string]string{}

	// walk through data
	for key, value := range data {
		// remove table field prefix
		dataToReturn[StringReplace(key, tableFieldPrefix, replacer, -1)] = value
	}

	return dataToReturn

}

// Remove field prefix from keys
// to access from view without
// table prefix
// Use GetSearchKeys from model
func StripTableFieldPrefixFromKeysForView(tableFieldPrefix string, keys []string, replaceWith string) map[string]string {
	// init array to return
	var dataToReturn = map[string]string{}

	// walk through data
	for _, key := range keys {
		// remove table field prefix
		dataToReturn[StringReplace(key, tableFieldPrefix, replaceWith, -1)] = key
	}

	return dataToReturn
}

func StripTablePrifixFromFields(fields []string, tableFieldPrefix, replaceWith string, capitalize bool) []string {

	// init array to return
	var dataToReturn = []string{}

	key := ""

	// walk through data
	for _, field := range fields {
		// remove table field prefix
		key = StringReplace(field, tableFieldPrefix, replaceWith, -1)
		if capitalize {
			key = strings.Title(key)
		}
		dataToReturn = append(dataToReturn, key)
	}

	return dataToReturn

}

func Removeprefix(MapT map[int]map[string]string, prefix string, replacement string) map[int]map[string]string {
	for _, innerMap := range MapT {
		for key, value := range innerMap {
			if strings.HasPrefix(key, prefix) {
				newKey := strings.Replace(key, prefix, replacement, 1)
				innerMap[newKey] = value
				delete(innerMap, key)
			}
		}
	}
	return MapT
}

func AppendToArray(existingArray, arrayToAppend []string) []string {
	for _, value := range arrayToAppend {
		existingArray = append(existingArray, strings.TrimSpace(value))
	}
	return existingArray
}

func AppendDataToSlice(existingSlice, sliceToAppend map[string]string) map[string]string {
	for key, value := range sliceToAppend {
		existingSlice[key] = value
	}
	return existingSlice
}

func AppendDataToTwoDimensionalSlice(existingSlice, sliceToAppend map[int]map[string]string) map[int]map[string]string {
	if len(sliceToAppend) != 0 {
		var iterateTill = len(existingSlice) + len(sliceToAppend)
		iterator := 0
		for i := len(existingSlice); i < iterateTill; i++ {
			existingSlice[i] = sliceToAppend[iterator]
			iterator++
		}
	}
	return existingSlice
}

// null value from db is %!s(<nil>),
// therefore this checks if it is null
// or empty
func IsNullOrEmpty(stringToVerify string) bool {

	// %!s(<nil>) == NULL
	if "%!s(<nil>)" == stringToVerify || "" == stringToVerify {
		return true
	}
	return false

}

// substrings a given string
// equivivalent of php substr()
func StringSubString(s string, subFrom, subTo int) string {
	// Take substring of first word with runes.
	// ... This handles any kind of rune in the string.
	runes := []rune(s)
	// ... Convert back into a string from rune slice.
	return string(runes[subFrom:subTo])
}

// strips html tags and attribs in textual
// representation
func EscapeHtml(htmlString string) string {
	return html.EscapeString(htmlString)
}

// parse textual representation of html
// to html
func UnescapeHtml(s string) string {
	return html.UnescapeString(s)
}

// Converts an interface to int
// to do so it must firstly convert interface val to float64
// and then rounds it to an int to avoid panic error if the
// value is a decimal type
func InterfaceToInt(interfaceToConvert map[string]interface{}) map[string]int {
	var mapToReturn = map[string]int{}
	for key, val := range interfaceToConvert {
		//		fmt.Println(val)
		mapToReturn[key] = RoundToEvenNumber(val.(float64))
	}
	return mapToReturn
}

// to avoid unwanted white spaces in email
// trim the input value
func CleanEmail(email string) string {
	return StringToLower(TrimWhiteSpaces(email))
}

// **** +++++ *****
// Please note it is not a mechanism to protect any data
// its main intention it to avoid hackers but also user
// to guess the next obvious numer. Therefore, do not
// use for and sensitive data
// ****************
// A mechanism to encrypt and decrypt ids passed via url
// Idea is to make it hard for hackers to guess the next
// obvious id ex.: if user_id=10 the successor would be 11
func EncryptIdForUrl(idToEncrypt string) string {
	//	fmt.Println("Data to encrypt:", idToEncrypt)
	return fmt.Sprintf("%d%d%d", GenerateRandomNumberInRange(10, 99), ((StringToInt(idToEncrypt) + encryption_adder) * encryption_multiplier), GenerateRandomNumberInRange(101, 999))
}

func DecryptIdFromUrl(idToDecrypt string) (int, string) {
	// First things first
	// 1- check if sting len is at least 9
	// because the min most id returns an
	// equation of (1 + 97)*13 = 1274 + rand2digit(99) + rand3digit(999) = 991274999
	// 2- make sure it is an in and not malfunctioned string
	if len(idToDecrypt) < 9 || StringToInt(idToDecrypt) == 0 {
		return 0, "0"
	}

	//	fmt.Println(idToDecrypt)
	// remove the leading two digits
	id := StringSubString(idToDecrypt, 2, len(idToDecrypt))
	//	fmt.Println(id)
	// strip the trailing 3 digits
	id = StringSubString(id, 0, len(id)-3)
	//	fmt.Println(id)
	// finally convert string to int and
	//	reverse engineer encryption
	var idToReturn = ((StringToInt(id) / encryption_multiplier) - encryption_adder)
	return idToReturn, IntToStr(idToReturn)
}

// obfuscate a string
// IBAN and SWIFT for example
func ObfuscateString(stringToObfuscate string, obfuscationLength int) string {
	var lengthOfStringToObfuscate = len(stringToObfuscate)
	return fmt.Sprintf("%sxxxxxxxx%s", StringSubString(stringToObfuscate, 0, obfuscationLength), StringSubString(stringToObfuscate, (lengthOfStringToObfuscate-obfuscationLength), lengthOfStringToObfuscate))
}

func MakeBuiltupareaGroupedElements(groupedBuiltupareas map[int]map[string]string, categoryName, builtupareaName string) (string, string) {

	var groupedCategoryList string = ""
	var groupedBuiltuparea string = ""
	var formatPattern = ""
	for i := 0; i < len(groupedBuiltupareas); i++ {
		if i == 0 {
			formatPattern = "%s%s"
		} else {
			formatPattern = "%s, %s"
		}
		groupedCategoryList = fmt.Sprintf(formatPattern,
			groupedCategoryList, fmt.Sprintf("%s (%s)",
				groupedBuiltupareas[i][categoryName],
				groupedBuiltupareas[i][builtupareaName],
			),
		)
		groupedBuiltuparea = fmt.Sprintf(formatPattern, groupedBuiltuparea, groupedBuiltupareas[i][builtupareaName])
	}
	return groupedCategoryList, groupedBuiltuparea
}

func IsCommaSeperatedFloat(stringToValidate string) bool {
	stringToValidate = NormalizeGermanFormat(stringToValidate)
	return stringToValidate != "" && (govalidator.IsInt(stringToValidate) || govalidator.IsFloat(FindReplace(stringToValidate, ",", ".")))
}

func NormalizeGermanFormat(stringToConvert string) string {
	// s := strings.Replace(stringToConvert, ",", ".", -1)
	// return strings.Replace(s, ".", "", 1)

	var s = strings.Replace(stringToConvert, ".", "", -1)
	return strings.Replace(s, ",", ".", 1)
}

// func NumberToGermanFormat(number string) string {
// 	if ContainsString(".", number) {
// 		return messagePrinter.Sprintf("%f", StringToFloat(number, 64))
// 	}
// 	return messagePrinter.Sprintf("%d", StringToInt(number))
// }

func NumberToGermanFormat(number interface{}) string {

	const float64Type float64 = 0.0
	var numberTypeof = reflect.TypeOf(number)
	var float64TypeOf = reflect.TypeOf(float64Type)

	if numberTypeof.Kind() == float64TypeOf.Kind() {
		return FindReplace(messagePrinter.Sprintf("%.2f", number), ",00", "")
	}
	return messagePrinter.Sprintf("%d", number)
}

func StructToArray(structToParse interface{}) map[string]string {

	v := reflect.ValueOf(structToParse)
	typeOfS := v.Type()

	var arrayToReturn = map[string]string{}
	var key string
	for i := 0; i < v.NumField(); i++ {
		key = StringToLower(typeOfS.Field(i).Name)
		// fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).String())
		arrayToReturn[key] = v.Field(i).String()
	}
	return arrayToReturn
}

func ActionpdateAvatar(base64String string) *os.File {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64String))
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err = jpeg.Encode(f, m, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
		log.Fatal(err)
	}

	return f
}

func GnerateOTP(length int) (string, error) {
	rand.Seed(time.Now().UnixNano())
	const letters = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b), nil
}

func Sendotp(
	accountSid string,
	authToken string,
	from string,
	to string,
	otp string) {

	message := fmt.Sprintf("Your OTP is %s", otp)

	// إنشاء عميل Twilio
	client := twilio.NewClient(accountSid, authToken, nil)

	// Send a message
	msg, err := client.Messages.SendMessage(from, to, message, nil)

	fmt.Println(msg)
	fmt.Println(err)
	// التحقق من الخطأ
	if err != nil {
		fmt.Println("Error sending message: ", err.Error())
	} else {
		fmt.Println("Message sent successfully")
	}
}

func RemoveprefixToArray(MapT []map[string]string, prefix string, replacement string) []map[string]string {
	for _, innerMap := range MapT {
		for key, value := range innerMap {
			if strings.HasPrefix(key, prefix) {
				newKey := strings.Replace(key, prefix, replacement, 1)
				innerMap[newKey] = value
				delete(innerMap, key)
			}
		}
	}

	return MapT
}

func RemoveprefixToElemnt(MapT map[string]string, prefix string, replacement string) map[string]string {

	for key, value := range MapT {
		if strings.HasPrefix(key, prefix) {
			newKey := strings.Replace(key, prefix, replacement, 1)
			MapT[newKey] = value
			delete(MapT, key)
		}
	}

	return MapT
}

// Creates a random alphanumeric string of given length
func GenerateRandomNumericString(length int) string {
	var pattern = []rune("0123456789")
	return createRandomString(length, pattern)
}
