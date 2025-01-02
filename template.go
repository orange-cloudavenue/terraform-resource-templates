package main

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/orange-cloudavenue/terraform-provider-resource-templates/pkg/file"

	_ "embed"
)

type templateDef struct {
	CategoryName          string
	ResourceName          string
	Name                  string
	PackageName           string
	LowerCamelName        string
	SnakeName             string
	CamelName             string
	Filename              string
	TestDir               string
	SchemaDir             string
	FullSnakeResourceName string
	FullCamelResourceName string
}

//go:embed templates/datasource.go.tmpl
var templateDatasource string

//go:embed templates/resource.go.tmpl
var templateResource string

//go:embed templates/acc_test_resource.go.tmpl
var templateAccTestResource string

//go:embed templates/acc_test_datasource.go.tmpl
var templateAccTestDataSource string

//go:embed templates/schema.go.tmpl
var templateSchema string

//go:embed templates/unit_test_schema.go.tmpl
var templateUnitTestSchema string

//go:embed templates/base.go.tmpl
var templateBase string

//go:embed templates/types.go.tmpl
var templateTypes string

//go:embed templates/common_templates.go.tmpl
var templateCommonTemplates string

func genTemplateConf(categoryName, resourceName, packageName, testDir, fileName, schemaDir string) templateDef {
	t := templateDef{
		CategoryName:          categoryName,
		ResourceName:          resourceName,
		PackageName:           packageName,
		LowerCamelName:        strcase.ToLowerCamel(resourceName),
		CamelName:             strcase.ToCamel(resourceName),
		SnakeName:             strcase.ToSnake(resourceName),
		Filename:              fileName,
		TestDir:               testDir,
		SchemaDir:             schemaDir,
		FullSnakeResourceName: categoryName + "_" + strcase.ToSnake(resourceName),
		FullCamelResourceName: categoryName + "_" + strcase.ToCamel(resourceName),
	}

	if resourceName == "" {
		t.LowerCamelName = strcase.ToLowerCamel(categoryName)
		t.CamelName = strcase.ToCamel(categoryName)
		t.SnakeName = strcase.ToSnake(categoryName)
		t.FullSnakeResourceName = strcase.ToSnake(categoryName)
		t.FullCamelResourceName = strcase.ToCamel(categoryName)
	}

	return t
}

func (t templateDef) createTemplateFiles(tfTypes string) error {
	templateS := templateDatasource
	templateAccTest := templateAccTestDataSource
	if tfTypes == "resource" {
		templateS = templateResource
		templateAccTest = templateAccTestResource
	}

	// * xx.go
	if err := createTemplate(t, t.Filename, templateCommonTemplates+templateS); err != nil {
		return err
	}

	// * testDir/categoryName_xx_test.go
	// for acc test
	if err := createTemplateIfNotExists(t, t.TestDir+"/"+t.CategoryName+"_"+fileNameWithoutExtAndPath(t.Filename)+"_test.go", templateCommonTemplates+templateAccTest); err != nil {
		return err
	}

	// * xx_schema.go
	// if file not already exists create schema file
	if err := createTemplateIfNotExists(t, t.SchemaDir+"/"+fileNameWithoutResourceOrDataSource(t.Filename)+"_schema.go", templateCommonTemplates+templateSchema); err != nil {
		return err
	}

	// * xx_schema_test.go
	// if file not already exists create schema test file
	if err := createTemplateIfNotExists(t, t.SchemaDir+"/"+fileNameWithoutResourceOrDataSource(t.Filename)+"_schema_test.go", templateCommonTemplates+templateUnitTestSchema); err != nil {
		return err
	}

	// * xx_types.go
	// if file not already exists create types file
	if err := createTemplateIfNotExists(t, t.SchemaDir+"/"+fileNameWithoutResourceOrDataSource(t.Filename)+"_types.go", templateCommonTemplates+templateTypes); err != nil {
		return err
	}

	// * base.go
	// if base file not already exists create it
	if err := createTemplateIfNotExists(t, t.SchemaDir+"/"+"base.go", templateCommonTemplates+templateBase); err != nil {
		return err
	}

	return nil
}

func fileNameWithoutExtAndPath(fileName string) string {
	f := filepath.Base(fileName)
	return strings.TrimSuffix(f, filepath.Ext(f))
}

// fileNameWithoutResourceOrDataSource returns the filename without the resource or datasource prefix.
func fileNameWithoutResourceOrDataSource(fileName string) string {
	f := fileNameWithoutExtAndPath(fileName)
	f = strings.TrimSuffix(f, "_resource")   // remove _resource
	f = strings.TrimSuffix(f, "_datasource") // remove _datasource
	return f
}

// createTemplateIfNotExists creates the template file if it does not exist.
func createTemplateIfNotExists(t templateDef, fileName, templateString string) error {
	if !file.IsFileExists(fileName) {
		return createTemplate(t, fileName, templateString)
	}

	return nil
}

// createTemplate creates the template file.
func createTemplate(t templateDef, fileName, templateString string) error {
	var tplTypes bytes.Buffer
	tmplTypes, err := template.New("template").Parse(templateString)
	if err != nil {
		return err
	}

	if err := tmplTypes.Execute(&tplTypes, t); err != nil {
		return err
	}

	// 0o600 syntax https://stackoverflow.com/questions/5624359/write-file-with-specific-permissions-in-python/5624691#5624691
	return os.WriteFile(fileName, tplTypes.Bytes(), 0o600)
}
