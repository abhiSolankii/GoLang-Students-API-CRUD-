package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/abhiSolankii/students-api-go-lang/internal/storage"
	"github.com/abhiSolankii/students-api-go-lang/internal/types"
	"github.com/abhiSolankii/students-api-go-lang/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			slog.Error("Error decoding JSON:", slog.String("Error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			slog.Error("Error decoding JSON:", slog.String("Error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//Request validation
		if err := validator.New().Struct(student); err != nil {
			slog.Error("Error validating request:", slog.String("Error", err.Error()))
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		slog.Info("Creating a student")

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		if err != nil {
			slog.Error("Error creating student:", slog.String("Error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("Student created succefully", slog.String("userId", fmt.Sprint(lastId)))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if idStr == "" {
			slog.Error("id is empty")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is required")))
			return
		}
		//Type of id is string so we need to convert it to int64
		id, err := strconv.ParseInt(idStr, 10, 64) // 10 is base of the number and 64 is size of the int
		if err != nil {
			slog.Error("failed to parse id", slog.String("err", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		slog.Info("Getting a student", slog.Int64("id", id))
		student, err := storage.GetStudentById(id)
		if err != nil {
			slog.Error("error getting student", slog.String("id", idStr))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			slog.Error("error getting students", slog.String("Error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}
