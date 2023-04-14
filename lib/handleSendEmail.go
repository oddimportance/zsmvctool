package lib

import (
	"bytes"
	"fmt"
	"html"
	"html/template"

	"zsmvctool/persistence"

	"gopkg.in/gomail.v2"
)

// * Instantiate Mailer Class
// var _mailer = new(lib.HandleSMTPEmail)
//
// * Init necessary classes
// _mailer.Init(h._handleRedirectAndPanic)
//
// * Set Subject
// _mailer.SetSubject("Hey it works, Alhumdulillah!")
//
// * Prepare key value pair to be parsed in mail
// var keyValuesToBeParsed = []map[string]string{}
//
// * Attach a General file
// _mailer.Attach(PathOfFileToAttach)
//
// * A list of possible Keys with respected values
//
//	keyValuesToBeParsed = append(keyValuesToBeParsed, map[string]string{
//			"mailTo":           "oddimportance@gmail.com",
//			"First_name":        "Odd Impotance",
//			"ConfirmationLink": "http://www.google.de",
//
//			** A Constant that sets inidividual files as attachments **
//			"filesToAttach":    h._envConfigVars.GetConfVar("privateStorage") + "/test.csv"})
//
// * Set common css (CSS file must be in the same template folder)
// _mailer.SetCSS("commonCSS")
//
// * Now that every thing is finally set, trigger the mail
// _mailer.SendMail("register", rest)
type HandleSMTPEmail struct {
	_envConfigVars          GetEnvConfigVars
	smtpDetails             map[string]string
	_handleRedirectAndPanic *HandleRedirectAndPanic
	_mailer                 *gomail.Message
	defaultLang             string
	emailTemplatePath       string
	_readFile               ReadFile
	cssString               string
	subject                 string
}

func (h *HandleSMTPEmail) Init(defaultLang string, _handleRedirectAndPanic *HandleRedirectAndPanic) {
	h.defaultLang = defaultLang
	h._handleRedirectAndPanic = _handleRedirectAndPanic
	h.setEnvConfigVars()
	h.setSmtpDetails()
	h.SetEmailTemplatePath(h._envConfigVars.GetConfVar("emailTemplatePath"))
	h.initReadFile()
	h.setDailer()

}

func (h *HandleSMTPEmail) setDailer() {
	h._mailer = gomail.NewMessage()
}

func (h *HandleSMTPEmail) setSmtpDetails() {

	var smtpDetailsFromJson = SplitString(h._envConfigVars.GetConfVar("smtpDetails"), ",")
	h.smtpDetails = map[string]string{
		"smtpServer":       TrimWhiteSpaces(smtpDetailsFromJson[0]),
		"smtpPort":         TrimWhiteSpaces(smtpDetailsFromJson[1]),
		"smtpFrom":         TrimWhiteSpaces(smtpDetailsFromJson[2]),
		"smtpAuthUser":     TrimWhiteSpaces(smtpDetailsFromJson[3]),
		"smtpAuthPassword": TrimWhiteSpaces(smtpDetailsFromJson[4])}
}

func (h *HandleSMTPEmail) setEnvConfigVars() {
	h._envConfigVars = GetEnvConfigVars{}
	h._envConfigVars.Initiate(persistence.EnvConfigVarsFilePath)
}

func (h *HandleSMTPEmail) initReadFile() {
	h._readFile = ReadFile{}
}

func (h *HandleSMTPEmail) SetSubject(subject string) {
	h.subject = subject
}

// Public method to set path for email templates
// This method is called during the Init, which
// tries to set the default path from Environment
// Config Vars
func (h *HandleSMTPEmail) SetEmailTemplatePath(emailTemplatePath string) {
	h.emailTemplatePath = emailTemplatePath
}

// Reads the content of a css file and sets as
// string
// @ param cssFile string
func (h *HandleSMTPEmail) SetCSS(cssFile string) {
	h.cssString = string(h._readFile.Initiate(fmt.Sprintf("%s/%s.css", h.emailTemplatePath, cssFile)))
}

func (h *HandleSMTPEmail) getContentFromTemplate(emailTemplate string) string {

	pathToFile := fmt.Sprintf("%s/%s/%s.html", h.emailTemplatePath, h.defaultLang, emailTemplate)
	if h.defaultLang == "" || !h._readFile.FileExists(pathToFile) {
		pathToFile = fmt.Sprintf("%s/%s.html", h.emailTemplatePath, emailTemplate)
	}

	return FindReplace(
		string(h._readFile.Initiate(pathToFile)),
		"##CSS##",
		h.cssString)
}

func (h *HandleSMTPEmail) SendMail(emailTemplate string, mailData map[string]string) {
	var tmpl *template.Template = h.makeTemplate(emailTemplate)
	//	for _, data := range mailData {
	//	h.checkForFilesToAttach(data["filesToAttach"])
	h.checkForFilesToAttach(mailData["filesToAttach"])
	//		h.send(data["mailTo"], h.parseTemplate(tmpl, data))
	h.send(mailData["mailTo"], h.parseTemplate(tmpl, mailData))
	// unset all file
	// h._mailer.UnsetAllAttachedFiles()
	//	}
}

func (h *HandleSMTPEmail) SendMailWithoutTemplate(mailTo, mailBody string) {
	var tmpl *template.Template = h.makeTemplate("custom_mail_content")
	h.send(mailTo, h.parseTemplate(tmpl, map[string]string{"CustomMailContent": mailBody, "Subject": h.subject}))
}

func (h *HandleSMTPEmail) checkForFilesToAttach(filePathsString string) {
	if filePathsString == "" {
		return
	}

	var filePaths = SplitString(filePathsString, ",")
	for _, file := range filePaths {
		h.AttachFile(TrimWhiteSpaces(file))
	}
}

func (h *HandleSMTPEmail) makeTemplate(emailTemplate string) *template.Template {

	return template.Must(template.New(fmt.Sprintf("%s.html", emailTemplate)).Funcs(
		template.FuncMap{
			"TemplateFunctionStringToHtml": h.TemplateFunctionStringToHtml,
		}).Parse(h.getContentFromTemplate(emailTemplate)))

}

func (h *HandleSMTPEmail) parseTemplate(tmpl *template.Template, valuesToParse map[string]string) string {

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, valuesToParse); err != nil {
		fmt.Println(err)
	}

	return tpl.String()

}

func (h *HandleSMTPEmail) SetCC(cc string) {
	h._mailer.SetAddressHeader("Cc", cc, "")
}

func (h *HandleSMTPEmail) AttachFile(pathToFile string) {
	h._mailer.Attach(pathToFile)
}

func (h *HandleSMTPEmail) send(mailTo string, body string) {

	from := h.smtpDetails["smtpFrom"]
	user := h.smtpDetails["smtpAuthUser"]
	pass := h.smtpDetails["smtpAuthPassword"]

	//h._mailer := gomail.NewMessage()
	h._mailer.SetHeader("From", from)
	h._mailer.SetHeader("To", mailTo)
	// Set CC
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	h._mailer.SetHeader("Subject", h.subject)
	h._mailer.SetBody("text/html", body)
	// attach files like following
	//m.Attach(h._envConfigVars.GetConfVar("privateStorage") + "/bg-title-01.jpg")
	d := gomail.NewDialer(h.smtpDetails["smtpServer"], StringToInt(h.smtpDetails["smtpPort"]), user, pass)
	d.SSL = false

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(h._mailer); err != nil {
		//panic(err)
		fmt.Println(err)
	}
}

func (h *HandleSMTPEmail) TemplateFunctionStringToHtml(s string) template.HTML {
	//return template.HTML(s)
	return template.HTML(html.UnescapeString(s))
}
