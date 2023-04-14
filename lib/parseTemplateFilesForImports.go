package lib

import (
	"fmt"
	//"io/ioutil"
	"regexp"
	"strings"
)

type ParseTemplateFilesForImports struct {
	_handleFilePath             *HandleFilePath
	templatePath                string
	filesToBeParsed             []string
	localCSSIncludes            []string
	externalCSSIncludes         []string
	localJSIncludesForHeader    []string
	externalJSIncludesForHeader []string
	localJSIncludesForFooter    []string
	externalJSIncludesForFooter []string
	filesAlreadyParsed          map[string]int
	_readFile                   ReadFile
}

func (p *ParseTemplateFilesForImports) InitParsing(handleFilePath *HandleFilePath, baseLayout string, controller string, action string) []string {

	//p._readFile = new(*ReadFile)

	p.setHandleFilePath(handleFilePath)

	p.setTemplatePath()

	p.initFilesAlreadyParsed()

	p.addBaseLyoutAndActionTemplate(baseLayout, controller, action)

	p.initiateParsingImports()

	//fmt.Println(p.filesToBeParsed)

	return p.filesToBeParsed
}

func (p *ParseTemplateFilesForImports) GetLocalCSSFiles() []string {
	return p.localCSSIncludes
}

func (p *ParseTemplateFilesForImports) GetExternalCSSFiles() []string {
	return p.externalCSSIncludes
}

func (p *ParseTemplateFilesForImports) GetLocalJSFilesForHeader() []string {
	return p.localJSIncludesForHeader
}

func (p *ParseTemplateFilesForImports) GetExternalJSFilesForHeader() []string {
	return p.externalJSIncludesForHeader
}

func (p *ParseTemplateFilesForImports) GetLocalJSFilesForFooter() []string {
	return p.localJSIncludesForFooter
}

func (p *ParseTemplateFilesForImports) GetExternalJSFilesForFooter() []string {
	return p.externalJSIncludesForFooter
}

func (p *ParseTemplateFilesForImports) addBaseLyoutAndActionTemplate(baseLayout string, controller string, action string) {

	p.setFileNameToParseList(fmt.Sprintf("layouts/%s", baseLayout))
	p.setFileNameToParseList(fmt.Sprintf("views/%s/%s", controller, action))

}

func (p *ParseTemplateFilesForImports) setFileNameToParseList(fileName string) {

	var file = p.makeFilePath(fileName)

	if p.isFileNotInParseList(file) {
		p.filesAlreadyParsed[file] = 1
		p.filesToBeParsed = append(p.filesToBeParsed, file)
	}

}

func (p *ParseTemplateFilesForImports) isFileNotInParseList(file string) bool {

	if p.filesAlreadyParsed[file] != 1 {
		return true
	}

	return false

}

func (p *ParseTemplateFilesForImports) makeFilePath(fileName string) string {
	return fmt.Sprintf("%s/%s.html", p.templatePath, fileName)
}

func (p *ParseTemplateFilesForImports) initiateParsingImports() {

	//var file = fmt.Sprintf("%s/layouts/%s.html", p.templatePath, baseLayout)

	var tempFileHolder []string

	var j int = len(p.filesToBeParsed)

	for i := 0; i < j; i++ {

		//fmt.Println(p.filesToBeParsed[i])
		tempFileHolder = p.readTemplateContent(p.filesToBeParsed[i])
		for _, k := range tempFileHolder {
			//fileHolderToReturn = append(fileHolderToReturn, k)
			p.setFileNameToParseList(strings.TrimSpace(k))
			j++
		}

	}
}

func (p *ParseTemplateFilesForImports) readTemplateContent(file string) []string {

	var fileContent = string(p._readFile.Initiate(file))

	//	var filesToLoad []string = p.parseFileContentForImports(fileContent)
	var filesToLoad []string = p.parseIncludes(fileContent, "import")

	// set the local css to be included
	p.localCSSIncludes = p.appendFileNamesToIncludes(p.localCSSIncludes, p.parseIncludes(fileContent, "setLocalCSS"))
	// set the external css to be included
	p.externalCSSIncludes = p.appendFileNamesToIncludes(p.externalCSSIncludes, p.parseIncludes(fileContent, "setExternalCSS"))

	// set the local js to be included in footer
	p.localJSIncludesForFooter = p.appendFileNamesToIncludes(p.localJSIncludesForFooter, p.parseIncludes(fileContent, "setLocalJSForFooter"))
	// set the external js to be included in footer
	p.externalJSIncludesForFooter = p.appendFileNamesToIncludes(p.externalJSIncludesForFooter, p.parseIncludes(fileContent, "setExternalJSForFooter"))

	// set the local js to be included in header
	p.localJSIncludesForHeader = p.appendFileNamesToIncludes(p.localJSIncludesForHeader, p.parseIncludes(fileContent, "setLocalJSForHeader"))
	// set the external js to be included in header
	p.externalJSIncludesForHeader = p.appendFileNamesToIncludes(p.externalJSIncludesForHeader, p.parseIncludes(fileContent, "setExternalJSForHeader"))

	//fmt.Println(string(b))
	return filesToLoad
}

func (p *ParseTemplateFilesForImports) parseIncludes(fileContent, parseMarker string) []string {
	re := regexp.MustCompile(fmt.Sprintf(`.*::%s{.*}end%s`, parseMarker, parseMarker))
	matches := re.FindStringSubmatch(fileContent)
	//	fmt.Println(matches)
	//	fmt.Printf("Array Size %d\n", len(matches))

	if len(matches) == 1 {
		//		fmt.Println(matches[0])
		return p.extractFileNamesFromString(matches[0], parseMarker)
	}
	return matches
}

func (p *ParseTemplateFilesForImports) appendFileNamesToIncludes(fileNamesArray, filesToAppend []string) []string {
	for _, fileName := range filesToAppend {
		if strings.TrimSpace(fileName) != "" {
			fileNamesArray = append(fileNamesArray, strings.TrimSpace(fileName))
		}
	}
	return fileNamesArray
}

func (p *ParseTemplateFilesForImports) extractFileNamesFromString(importString, parseMarker string) []string {
	//fmt.Println(p.stripQuotes(importString))
	var cleanString = p.cleanString(importString, parseMarker)

	//fmt.Println(cleanString)

	if cleanString == "" {
		return nil
	}

	return strings.Split(cleanString, ",")
}

func (p *ParseTemplateFilesForImports) cleanString(importString, marker string) string {
	return p.stripQuotes(p.stripImports(importString, marker))
}

func (p *ParseTemplateFilesForImports) stripQuotes(importString string) string {
	return strings.Replace(importString, "\"", "", -1)
	//return importString
}

func (p *ParseTemplateFilesForImports) stripImports(importString, marker string) string {
	return strings.TrimLeft(strings.TrimRight(importString, fmt.Sprintf("}end%s", marker)), fmt.Sprintf("::%s{", marker))
}

func (p *ParseTemplateFilesForImports) setTemplatePath() {

	p.templatePath = p._handleFilePath.GetTemplatePath()
	//	fmt.Println(p.templatePath)
}

func (p *ParseTemplateFilesForImports) setHandleFilePath(handleFilePath *HandleFilePath) {
	p._handleFilePath = handleFilePath

}

func (p *ParseTemplateFilesForImports) initFilesAlreadyParsed() {
	p.filesAlreadyParsed = make(map[string]int)
}
