package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/arthurkushman/go-hungarian"
	"github.com/tealeg/xlsx"
)

const (
	SheetCourse       = "Course"
	SheetTeacher      = "Teacher"
	SheetRegistration = "Registration"

	BasicType    = 1
	TempPriority = 999
)

type Course struct {
	ID        string
	Name      string
	Type      int
	NoClasses int
}

type Teacher struct {
	ID       string
	Name     string
	MaxClass int
}

type TeacherRegistration struct {
	CourseID  string
	TeacherID string
	Priority  int
}

type TeacherAssignmentResult struct {
	No        int
	ClassID   string
	CourseID  string
	TeacherID string
	Priority  int
}

func main() {
	// Read excel file
	excelFileName := "dataSample4BTL.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println("Read excel file error: ", err.Error())
		return
	}

	var (
		courses      []*Course
		teachers     []*Teacher
		registration = make(map[string][]*TeacherRegistration)
	)

	for _, sheet := range xlFile.Sheets {
		if sheet.Name == SheetCourse { // Read Course
			for idx, row := range sheet.Rows {
				if idx == 0 {
					continue
				}
				if len(row.Cells) != 4 {
					continue
				}
				cType, _ := row.Cells[2].Int()
				cNoClass, _ := row.Cells[3].Int()
				if cNoClass == -1 {
					continue
				}
				c := Course{
					ID:        strings.TrimSpace(row.Cells[0].String()),
					Name:      strings.TrimSpace(row.Cells[1].String()),
					Type:      cType,
					NoClasses: cNoClass,
				}
				courses = append(courses, &c)
			}
		}

		if sheet.Name == SheetTeacher { // Read Teahcher
			for idx, row := range sheet.Rows {
				if idx == 0 {
					continue
				}
				if len(row.Cells) != 3 {
					continue
				}
				maxClass, _ := row.Cells[2].Int()
				if maxClass <= 0 {
					continue
				}
				t := Teacher{
					ID:       strings.TrimSpace(row.Cells[0].String()),
					Name:     strings.TrimSpace(row.Cells[1].String()),
					MaxClass: maxClass,
				}
				teachers = append(teachers, &t)
			}
		}

		if sheet.Name == SheetRegistration { // Read Registration
			header := make(map[int]string)
			for idx, row := range sheet.Rows {
				if idx == 0 {
					for i := 1; i < len(row.Cells); i++ {
						header[i] = strings.TrimSpace(row.Cells[i].String())
					}
					continue
				}

				courseID := strings.TrimSpace(row.Cells[0].String())
				var teacherRegistration []*TeacherRegistration
				for i := 1; i < len(row.Cells); i++ {
					priority, _ := row.Cells[i].Int()
					teacherRegistration = append(teacherRegistration, &TeacherRegistration{
						CourseID:  courseID,
						TeacherID: header[i],
						Priority:  priority,
					})
				}
				registration[courseID] = teacherRegistration
			}
		}
	}

	// Algo
	sort.Slice(courses, func(i, j int) bool {
		return courses[i].NoClasses < courses[j].NoClasses
	})

	printInput(courses, teachers, registration)
	Run(courses, teachers, registration)
}

func printInput(courses []*Course, teachers []*Teacher, registration map[string][]*TeacherRegistration) {
	fmt.Println("course: ")
	for _, v := range courses {
		fmt.Println("", *v)
	}
	fmt.Println("teacher: ")
	for _, v := range teachers {
		fmt.Println("", *v)
	}
	fmt.Println("registration: ")
	for k, v := range registration {
		tr := k + ": "
		for _, m := range v {
			tr += m.CourseID + " " + m.TeacherID + " " + strconv.Itoa(m.Priority) + " "
		}
		fmt.Println(tr)
	}
}

func Run(courses []*Course, teachers []*Teacher, registration map[string][]*TeacherRegistration) {
	assignBasic(courses, teachers, registration)
}

func assignBasic(courses []*Course, teachers []*Teacher, registration map[string][]*TeacherRegistration) {
	var result []TeacherAssignmentResult

	basicCount := countBasicCourses(courses)
	for basicCount > 0 {
		teacherAssignmentResult := hungarianAssignment(courses, teachers, registration, BasicType)
		result = append(result, teacherAssignmentResult...)

		basicCount = countBasicCourses(courses)
	}

	fmt.Println("assignBasic: ", result)
}

func countBasicCourses(courses []*Course) int {
	var count int
	for _, c := range courses {
		if c.Type == BasicType && c.NoClasses > 0 {
			count++
		}
	}
	return count
}

func hungarianAssignment(
	courses []*Course,
	teachers []*Teacher,
	registration map[string][]*TeacherRegistration,
	courseType int,
) []TeacherAssignmentResult {
	var rowCourses []Course
	var colTeachers []Teacher
	for _, c := range courses {
		if c.Type == courseType && c.NoClasses > 0 {
			rowCourses = append(rowCourses, *c)
			c.NoClasses--
		}
	}
	for _, t := range teachers {
		if t.MaxClass > 0 {
			colTeachers = append(colTeachers, *t)
			t.MaxClass--
		}
	}

	// Build matrix
	lrow := len(rowCourses)
	lcol := len(colTeachers)
	var need, size int
	if lrow > lcol {
		need = lrow - lcol
		size = lrow
	} else if lrow < lcol {
		need = lcol - lrow
		size = lcol
	} else {
		size = lrow
	}

	registrationMatrix := make([][]TeacherRegistration, size)
	for i := 0; i < lrow; i++ {
		registrationMatrix[i] = make([]TeacherRegistration, size)

		for j := 0; j < lcol; j++ {
			cID := rowCourses[i].ID
			teacherRegistration := registration[cID]

			teacherPriority := teacherRegistration[j].Priority
			if teacherPriority < 0 { // Teacher does not invole in this course so set by temp priority
				teacherPriority = TempPriority
			}

			registrationMatrix[i][j] = TeacherRegistration{
				CourseID:  cID,
				TeacherID: teacherRegistration[j].TeacherID,
				Priority:  teacherPriority,
			}

		}
	}

	if need > 0 && lrow < lcol {
		registrationMatrix = fillRowTemPriority(registrationMatrix, need, size)
	} else if need > 0 && lrow > lcol {
		registrationMatrix = fillColTemPriority(registrationMatrix, need, size)
	}

	fmt.Println("matrix: ", registrationMatrix)
	matrix := extractPriority(registrationMatrix, size)
	hungarianResult := hungarian.SolveMin(matrix)
	fmt.Println("result: ", hungarianResult)

	return buildResult(registrationMatrix, size, hungarianResult)
}

func fillRowTemPriority(m [][]TeacherRegistration, need, size int) [][]TeacherRegistration {
	for i := size - need; i < size; i++ {
		m[i] = make([]TeacherRegistration, size)
		for j := 0; j < size; j++ {
			m[i][j].Priority = TempPriority
		}
	}

	return m
}

func fillColTemPriority(m [][]TeacherRegistration, need, size int) [][]TeacherRegistration {
	for i := 0; i < size; i++ {
		for j := size - need; j < size; j++ {
			m[i][j].Priority = TempPriority
		}
	}

	return m
}

func extractPriority(registrationMatrix [][]TeacherRegistration, size int) [][]float64 {
	m := make([][]float64, size)
	for i := 0; i < size; i++ {
		m[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			m[i][j] = float64(registrationMatrix[i][j].Priority)
		}
	}

	return m
}

func buildResult(registrationMatrix [][]TeacherRegistration, size int, hunarianResult map[int]map[int]float64) []TeacherAssignmentResult {
	var result []TeacherAssignmentResult
	for i := 0; i < size; i++ {
		val, ok := hunarianResult[i]
		if !ok {
			log.Fatalln("hunarianResult does not have i index: ", i)
		}

		for k, v := range val {
			if v >= float64(size) {
				break
			}

			tcRegis := registrationMatrix[k][int(v)]

			teacherAssignmentResult := TeacherAssignmentResult{
				ClassID:   tcRegis.CourseID,
				CourseID:  tcRegis.CourseID,
				TeacherID: tcRegis.TeacherID,
				Priority:  tcRegis.Priority,
			}
			result = append(result, teacherAssignmentResult)
		}
	}

	return result
}
