package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	// "github.com/arthurkushman/go-hungarian"

	hungarianAlgorithm "github.com/oddg/hungarian-algorithm"
	"github.com/tealeg/xlsx"
)

const (
	SheetCourse       = "Course"
	SheetTeacher      = "Teacher"
	SheetRegistration = "Registration"

	BasicType    = 1
	ElectiveType = 0
	TempPriority = 999

	NotApplicable = "N/A"
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
	excelFileName := os.Args[1]
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println("Read excel file error: ", err.Error())
		return
	}

	var (
		courses      []*Course
		teachers     = make(map[string]*Teacher)
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
				teachers[t.ID] = &t
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
					if priority < 0 { // Teacher does not invole in this course so set by temp priority
						priority = TempPriority
					}
					teacherRegistration = append(teacherRegistration, &TeacherRegistration{
						CourseID:  courseID,
						TeacherID: header[i],
						Priority:  priority,
					})
				}

				// Fill the right most cell have no value
				if len(row.Cells)-1 < len(header) {
					for i := len(row.Cells); i <= len(header); i++ {
						teacherRegistration = append(teacherRegistration, &TeacherRegistration{
							CourseID:  courseID,
							TeacherID: header[i],
							Priority:  TempPriority,
						})
					}
				}

				registration[courseID] = teacherRegistration
			}
		}
	}

	// Algo
	sort.Slice(courses, func(i, j int) bool {
		return courses[i].NoClasses < courses[j].NoClasses
	})

	fmt.Println("Input information: -----------------------")
	printInput(courses, teachers, registration)

	start := time.Now().UTC().UnixNano()
	Run(courses, teachers, registration)
	end := time.Now().UTC().UnixNano()
	fmt.Println("Execution time ns: ", end-start)

}

func printInput(courses []*Course, teachers map[string]*Teacher, registration map[string][]*TeacherRegistration) {
	fmt.Println("course: ")
	var basicSumNoClass, ElectiveSumNoClass int
	for _, v := range courses {
		if v.Type == BasicType {
			basicSumNoClass += v.NoClasses
		} else {
			ElectiveSumNoClass += v.NoClasses
		}

		fmt.Printf("  %v %v %v %v\n", v.ID, v.Name, v.Type, v.NoClasses)
	}
	fmt.Printf("Basic: %v; Elective: %v; Total: %v\n", basicSumNoClass, ElectiveSumNoClass, basicSumNoClass+ElectiveSumNoClass)

	fmt.Println("teacher: ")
	for _, v := range teachers {
		fmt.Printf("  %v %v %v\n", v.ID, v.Name, v.MaxClass)
	}

	// -----------------
	fmt.Println("registration: ")
	header := "     "
	var regisRow []string
	idx := 0
	for k, v := range registration {
		tr := "  " + k

		for _, m := range v {
			if idx == 0 {
				header += " " + m.TeacherID
			}

			tr += " " + fmt.Sprintf("%3v", m.Priority)
		}
		idx++
		regisRow = append(regisRow, tr)
	}
	sort.Slice(regisRow, func(i, j int) bool {
		return regisRow[i] < regisRow[j]
	})

	fmt.Println(header)
	for _, v := range regisRow {
		fmt.Println(v)
	}
}

func Run(courses []*Course, teachers map[string]*Teacher, registration map[string][]*TeacherRegistration) {
	var result []TeacherAssignmentResult

	// Assign basic course first
	basicCount := countCourseType(courses, BasicType)
	for basicCount > 0 {
		teacherAssignmentResult := courseAssignment(courses, teachers, registration, BasicType)

		// fmt.Println("Basic course result: ")
		// printTeacherAssignmentResult(teacherAssignmentResult)

		result = append(result, teacherAssignmentResult...)

		basicCount = countCourseType(courses, BasicType)
	}

	// Assign elective course after
	electiveCount := countCourseType(courses, ElectiveType)
	for electiveCount > 0 {
		teacherAssignmentResult := courseAssignment(courses, teachers, registration, ElectiveType)

		// fmt.Println("Elective course result: ")
		// printTeacherAssignmentResult(teacherAssignmentResult)

		result = append(result, teacherAssignmentResult...)

		electiveCount = countCourseType(courses, ElectiveType)
	}

	printResult(result, courses)
}

func countCourseType(courses []*Course, courseType int) int {
	var count int
	for _, c := range courses {
		if c.Type == courseType && c.NoClasses > 0 {
			count++
		}
	}
	return count
}

func printResult(result []TeacherAssignmentResult, courses []*Course) {
	courseInfo := make(map[string]*Course)
	for _, v := range courses {
		courseInfo[v.ID] = v
	}

	var (
		basicAssignedCount, basicPriority                   int
		electiveAssignedCount, electivePriority             int
		notApplicableBasicCount, notApplicableElectiveCount int
		courseMap                                           = make(map[string]int)
	)
	for i := 0; i < len(result); i++ {
		row := result[i]

		courseMap[row.CourseID]++
		classID := fmt.Sprintf("%v-%v", row.CourseID, courseMap[row.CourseID])

		var priority string
		if row.Priority == TempPriority {
			priority = NotApplicable
		} else {
			priority = strconv.Itoa(row.Priority)
		}

		// Print result info
		fmt.Printf("%v %v %v %v %v\n", i+1, classID, row.CourseID, row.TeacherID, priority)

		// Count by course type
		if row.Priority == TempPriority {
			if courseInfo[row.CourseID].Type == BasicType {
				notApplicableBasicCount++
			} else {
				notApplicableElectiveCount++
			}
		} else {
			if courseInfo[row.CourseID].Type == BasicType {
				basicAssignedCount++
				basicPriority += row.Priority
			} else {
				electiveAssignedCount++
				electivePriority += row.Priority
			}
		}
	}

	// Calculate standard dev
	totalPriority := basicPriority + electivePriority
	totalAssigned := basicAssignedCount + electiveAssignedCount
	averagePriority := totalPriority / totalAssigned
	var std float64
	for _, v := range result {
		if v.Priority == TempPriority {
			continue
		}
		std += math.Pow((float64(v.Priority) - float64(averagePriority)), 2)
	}
	std = math.Sqrt(std / (float64(totalAssigned) - 1))

	fmt.Println("Result: --------------------------------------------")
	fmt.Printf("Total Basic assigned / Total Basic not applicable: %v/%v\n", basicAssignedCount, notApplicableBasicCount)
	fmt.Printf("Total Elective assigned / Total Elective not applicable: %v/%v\n", electiveAssignedCount, notApplicableElectiveCount)
	fmt.Printf("Basic priority: %v; Elective priority: %v\n", basicPriority, electivePriority)
	fmt.Printf("Total priority / StdDev: %v/%.5v\n", totalPriority, std)
}

func courseAssignment(
	courses []*Course,
	teachers map[string]*Teacher,
	registration map[string][]*TeacherRegistration,
	courseType int,
) []TeacherAssignmentResult {
	var rowCourses []*Course
	for _, c := range courses {
		if c.Type == courseType && c.NoClasses > 0 {
			rowCourses = append(rowCourses, c)
		}
	}

	var result []TeacherAssignmentResult
	if len(rowCourses) == 0 {
		log.Fatalln("courseAssignment error no row course")
	} else if len(rowCourses) == 1 {
		result = oneCourseAssignment(rowCourses, teachers, registration)
	} else {
		result = hungarianAssignment(rowCourses, teachers, registration)
	}

	return result
}

func oneCourseAssignment(
	rowCourses []*Course,
	teachers map[string]*Teacher,
	registration map[string][]*TeacherRegistration,
) []TeacherAssignmentResult {
	course := rowCourses[0]
	teacherRegistrations, ok := registration[course.ID]
	if !ok {
		log.Fatalln("oneCourseAssignment error: this course have no registration", course.ID)
	}

	var (
		min    = TempPriority
		tcr    *TeacherRegistration
		result TeacherAssignmentResult
	)

	for _, v := range teacherRegistrations {
		teacher, ok := teachers[v.TeacherID]
		if !ok {
			continue
		}

		if teacher.MaxClass > 0 && min > v.Priority { // Find min priority
			tcr = v
			min = v.Priority
		}
	}

	if tcr != nil { // This course do have a teacher
		result = TeacherAssignmentResult{
			CourseID:  tcr.CourseID,
			TeacherID: tcr.TeacherID,
			Priority:  tcr.Priority,
		}

		teachers[tcr.TeacherID].MaxClass-- // Reduce the teacher max class
		// tcr.Priority = TempPriority        // The teacher have assigned to the course so remove from the registration list
	} else {
		result = TeacherAssignmentResult{
			CourseID:  course.ID,
			TeacherID: NotApplicable,
			Priority:  TempPriority,
		}
	}

	// This course have checked so reduce number of class by 1
	course.NoClasses--

	return []TeacherAssignmentResult{result}
}

func hungarianAssignment(
	rowCourses []*Course,
	teachers map[string]*Teacher,
	registration map[string][]*TeacherRegistration,
) []TeacherAssignmentResult {
	colTeachers := make(map[string]bool)
	for _, v := range rowCourses {
		teacherRegistrations, ok := registration[v.ID]
		if !ok {
			log.Fatalln("Registration does not have this classes: ", v.ID)
		}
		for _, m := range teacherRegistrations {
			teacher, ok := teachers[m.TeacherID]
			if !ok {
				continue
			}
			if teacher.MaxClass > 0 {
				colTeachers[m.TeacherID] = true
			}
		}
	}

	var teacherID []string
	for k := range colTeachers {
		teacherID = append(teacherID, k)
	}
	sort.Slice(teacherID, func(i, j int) bool {
		return teacherID[i] < teacherID[j]
	})

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

	registrationMatrix := make([][]*TeacherRegistration, size)
	for i := 0; i < lrow; i++ {
		registrationMatrix[i] = make([]*TeacherRegistration, size)

		for j := 0; j < lcol; j++ {
			cID := rowCourses[i].ID
			tcRegis, ok := registration[cID]
			if !ok {
				log.Fatalln("Build matrix error, registration does not have this class: ", cID)
			}

			// Get the teacher who have invole in the course from teacherIDs
			var tr *TeacherRegistration
			for _, v := range tcRegis {
				if v.TeacherID == teacherID[j] {
					tr = v
					break
				}
			}
			if tr == nil {
				for _, v := range tcRegis {
					fmt.Printf("%v ", v.TeacherID)
				}
				fmt.Printf("\n%v\n", teacherID)
				log.Fatalln("Get teacher from registration list not found: ", tcRegis)
			}

			registrationMatrix[i][j] = tr
		}
	}

	if need > 0 && lrow < lcol {
		registrationMatrix = fillRowTemPriority(registrationMatrix, need, size)
	} else if need > 0 && lrow > lcol {
		registrationMatrix = fillColTemPriority(registrationMatrix, need, size)
	}

	fmt.Println("----------------------------")
	// fmt.Println("registrationMatrix: ", registrationMatrix)
	matrix := extractPriority(registrationMatrix, size)
	// fmt.Println("Priority matrix:")
	// printPriorityMatrix(matrix)

	hungarianResult, err := hungarianAlgorithm.Solve(matrix)
	if err != nil {
		log.Fatalln("Hungarian assignment error: ", err.Error())
	}
	// fmt.Println("result: ", hungarianResult)

	return buildResult(rowCourses, teachers, registrationMatrix, size, hungarianResult)
}

func fillRowTemPriority(m [][]*TeacherRegistration, need, size int) [][]*TeacherRegistration {
	for i := size - need; i < size; i++ {
		m[i] = make([]*TeacherRegistration, size)
		for j := 0; j < size; j++ {
			m[i][j] = &TeacherRegistration{Priority: TempPriority}
		}
	}

	return m
}

func fillColTemPriority(m [][]*TeacherRegistration, need, size int) [][]*TeacherRegistration {
	for i := 0; i < size; i++ {
		for j := size - need; j < size; j++ {
			m[i][j] = &TeacherRegistration{Priority: TempPriority}
		}
	}

	return m
}

func extractPriority(registrationMatrix [][]*TeacherRegistration, size int) [][]int {
	m := make([][]int, size)
	for i := 0; i < size; i++ {
		m[i] = make([]int, size)
		for j := 0; j < size; j++ {
			m[i][j] = registrationMatrix[i][j].Priority
		}
	}

	return m
}

func buildResult(
	rowCourses []*Course,
	teachers map[string]*Teacher,
	registrationMatrix [][]*TeacherRegistration,
	size int,
	hunarianResult []int,
) []TeacherAssignmentResult {

	var result []TeacherAssignmentResult
	for i := 0; i < size; i++ {
		j := hunarianResult[i]
		tcr := registrationMatrix[i][j]
		if tcr.TeacherID == "" {
			continue
		}

		teacher, ok := teachers[tcr.TeacherID]
		if !ok {
			continue
		}

		var teacherAssignmentResult TeacherAssignmentResult
		if teacher.MaxClass > 0 && tcr.Priority != TempPriority {
			teacherAssignmentResult = TeacherAssignmentResult{
				CourseID:  tcr.CourseID,
				TeacherID: tcr.TeacherID,
				Priority:  tcr.Priority,
			}

			teacher.MaxClass-- // Reduce the teacher max class
			// tcr.Priority = TempPriority // The teacher have assigned to the course so remove from the registration list
		} else {
			teacherAssignmentResult = TeacherAssignmentResult{
				CourseID:  rowCourses[i].ID,
				TeacherID: NotApplicable,
				Priority:  TempPriority,
			}
		}

		rowCourses[i].NoClasses-- // This course have checked so reduce number of class by 1

		result = append(result, teacherAssignmentResult)
	}

	return result
}

func printPriorityMatrix(m [][]int) {
	for i := 0; i < len(m); i++ {
		row := "  "
		for j := 0; j < len(m); j++ {
			row += " " + strconv.Itoa(m[i][j])
		}
		fmt.Println(row)
	}
}

func printTeacherAssignmentResult(result []TeacherAssignmentResult) {
	for _, v := range result {
		fmt.Printf("   %v %v %v\n", v.CourseID, v.TeacherID, v.Priority)
	}
}
