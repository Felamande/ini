package ini

import "strconv"
import "reflect"

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

func checkValid(data *[]byte) (err error) {
	if (*data)[len(*data)-1] != '\n' {
		(*data) = append(*data, '\n')
	}
	return &invalidError{1, ']'}
}

func convToStringSlice(data *[]byte) (re []string) {
	//if e := checkValid(data); e != nil {
	//	return []string{"[failure]"}, e
	//}
	//
	var s scanner = scanner{data: data, sym: '\n', from: 0, length: len(*data)}
	s.checkTail()

	for s.from != s.length {
		thisfrom := s.from
		thisloc := s.scansym()
		if thisloc != thisfrom {
			re = append(re, string((*data)[thisfrom:thisloc]))
		}
	}
	return re
}

func Unmarshal1(data *[]byte, v interface{}) {
	lines := convToStringSlice(data)

	mainStruct := reflect.ValueOf(v).Elem()
	var section reflect.Value

	var s scanner

	for index, lineString := range lines {
		length := len(lineString)
		if string(lineString[0]) == "[" {
			sectionName := lineString[1 : length-1]
			preSection := section
			if section = mainStruct.FieldByName(sectionName); !section.CanSet() {
				if index == 0 {
					section = mainStruct.Field(0)
				} else {
					section = preSection
				}
			}
		} else {
			sd := []byte(lineString)
			s.data = &sd
			s.from = 0
			s.sym = '='
			s.length = length
			loc := s.scansym()
			if Key := section.FieldByName(lineString[0:loc]); Key.CanSet() {
				Key.SetString(lineString[loc+1 : length])
			}
		}
	}
}

func Unmarshal2(data *[]byte, v interface{}) {
	mainStruct := reflect.ValueOf(v).Elem()
	var section reflect.Value

	var s scanner = scanner{data: data, from: 0, length: len(*data)}
	s.checkTail()
	for s.from != s.length {
		if (*data)[s.from] == '[' {
			s.sym = ']'
			thisfrom := s.from
			thisloc := s.scansym()
			sectionName := (*data)[thisfrom+1 : thisloc]
			if section = mainStruct.FieldByName(string(sectionName)); !section.CanSet() {
				s.sym = '['
				s.scansym()
				if s.from != s.length {
					s.from = s.from - 1
				}
			}
			if s.length != s.from {
				for (*data)[s.from] == '\n' {
					s.from = s.from + 1
				}
			}
		} else {
			s.sym = '='
			thisfrom := s.from
			thisloc := s.scansym()
			if key := section.FieldByName(string((*data)[thisfrom:thisloc])); key.CanSet() {
				s.sym = '\n'
				thisfrom := s.from
				thisloc := s.scansym()
				key.SetString(string((*data)[thisfrom:thisloc]))
			} else {
				s.sym = '\n'
				s.scansym()
			}
			if s.length != s.from {
				for (*data)[s.from] == '\n' {
					s.from = s.from + 1
				}
			}
		}
	}
}
