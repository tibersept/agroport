package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

var db *sql.DB

type Worker struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"` // "tractor_driver", "harvester_driver", etc.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Field struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Coordinates json.RawMessage `json:"coordinates"` // GeoJSON polygon
	Area        float64         `json:"area"`        // in decares
	CropType    string          `json:"crop_type"`
	Period      string          `json:"period"` // the period of operation
	Region      string          `json:"region"` // region name
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Schedule struct {
	ID        int       `json:"id"`
	WorkerID  int       `json:"worker_id"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Worker    *Worker   `json:"worker,omitempty"`
}

type Operation struct {
	ID          int        `json:"id"`
	ScheduleID  *int       `json:"schedule_id"`
	WorkerID    int        `json:"worker_id"`
	FieldID     int        `json:"field_id"`
	Type        string     `json:"type"` // "plowing", "seeding", "harvesting", etc.
	Description string     `json:"description"`
	Status      string     `json:"status"` // "planned", "in_progress", "completed", "cancelled"
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	CompletedAt *time.Time `json:"completed_at"`
	Notes       string     `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Schedule    *Schedule  `json:"schedule,omitempty"`
	Worker      *Worker    `json:"worker,omitempty"`
	Field       *Field     `json:"field,omitempty"`
}

type DailyReport struct {
	Date             time.Time          `json:"date"`
	TotalWorkers     int                `json:"total_workers"`
	TotalOperations  int                `json:"total_operations"`
	CompletedOps     int                `json:"completed_operations"`
	InProgressOps    int                `json:"in_progress_operations"`
	OperationsByType map[string]int     `json:"operations_by_type"`
	WorkerStats      []WorkerDailyStats `json:"worker_stats"`
	FieldStats       []FieldDailyStats  `json:"field_stats"`
}

type WorkerDailyStats struct {
	WorkerID     int     `json:"worker_id"`
	WorkerName   string  `json:"worker_name"`
	Operations   int     `json:"operations"`
	HoursWorked  float64 `json:"hours_worked"`
	FieldsWorked int     `json:"fields_worked"`
}

type FieldDailyStats struct {
	FieldID      int     `json:"field_id"`
	FieldName    string  `json:"field_name"`
	Operations   int     `json:"operations"`
	HoursWorked  float64 `json:"hours_worked"`
	WorkersCount int     `json:"workers_count"`
}

func InitDB(database *sql.DB) {
	db = database
}

func RunMigrations() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS workers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE,
			phone VARCHAR(50),
			role VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS fields (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			coordinates JSONB NOT NULL,
			area DECIMAL(10,2) DEFAULT 0,
			crop_type VARCHAR(100) NOT NULL,
			period VARCHAR(100) NOT NULL,
			region VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS schedules (
			id SERIAL PRIMARY KEY,
			worker_id INTEGER REFERENCES workers(id) ON DELETE CASCADE,
			date DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS operations (
			id SERIAL PRIMARY KEY,
			schedule_id INTEGER REFERENCES schedules(id) ON DELETE CASCADE,
			worker_id INTEGER REFERENCES workers(id) ON DELETE CASCADE,
			field_id INTEGER REFERENCES fields(id) ON DELETE CASCADE,
			type VARCHAR(100) NOT NULL,
			description TEXT,
			status VARCHAR(50) DEFAULT 'planned',
			start_time TIMESTAMP,
			end_time TIMESTAMP,
			completed_at TIMESTAMP,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_schedules_worker_date ON schedules(worker_id, date)`,
		`CREATE INDEX IF NOT EXISTS idx_operations_schedule ON operations(schedule_id) WHERE schedule_id IS NOT NULL`,
		`CREATE INDEX IF NOT EXISTS idx_operations_worker ON operations(worker_id)`,
		`CREATE INDEX IF NOT EXISTS idx_operations_field ON operations(field_id)`,
		`CREATE INDEX IF NOT EXISTS idx_operations_worker_date ON operations(worker_id, DATE(start_time))`,
		`CREATE INDEX IF NOT EXISTS idx_operations_field_date ON operations(field_id, DATE(start_time))`,
		`CREATE INDEX IF NOT EXISTS idx_operations_status ON operations(status)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}

// Worker methods
func CreateWorker(worker *Worker) error {
	query := `INSERT INTO workers (name, email, phone, role)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, updated_at`
	return db.QueryRow(query, worker.Name, worker.Email, worker.Phone, worker.Role).
		Scan(&worker.ID, &worker.CreatedAt, &worker.UpdatedAt)
}

func GetWorkers() ([]Worker, error) {
	query := `SELECT id, name, email, phone, role, created_at, updated_at FROM workers ORDER BY name`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workers []Worker
	for rows.Next() {
		var w Worker
		err := rows.Scan(&w.ID, &w.Name, &w.Email, &w.Phone, &w.Role, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}
		workers = append(workers, w)
	}
	return workers, nil
}

func GetWorkerByID(id int) (*Worker, error) {
	var w Worker
	query := `SELECT id, name, email, phone, role, created_at, updated_at FROM workers WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&w.ID, &w.Name, &w.Email, &w.Phone, &w.Role, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func UpdateWorker(worker *Worker) error {
	query := `UPDATE workers SET name = $1, email = $2, phone = $3, role = $4, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $5 RETURNING updated_at`
	return db.QueryRow(query, worker.Name, worker.Email, worker.Phone, worker.Role, worker.ID).
		Scan(&worker.UpdatedAt)
}

func DeleteWorker(id int) error {
	query := `DELETE FROM workers WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// Field methods
func CreateField(field *Field) error {
	query := `INSERT INTO fields (name, description, coordinates, area, crop_type, period, region)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  RETURNING id, created_at, updated_at`
	return db.QueryRow(query, field.Name, field.Description, field.Coordinates, field.Area, field.CropType, field.Period, field.Region).
		Scan(&field.ID, &field.CreatedAt, &field.UpdatedAt)
}

func GetFields() ([]Field, error) {
	query := `SELECT id, name, description, coordinates, area, crop_type, period, region, created_at, updated_at FROM fields ORDER BY name`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []Field
	for rows.Next() {
		var f Field
		err := rows.Scan(&f.ID, &f.Name, &f.Description, &f.Coordinates, &f.Area, &f.CropType, &f.Period, &f.Region, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, nil
}

func GetFieldByID(id int) (*Field, error) {
	var f Field
	query := `SELECT id, name, description, coordinates, area, crop_type, period, region, created_at, updated_at FROM fields WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&f.ID, &f.Name, &f.Description, &f.Coordinates, &f.Area, &f.CropType, &f.Period, &f.Region, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func UpdateField(field *Field) error {
	query := `UPDATE fields SET name = $1, description = $2, coordinates = $3, area = $4, crop_type = $5, period = $6, region = $7, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $8 RETURNING updated_at`
	return db.QueryRow(query, field.Name, field.Description, field.Coordinates, field.Area, field.CropType, field.Period, field.Region, field.ID).
		Scan(&field.UpdatedAt)
}

func DeleteField(id int) error {
	query := `DELETE FROM fields WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// Schedule methods
func CreateSchedule(schedule *Schedule) error {
	query := `INSERT INTO schedules (worker_id, date)
			  VALUES ($1, $2)
			  RETURNING id, created_at, updated_at`
	return db.QueryRow(query, schedule.WorkerID, schedule.Date).
		Scan(&schedule.ID, &schedule.CreatedAt, &schedule.UpdatedAt)
}

func GetSchedules() ([]Schedule, error) {
	query := `SELECT s.id, s.worker_id, s.date, s.created_at, s.updated_at, w.name
			  FROM schedules s
			  LEFT JOIN workers w ON s.worker_id = w.id
			  ORDER BY s.date DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var s Schedule
		var workerName sql.NullString
		err := rows.Scan(&s.ID, &s.WorkerID, &s.Date, &s.CreatedAt, &s.UpdatedAt, &workerName)
		if err != nil {
			return nil, err
		}
		if workerName.Valid {
			s.Worker = &Worker{ID: s.WorkerID, Name: workerName.String}
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func GetScheduleByID(id int) (*Schedule, error) {
	var s Schedule
	query := `SELECT s.id, s.worker_id, s.date, s.created_at, s.updated_at, w.name
			  FROM schedules s
			  LEFT JOIN workers w ON s.worker_id = w.id
			  WHERE s.id = $1`
	var workerName sql.NullString
	err := db.QueryRow(query, id).Scan(&s.ID, &s.WorkerID, &s.Date, &s.CreatedAt, &s.UpdatedAt, &workerName)
	if err != nil {
		return nil, err
	}
	if workerName.Valid {
		s.Worker = &Worker{ID: s.WorkerID, Name: workerName.String}
	}
	return &s, nil
}

func GetWorkerSchedules(workerID int) ([]Schedule, error) {
	query := `SELECT s.id, s.worker_id, s.date, s.created_at, s.updated_at, w.name
			  FROM schedules s
			  LEFT JOIN workers w ON s.worker_id = w.id
			  WHERE s.worker_id = $1
			  ORDER BY s.date DESC`
	rows, err := db.Query(query, workerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var s Schedule
		var workerName sql.NullString
		err := rows.Scan(&s.ID, &s.WorkerID, &s.Date, &s.CreatedAt, &s.UpdatedAt, &workerName)
		if err != nil {
			return nil, err
		}
		if workerName.Valid {
			s.Worker = &Worker{ID: s.WorkerID, Name: workerName.String}
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func UpdateSchedule(schedule *Schedule) error {
	query := `UPDATE schedules SET worker_id = $1, date = $2, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $3 RETURNING updated_at`
	return db.QueryRow(query, schedule.WorkerID, schedule.Date, schedule.ID).
		Scan(&schedule.UpdatedAt)
}

func DeleteSchedule(id int) error {
	query := `DELETE FROM schedules WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// Operation methods
func CreateOperation(operation *Operation) error {
	query := `INSERT INTO operations (schedule_id, worker_id, field_id, type, description, status, start_time, end_time, notes)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  RETURNING id, created_at, updated_at`
	return db.QueryRow(query, operation.ScheduleID, operation.WorkerID, operation.FieldID, operation.Type,
		operation.Description, operation.Status, operation.StartTime, operation.EndTime, operation.Notes).
		Scan(&operation.ID, &operation.CreatedAt, &operation.UpdatedAt)
}

func GetOperations() ([]Operation, error) {
	query := `SELECT o.id, o.schedule_id, o.worker_id, o.field_id, o.type, o.description, o.status,
					 o.start_time, o.end_time, o.completed_at, o.notes, o.created_at, o.updated_at,
					 w.name, f.name
			  FROM operations o
			  LEFT JOIN workers w ON o.worker_id = w.id
			  LEFT JOIN fields f ON o.field_id = f.id
			  ORDER BY o.start_time DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operations []Operation
	for rows.Next() {
		var o Operation
		var workerName, fieldName sql.NullString
		err := rows.Scan(&o.ID, &o.ScheduleID, &o.WorkerID, &o.FieldID, &o.Type, &o.Description, &o.Status,
			&o.StartTime, &o.EndTime, &o.CompletedAt, &o.Notes, &o.CreatedAt, &o.UpdatedAt, &workerName, &fieldName)
		if err != nil {
			return nil, err
		}
		if workerName.Valid {
			o.Worker = &Worker{ID: o.WorkerID, Name: workerName.String}
		}
		if fieldName.Valid {
			o.Field = &Field{ID: o.FieldID, Name: fieldName.String}
		}
		operations = append(operations, o)
	}
	return operations, nil
}

func GetOperationByID(id int) (*Operation, error) {
	var o Operation
	query := `SELECT o.id, o.schedule_id, o.worker_id, o.field_id, o.type, o.description, o.status,
					 o.start_time, o.end_time, o.completed_at, o.notes, o.created_at, o.updated_at,
					 w.name, f.name
			  FROM operations o
			  LEFT JOIN workers w ON o.worker_id = w.id
			  LEFT JOIN fields f ON o.field_id = f.id
			  WHERE o.id = $1`
	var workerName, fieldName sql.NullString
	err := db.QueryRow(query, id).Scan(&o.ID, &o.ScheduleID, &o.WorkerID, &o.FieldID, &o.Type, &o.Description, &o.Status,
		&o.StartTime, &o.EndTime, &o.CompletedAt, &o.Notes, &o.CreatedAt, &o.UpdatedAt, &workerName, &fieldName)
	if err != nil {
		return nil, err
	}
	if workerName.Valid {
		o.Worker = &Worker{ID: o.WorkerID, Name: workerName.String}
	}
	if fieldName.Valid {
		o.Field = &Field{ID: o.FieldID, Name: fieldName.String}
	}
	return &o, nil
}

func UpdateOperation(operation *Operation) error {
	query := `UPDATE operations SET schedule_id = $1, worker_id = $2, field_id = $3, type = $4, description = $5,
			  status = $6, start_time = $7, end_time = $8, notes = $9, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $10 RETURNING updated_at`
	return db.QueryRow(query, operation.ScheduleID, operation.WorkerID, operation.FieldID, operation.Type,
		operation.Description, operation.Status, operation.StartTime, operation.EndTime, operation.Notes, operation.ID).
		Scan(&operation.UpdatedAt)
}

func CompleteOperation(id int) error {
	query := `UPDATE operations SET status = 'completed', completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func StartOperation(id int) error {
	query := `UPDATE operations SET status = 'in_progress', start_time = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func RejectOperation(id int) error {
	query := `UPDATE operations SET worker_id = NULL, schedule_id = NULL updated_at = CURRENT_TIMESTAMP
			  WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func DeleteOperation(id int) error {
	query := `DELETE FROM operations WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// Report methods
func GetDailyReport(date time.Time) (*DailyReport, error) {
	report := &DailyReport{
		Date:             date,
		OperationsByType: make(map[string]int),
	}

	// Get total workers active on this date
	err := db.QueryRow(`SELECT COUNT(DISTINCT worker_id) FROM operations WHERE DATE(start_time) = $1`, date).
		Scan(&report.TotalWorkers)
	if err != nil {
		return nil, err
	}

	// Get operation statistics
	err = db.QueryRow(`SELECT COUNT(*),
						COUNT(CASE WHEN status = 'completed' THEN 1 END),
						COUNT(CASE WHEN status = 'in_progress' THEN 1 END)
					   FROM operations WHERE DATE(start_time) = $1`, date).
		Scan(&report.TotalOperations, &report.CompletedOps, &report.InProgressOps)
	if err != nil {
		return nil, err
	}

	// Get operations by type
	rows, err := db.Query(`SELECT type, COUNT(*) FROM operations WHERE DATE(start_time) = $1 GROUP BY type`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var opType string
		var count int
		if err := rows.Scan(&opType, &count); err != nil {
			return nil, err
		}
		report.OperationsByType[opType] = count
	}

	// Get worker statistics
	workerRows, err := db.Query(`
		SELECT o.worker_id, w.name, COUNT(*),
			   COALESCE(SUM(EXTRACT(EPOCH FROM (COALESCE(o.end_time, NOW()) - o.start_time))/3600), 0),
			   COUNT(DISTINCT o.field_id)
		FROM operations o
		JOIN workers w ON o.worker_id = w.id
		WHERE DATE(o.start_time) = $1
		GROUP BY o.worker_id, w.name`, date)
	if err != nil {
		return nil, err
	}
	defer workerRows.Close()

	for workerRows.Next() {
		var ws WorkerDailyStats
		if err := workerRows.Scan(&ws.WorkerID, &ws.WorkerName, &ws.Operations, &ws.HoursWorked, &ws.FieldsWorked); err != nil {
			return nil, err
		}
		report.WorkerStats = append(report.WorkerStats, ws)
	}

	// Get field statistics
	fieldRows, err := db.Query(`
		SELECT o.field_id, f.name, COUNT(*),
			   COALESCE(SUM(EXTRACT(EPOCH FROM (COALESCE(o.end_time, NOW()) - o.start_time))/3600), 0),
			   COUNT(DISTINCT o.worker_id)
		FROM operations o
		JOIN fields f ON o.field_id = f.id
		WHERE DATE(o.start_time) = $1
		GROUP BY o.field_id, f.name`, date)
	if err != nil {
		return nil, err
	}
	defer fieldRows.Close()

	for fieldRows.Next() {
		var fs FieldDailyStats
		if err := fieldRows.Scan(&fs.FieldID, &fs.FieldName, &fs.Operations, &fs.HoursWorked, &fs.WorkersCount); err != nil {
			return nil, err
		}
		report.FieldStats = append(report.FieldStats, fs)
	}

	return report, nil
}
