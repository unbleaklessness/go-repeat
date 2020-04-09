package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func findTemplate(name string, templates []template) (template, int) {

	for i, template := range templates {
		if template.Name == name {
			return template, i
		}
	}

	return template{}, -1
}

func readTemplates() ([]template, ierrori) {

	templatesFilePath, ie := getTemplatesFilePath()
	if ie != nil {
		return nil, ie
	}

	templatesBytes, e := ioutil.ReadFile(templatesFilePath)
	if e != nil {
		return nil, ierror{m: "Could not read templates file", e: e}
	}

	templates := make([]template, 0)

	e = json.Unmarshal(templatesBytes, &templates)
	if e != nil {
		return nil, ierror{m: "Could not parse templates file", e: e}
	}

	return templates, nil
}

func writeTemplates(templates []template) ierrori {

	newTemplatesBytes, e := json.Marshal(templates)
	if e != nil {
		return ierror{m: "Could not save templates", e: e}
	}

	templatesFilePath, ie := getTemplatesFilePath()
	if ie != nil {
		return ie
	}

	e = ioutil.WriteFile(templatesFilePath, newTemplatesBytes, os.ModePerm)
	if e != nil {
		return ierror{m: "Could not save templates", e: e}
	}

	return nil
}

func addTemplate(name string, p string) ierrori {

	p = filepath.Clean(p)

	bytes, e := ioutil.ReadFile(p)
	if e != nil {
		return ierror{m: "Could not read provided file", e: e}
	}

	templates, ie := readTemplates()
	if ie != nil {
		return ie
	}

	if _, i := findTemplate(name, templates); i != -1 {
		return ierror{m: "Template with given name is aready exists"}
	}

	newTemplate := template{
		Name:     name,
		FileName: filepath.Base(p),
		Bytes:    bytes,
	}

	templates = append(templates, newTemplate)

	ie = writeTemplates(templates)
	if ie != nil {
		return ie
	}

	return nil
}

func useTemplate(unitDirectoryPath string, templateName string, isQuestion bool, inline string) ierrori {

	templates, ie := readTemplates()
	if ie != nil {
		return ie
	}

	template, i := findTemplate(templateName, templates)
	if i == -1 {
		return ierror{m: "Could not find template with given name"}
	}

	templateFilePath := ""

	if isQuestion {
		templateFilePath = filepath.Join(unitDirectoryPath, questionDirectoryName, template.FileName)
	} else {
		templateFilePath = filepath.Join(unitDirectoryPath, answerDirectoryName, template.FileName)
	}

	data := []byte{}

	if len(inline) > 0 {
		data = []byte(inline)
	} else {
		data = template.Bytes
	}

	e := ioutil.WriteFile(templateFilePath, data, os.ModePerm)
	if e != nil {
		return ierror{m: "Could not write template file", e: e}
	}

	if len(inline) <= 0 {
		e = open(templateFilePath)
		if e != nil {
			return ierror{m: "Could not open template file", e: e}
		}
	}

	return nil
}

func deleteTemplate(name string) ierrori {

	templates, ie := readTemplates()
	if ie != nil {
		return ie
	}

	_, i := findTemplate(name, templates)
	if i == -1 {
		return ierror{m: "Could not find template with provided name"}
	}

	templates[i] = templates[len(templates)-1]
	templates = templates[:len(templates)-1]

	ie = writeTemplates(templates)
	if ie != nil {
		return ie
	}

	return nil
}

func listTemplates() ierrori {

	templates, ie := readTemplates()
	if ie != nil {
		return ie
	}

	for i, template := range templates {
		fmt.Printf("%d) %s / %s\n", i+1, template.Name, template.FileName)
	}

	return nil
}
