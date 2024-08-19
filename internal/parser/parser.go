package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type fields map[string]string

type Struct struct {
	name   string
	fields fields
}

func NewStruct(name string) *Struct {
	return &Struct{
		name:   name,
		fields: make(fields),
	}
}

func (s *Struct) AddArguments(fieldName, fieldType string) {
	s.fields[fieldName] = fieldType
}

func (s *Struct) WriteGetters() []string {
	var sb strings.Builder
	res := make([]string, 0)

	for fieldName, fieldType := range s.fields {
		sb.WriteString(fmt.Sprintf("func(%s %s) Get%s() %s {\n", strings.ToLower(string(s.name[0])), s.name, strings.ToUpper(fieldName[:1])+fieldName[1:], fieldType))
		sb.WriteString(fmt.Sprintf("\treturn %s.%s\n", strings.ToLower(string(s.name[0])), fieldName))
		sb.WriteString("}\n")
		res = append(res, sb.String())
		sb.Reset()
	}
	return res
}

func (s *Struct) WriteSetters() []string {
	var sb strings.Builder
	res := make([]string, 0)

	for fieldName, fieldType := range s.fields {
		sb.WriteString(fmt.Sprintf("func(%s %s) Set%s(val %s) {\n", strings.ToLower(string(s.name[0])), s.name, strings.ToUpper(fieldName[:1])+fieldName[1:], fieldType))
		sb.WriteString(fmt.Sprintf("\t%s.%s = val\n", strings.ToLower(string(s.name[0])), fieldName))
		sb.WriteString("}\n")
		res = append(res, sb.String())
		sb.Reset()
	}
	return res

}

func ParseFile(fileName string) ([]string, error) {
	structStringArr, err := readStructsFromFile(fileName)
	if err != nil {
		return nil, err
	}
	structArr := make([]Struct, 0)
	for _, item := range structStringArr {
		structArr = append(structArr, parseStruct(item))
	}

	funcs := make([]string, 0)
	for _, item := range structArr {
		funcs = append(funcs, item.WriteGetters()...)
		funcs = append(funcs, item.WriteSetters()...)
	}
	return funcs, nil
}

func readStructsFromFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var structsArray []string
	flag := false
	var sb strings.Builder

	for scanner.Scan() {
		currentItem := scanner.Text()

		if currentItem == "type" {
			flag = true
		}
		if flag {
			sb.WriteString(currentItem)
			sb.WriteRune(' ')
		}
		if flag && currentItem == "}" {
			flag = false
			structsArray = append(structsArray, sb.String())
			sb.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return structsArray, nil
}

func parseStruct(structText string) Struct {
	structArr := strings.Fields(structText)
	structName := structArr[1]
	res := NewStruct(structName)

	isArguments := false
	for i := 1; i < len(structArr); i++ {
		if strings.HasSuffix(structArr[i], "{") {
			isArguments = true
			continue
		}
		if structArr[i] == "}" {
			break
		}
		if isArguments {
			res.AddArguments(structArr[i], structArr[i+1])
			i++
		}
	}

	return *res

}
