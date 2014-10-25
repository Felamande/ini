package ini

import "strconv"
import "reflect"

const (
	bLMBrac = '['
	bRMBrac = ']'
	bColon  = ':'
	bComma  = ','
	bEqual  = '='
	bNega   = '-'
	bPoint  = '.'
)

type invalidError struct {
	line int
	sym  byte
}

func (e *invalidError) Error() string {
	sline := strconv.Itoa(e.line)
	return "syntax error: line " + sline + "missing " + string(e.sym)
}

type scanner struct {
	data   *[]byte
	sym    byte
	from   int
	length int
}

func (s *scanner) scansym() (loc int) {

	var i int
	for i = s.from; (*s.data)[i] != s.sym; i++ {
	}
	s.from = i + 1
	return i
}

func (s *scanner) checkTail() {

	if (*s.data)[s.length-1] != '\n' {
		*s.data = append(*s.data, '\n')
		s.length = s.length + 1
	}

}

func checkValid(pdata *[]byte) (err error) {
	if (*pdata)[len(*pdata)-1] != '\n' {
		(*pdata) = append(*pdata, '\n')
	}
	//for
	return &invalidError{1, ']'}
}

func convToStringSlice(data *[]byte) (re []string, err error) {
	//if e := checkValid(data); e != nil {
	//	return []string{"[failure]"}, e
	//}
	//
	var s scanner = scanner{data: data, sym: '\n', from: 0, length: len(*data)}
	s.checkTail()

	for s.from != s.length {
		thisfrom := s.from
		thisloc := s.scansym()
		re = append(re, string((*data)[thisfrom:thisloc]))
	}
	return re, nil
}

func Unmarshal(data *[]byte, v interface{}) {
	lines, _ := convToStringSlice(data)

	mainStruct := reflect.ValueOf(v).Elem()
	var section reflect.Value

	var s scanner

	for _, lineString := range lines {
		if string(lineString[0]) == "[" {
			sectionName := lineString[1 : len(lineString)-1]
			section = mainStruct.FieldByName(sectionName)
		} else {
			sd := []byte(lineString)
			s.data = &sd
			s.from = 0
			s.sym = bEqual
			s.length = len(lineString)
			loc := s.scansym()
			section.FieldByName(lineString[0:loc]).SetString(lineString[loc+1 : s.length])
		}
	}
}
