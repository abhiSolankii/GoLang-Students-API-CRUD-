package storage

import "github.com/abhiSolankii/students-api-go-lang/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudentById(id int64, student types.Student) (types.Student, error)
}
