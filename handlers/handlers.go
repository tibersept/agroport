package handlers

import (
	"agroport/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx/v3"
)

type Handler struct {
	db *sql.DB
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "error", Message: message})
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Worker handlers
func (h *Handler) CreateWorker(w http.ResponseWriter, r *http.Request) {
	var worker models.Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if worker.Name == "" || worker.Role == "" {
		h.respondWithError(w, http.StatusBadRequest, "Name and role are required")
		return
	}

	if err := models.CreateWorker(&worker); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create worker")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, SuccessResponse{
		Message: "Worker created successfully",
		Data:    worker,
	})
}

func (h *Handler) GetWorkers(w http.ResponseWriter, r *http.Request) {
	workers, err := models.GetWorkers()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch workers")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Workers retrieved successfully",
		Data:    workers,
	})
}

func (h *Handler) GetWorker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid worker ID")
		return
	}

	worker, err := models.GetWorkerByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.respondWithError(w, http.StatusNotFound, "Worker not found")
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch worker")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Worker retrieved successfully",
		Data:    worker,
	})
}

func (h *Handler) UpdateWorker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid worker ID")
		return
	}

	var worker models.Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	worker.ID = id
	if err := models.UpdateWorker(&worker); err != nil {
		if err == sql.ErrNoRows {
			h.respondWithError(w, http.StatusNotFound, "Worker not found")
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to update worker")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Worker updated successfully",
		Data:    worker,
	})
}

func (h *Handler) DeleteWorker(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid worker ID")
		return
	}

	if err := models.DeleteWorker(id); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete worker")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Worker deleted successfully",
	})
}

// Field handlers
func (h *Handler) CreateField(w http.ResponseWriter, r *http.Request) {
	var field models.Field
	if err := json.NewDecoder(r.Body).Decode(&field); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if field.Name == "" {
		h.respondWithError(w, http.StatusBadRequest, "Field name is required")
		return
	}

	if err := models.CreateField(&field); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create field")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, SuccessResponse{
		Message: "Field created successfully",
		Data:    field,
	})
}

func (h *Handler) GetFields(w http.ResponseWriter, r *http.Request) {
	fields, err := models.GetFields()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch fields")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Fields retrieved successfully",
		Data:    fields,
	})
}

func (h *Handler) GetField(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid field ID")
		return
	}

	field, err := models.GetFieldByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.respondWithError(w, http.StatusNotFound, "Field not found")
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch field")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Field retrieved successfully",
		Data:    field,
	})
}

func (h *Handler) UpdateField(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid field ID")
		return
	}

	var field models.Field
	if err := json.NewDecoder(r.Body).Decode(&field); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	field.ID = id
	if err := models.UpdateField(&field); err != nil {
		if err == sql.ErrNoRows {
			h.respondWithError(w, http.StatusNotFound, "Field not found")
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to update field")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Field updated successfully",
		Data:    field,
	})
}

func (h *Handler) DeleteField(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid field ID")
		return
	}

	if err := models.DeleteField(id); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete field")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Field deleted successfully",
	})
}

// Schedule handlers
func (h *Handler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var schedule models.Schedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if schedule.WorkerID == 0 {
		h.respondWithError(w, http.StatusBadRequest, "Worker ID and Field ID are required")
		return
	}

	if err := models.CreateSchedule(&schedule); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create schedule")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, SuccessResponse{
		Message: "Schedule created successfully",
		Data:    schedule,
	})
}

func (h *Handler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := models.GetSchedules()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch schedules")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Schedules retrieved successfully",
		Data:    schedules,
	})
}

func (h *Handler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	schedule, err := models.GetScheduleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.respondWithError(w, http.StatusNotFound, "Schedule not found")
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch schedule")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Schedule retrieved successfully",
		Data:    schedule,
	})
}

func (h *Handler) GetWorkerSchedules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workerID, err := strconv.Atoi(vars["workerId"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid worker ID")
		return
	}

	schedules, err := models.GetWorkerSchedules(workerID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch worker schedules")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Worker schedules retrieved successfully",
		Data:    schedules,
	})
}

func (h *Handler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	var schedule models.Schedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	schedule.ID = id
	if err := models.UpdateSchedule(&schedule); err != nil {
		if err == sql.ErrNoRows {
			h.respondWithError(w, http.StatusNotFound, "Schedule not found")
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to update schedule")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Schedule updated successfully",
		Data:    schedule,
	})
}

func (h *Handler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	if err := models.DeleteSchedule(id); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete schedule")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Schedule deleted successfully",
	})
}

// Operation handlers
func (h *Handler) CreateOperation(w http.ResponseWriter, r *http.Request) {
	var operation models.Operation
	if err := json.NewDecoder(r.Body).Decode(&operation); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if operation.WorkerID == 0 || operation.FieldID == 0 || operation.Type == "" {
		h.respondWithError(w, http.StatusBadRequest, "Worker ID, Field ID, and Type are required")
		return
	}

	if operation.Status == "" {
		operation.Status = "planned"
	}

	if err := models.CreateOperation(&operation); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create operation")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, SuccessResponse{
		Message: "Operation created successfully",
		Data:    operation,
	})
}

func (h *Handler) GetOperations(w http.ResponseWriter, r *http.Request) {
	operations, err := models.GetOperations()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch operations")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Operations retrieved successfully",
		Data:    operations,
	})
}

func (h *Handler) GetOperation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid operation ID")
		return
	}

	operation, err := models.GetOperationByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.respondWithError(w, http.StatusNotFound, "Operation not found")
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch operation")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Operation retrieved successfully",
		Data:    operation,
	})
}

func (h *Handler) UpdateOperation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid operation ID")
		return
	}

	var operation models.Operation
	if err := json.NewDecoder(r.Body).Decode(&operation); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	operation.ID = id
	if err := models.UpdateOperation(&operation); err != nil {
		if err == sql.ErrNoRows {
			h.respondWithError(w, http.StatusNotFound, "Operation not found")
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to update operation")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Operation updated successfully",
		Data:    operation,
	})
}

func (h *Handler) CompleteOperation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid operation ID")
		return
	}

	if err := models.CompleteOperation(id); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to complete operation")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Operation completed successfully",
	})
}

func (h *Handler) DeleteOperation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid operation ID")
		return
	}

	if err := models.DeleteOperation(id); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete operation")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Operation deleted successfully",
	})
}

// Report handlers
func (h *Handler) GetDailyReport(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	var date time.Time
	var err error

	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			h.respondWithError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
			return
		}
	} else {
		date = time.Now()
	}

	report, err := models.GetDailyReport(date)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to generate daily report")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Daily report generated successfully",
		Data:    report,
	})
}

func (h *Handler) GetMonthlyReport(w http.ResponseWriter, r *http.Request) {
	// Placeholder for monthly report logic
	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Monthly report endpoint - implementation needed",
	})
}

func (h *Handler) GetYearlyReport(w http.ResponseWriter, r *http.Request) {
	// Placeholder for yearly report logic
	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Yearly report endpoint - implementation needed",
	})
}

func (h *Handler) ExportDailyReport(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	var date time.Time
	var err error

	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			h.respondWithError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
			return
		}
	} else {
		date = time.Now()
	}

	report, err := models.GetDailyReport(date)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to generate daily report")
		return
	}

	// Create Excel file
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Daily Report")
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create Excel sheet")
		return
	}

	// Add report summary
	row := sheet.AddRow()
	row.AddCell().Value = "Daily Report"
	row.AddCell().Value = date.Format("2006-01-02")

	sheet.AddRow() // Empty row

	// Summary section
	summaryRow := sheet.AddRow()
	summaryRow.AddCell().Value = "Summary"

	row = sheet.AddRow()
	row.AddCell().Value = "Total Workers"
	row.AddCell().Value = fmt.Sprintf("%d", report.TotalWorkers)

	row = sheet.AddRow()
	row.AddCell().Value = "Total Operations"
	row.AddCell().Value = fmt.Sprintf("%d", report.TotalOperations)

	row = sheet.AddRow()
	row.AddCell().Value = "Completed Operations"
	row.AddCell().Value = fmt.Sprintf("%d", report.CompletedOps)

	row = sheet.AddRow()
	row.AddCell().Value = "In Progress Operations"
	row.AddCell().Value = fmt.Sprintf("%d", report.InProgressOps)

	sheet.AddRow() // Empty row

	// Operations by type
	if len(report.OperationsByType) > 0 {
		typeRow := sheet.AddRow()
		typeRow.AddCell().Value = "Operations by Type"

		for opType, count := range report.OperationsByType {
			row = sheet.AddRow()
			row.AddCell().Value = opType
			row.AddCell().Value = fmt.Sprintf("%d", count)
		}
		sheet.AddRow() // Empty row
	}

	// Worker statistics
	if len(report.WorkerStats) > 0 {
		workerHeaderRow := sheet.AddRow()
		workerHeaderRow.AddCell().Value = "Worker Statistics"

		headerRow := sheet.AddRow()
		headerRow.AddCell().Value = "Worker Name"
		headerRow.AddCell().Value = "Operations"
		headerRow.AddCell().Value = "Hours Worked"
		headerRow.AddCell().Value = "Fields Worked"

		for _, ws := range report.WorkerStats {
			row = sheet.AddRow()
			row.AddCell().Value = ws.WorkerName
			row.AddCell().Value = fmt.Sprintf("%d", ws.Operations)
			row.AddCell().Value = fmt.Sprintf("%.2f", ws.HoursWorked)
			row.AddCell().Value = fmt.Sprintf("%d", ws.FieldsWorked)
		}
		sheet.AddRow() // Empty row
	}

	// Field statistics
	if len(report.FieldStats) > 0 {
		fieldHeaderRow := sheet.AddRow()
		fieldHeaderRow.AddCell().Value = "Field Statistics"

		headerRow := sheet.AddRow()
		headerRow.AddCell().Value = "Field Name"
		headerRow.AddCell().Value = "Operations"
		headerRow.AddCell().Value = "Hours Worked"
		headerRow.AddCell().Value = "Workers Count"

		for _, fs := range report.FieldStats {
			row = sheet.AddRow()
			row.AddCell().Value = fs.FieldName
			row.AddCell().Value = fmt.Sprintf("%d", fs.Operations)
			row.AddCell().Value = fmt.Sprintf("%.2f", fs.HoursWorked)
			row.AddCell().Value = fmt.Sprintf("%d", fs.WorkersCount)
		}
	}

	// Set response headers for file download
	filename := fmt.Sprintf("daily_report_%s.xlsx", date.Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// Write Excel file to response
	if err := file.Write(w); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to write Excel file")
		return
	}
}

func (h *Handler) ExportMonthlyReport(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Monthly report export endpoint - implementation needed",
	})
}

func (h *Handler) ExportYearlyReport(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Yearly report export endpoint - implementation needed",
	})
}
