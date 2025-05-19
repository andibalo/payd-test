package db

import (
	"database/sql"
	"fmt"
	"github.com/andibalo/payd-test/backend/internal/config"
	"github.com/andibalo/payd-test/backend/pkg"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

func InitDB(cfg config.Config) *sql.DB {

	var err error
	db, err := sql.Open("sqlite", cfg.DBConnString())
	if err != nil {
		cfg.Logger().Error("Failed to connect to db", zap.Error(err))
		panic("Failed to connect to db")
	}

	cfg.Logger().Info("Connected to database")

	err = initDbTables(db, cfg)
	if err != nil {
		cfg.Logger().Error("Error init db tables", zap.Error(err))
		panic("Error init db tables")
	}

	if cfg.GetFlags().EnableSeedDB {
		err = seedDatabase(db, cfg)
		if err != nil {
			cfg.Logger().Error("Error seeding db", zap.Error(err))
		}
	}

	return db
}

func initDbTables(db *sql.DB, cfg config.Config) error {

	createUserTableQuery := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
    		email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(100) NOT NULL,
			updated_at TIMESTAMP,
			updated_by VARCHAR(100),
			deleted_at TIMESTAMP,
			deleted_by VARCHAR(100)
	);`

	_, err := db.Exec(createUserTableQuery)

	if err != nil {
		cfg.Logger().Error("Error create users table", zap.Error(err))
		return err
	}

	createShiftRoleEnumQuery := `CREATE TABLE IF NOT EXISTS shift_role_enum (
	 		id INTEGER PRIMARY KEY AUTOINCREMENT,
			role_name VARCHAR(100) NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(100) NOT NULL,
			updated_at TIMESTAMP,
			updated_by VARCHAR(100),
			deleted_at TIMESTAMP,
			deleted_by VARCHAR(100)
	);`

	_, err = db.Exec(createShiftRoleEnumQuery)

	if err != nil {
		cfg.Logger().Error("Error create shift_role_enum table", zap.Error(err))
		return err
	}

	createShiftTableQuery := `CREATE TABLE IF NOT EXISTS shifts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date DATE NOT NULL,
			start_time TIMESTAMP NOT NULL,  
			end_time TIMESTAMP NOT NULL,   
			role_id INTEGER NOT NULL, 
			location VARCHAR(255),              
    		is_active BOOLEAN NOT NULL DEFAULT TRUE, 
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(100) NOT NULL,
			updated_at TIMESTAMP,
			updated_by VARCHAR(100),
			deleted_at TIMESTAMP,
			deleted_by VARCHAR(100),
			FOREIGN KEY (role_id) REFERENCES shift_role_enum(id)
	);`

	_, err = db.Exec(createShiftTableQuery)

	if err != nil {
		cfg.Logger().Error("Error create shifts table", zap.Error(err))
		return err
	}

	createShiftRequestTableQuery := `CREATE TABLE IF NOT EXISTS shift_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		shift_id INTEGER NOT NULL,
		status VARCHAR(100) NOT NULL,     
		requested_by VARCHAR(100) NOT NULL,
		admin_actor VARCHAR(100),
		rejection_reason varchar(255),
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(100) NOT NULL,
		updated_at TIMESTAMP,
		updated_by VARCHAR(100),
		deleted_at TIMESTAMP,
		deleted_by VARCHAR(100),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (shift_id) REFERENCES shifts(id)
	);`

	_, err = db.Exec(createShiftRequestTableQuery)

	if err != nil {
		cfg.Logger().Error("Error create shift_requests table", zap.Error(err))
		return err
	}

	createWorkerShiftAssignmentsTableQuery := `CREATE TABLE IF NOT EXISTS worker_shift_assignments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL, 
		shift_id INTEGER NOT NULL, 
		assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		assigned_by VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(100) NOT NULL,
		updated_at TIMESTAMP,
		updated_by VARCHAR(100),
		deleted_at TIMESTAMP,
		deleted_by VARCHAR(100),
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (shift_id) REFERENCES shifts(shift_id)
	);`

	_, err = db.Exec(createWorkerShiftAssignmentsTableQuery)

	if err != nil {
		cfg.Logger().Error("Error create worker_shift_assignments table", zap.Error(err))
		return err
	}

	createWorkerShiftAvailabilitiesTableQuery := `CREATE TABLE IF NOT EXISTS worker_shift_availabilites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL, 
		available_date_from DATE NOT NULL,
		available_date_to DATE NOT NULL,
		available_time_from TIMESTAMP, 
		available_time_to TIMESTAMP,   
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);`

	_, err = db.Exec(createWorkerShiftAvailabilitiesTableQuery)

	if err != nil {
		cfg.Logger().Error("Error create worker_shift_availabilites table", zap.Error(err))
		return err
	}

	return nil
}

func seedDatabase(db *sql.DB, cfg config.Config) error {
	roles := []struct {
		RoleName  string
		CreatedBy string
		UpdatedBy string
		DeletedBy string
	}{
		{"Cleaner", "system", "", ""},
		{"Cook", "system", "", ""},
		{"Mover", "system", "", ""},
	}

	var existingRoleCount int

	err := db.QueryRow(`SELECT COUNT(*) FROM shift_role_enum`).Scan(&existingRoleCount)
	if err != nil {
		cfg.Logger().Error("Error checking for existing roles:", zap.Error(err))
		return err
	}

	if existingRoleCount == 0 {
		for _, role := range roles {
			query := `
				INSERT INTO shift_role_enum (role_name, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by)
				VALUES (?, CURRENT_TIMESTAMP, ?, ?, ?, ?, ?)
			`
			_, err := db.Exec(query, role.RoleName, role.CreatedBy, nil, nil, nil, nil)
			if err != nil {
				cfg.Logger().Error("Error inserting role :", zap.String("role", role.RoleName), zap.Error(err))
				return err
			} else {
				cfg.Logger().Info(fmt.Sprintf("Successfully inserted role: %s", role.RoleName))
			}
		}

	} else {
		cfg.Logger().Info("Roles already exist.")
	}

	var existingUserCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&existingUserCount)
	if err != nil {
		cfg.Logger().Error("Error checking for existing users:", zap.Error(err))
		return err
	}

	if existingUserCount == 0 {
		users := []struct {
			FirstName, LastName, Email, Password, Role string
		}{
			{"Alice", "Smith", "alice@example.com", "123", "WORKER"},
			{"Bob", "Johnson", "bob@example.com", "123", "WORKER"},
			{"Charlie", "Brown", "charlie@example.com", "123", "WORKER"},
			{"Admin", "Test", "admin@example.com", "123", "ADMIN"},
		}

		for _, user := range users {
			hasedPassword, err := pkg.HashPassword(user.Password)
			if err != nil {
				cfg.Logger().Error("Failed to hash password", zap.Error(err))

				return err
			}

			userQuery := `
				INSERT INTO users (first_name, last_name, email, password, role, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by)
				VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?, ?, ?, ?,?)
			`
			_, err = db.Exec(userQuery, user.FirstName, user.LastName, user.Email, hasedPassword, user.Role, "system", nil, nil, nil, nil)
			if err != nil {
				cfg.Logger().Error("Error inserting user:", zap.Error(err))
				return err
			}
			cfg.Logger().Info("Successfully inserted user: " + user.Email)
		}
	} else {
		cfg.Logger().Info("Users already exist.")
	}

	var existingShiftCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM shifts`).Scan(&existingShiftCount)
	if err != nil {
		cfg.Logger().Error("Error checking for existing shifts:", zap.Error(err))
		return err
	}

	if existingShiftCount == 0 {
		shifts := []struct {
			Date      string
			StartTime string
			EndTime   string
			RoleID    int
			Location  string
		}{
			{"2025-05-01", "08:00:00", "16:00:00", 1, "Location A"},
			{"2025-05-02", "09:00:00", "17:00:00", 2, "Location B"},
			{"2025-05-03", "10:00:00", "18:00:00", 3, "Location C"},
		}

		for _, shift := range shifts {
			shiftQuery := `
				INSERT INTO shifts (date, start_time, end_time, role_id, location, is_active, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by)
				VALUES (?, ?, ?, ?, ?, TRUE, CURRENT_TIMESTAMP, ?, ?, ?, ?, ?)
			`
			_, err := db.Exec(shiftQuery, shift.Date, shift.StartTime, shift.EndTime, shift.RoleID, shift.Location, "system", nil, nil, nil, nil)
			if err != nil {
				cfg.Logger().Error("Error inserting shift:", zap.Error(err))
				return err
			}
			cfg.Logger().Info("Successfully inserted shift: " + shift.Location)
		}
	} else {
		cfg.Logger().Info("Shifts already exist.")
	}

	var existingShiftRequestCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM shift_requests`).Scan(&existingShiftRequestCount)
	if err != nil {
		cfg.Logger().Error("Error checking for existing shift requests:", zap.Error(err))
		return err
	}

	if existingShiftRequestCount == 0 {
		shiftRequests := []struct {
			UserID      int
			ShiftID     int
			Status      string
			RequestedBy string
		}{
			{1, 1, "PENDING", "alice@example.com"},
			{2, 2, "REJECTED", "bob@example.com"},
			{3, 3, "APPROVED", "charlie@example.com"},
		}

		for _, request := range shiftRequests {
			shiftRequestQuery := `
				INSERT INTO shift_requests (user_id, shift_id, status, requested_by, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by)
				VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, ?, ?, ?, ?, ?)
			`
			_, err := db.Exec(shiftRequestQuery, request.UserID, request.ShiftID, request.Status, request.RequestedBy, "system", nil, nil, nil, nil)
			if err != nil {
				cfg.Logger().Error("Error inserting shift request:", zap.Error(err))
				return err
			}
			cfg.Logger().Info("Successfully inserted shift request.")
		}
	} else {
		cfg.Logger().Info("Shift requests already exist.")
	}

	var existingWorkerShiftAssignmentCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM worker_shift_assignments`).Scan(&existingWorkerShiftAssignmentCount)
	if err != nil {
		cfg.Logger().Error("Error checking for existing worker shift assignments:", zap.Error(err))
		return err
	}

	if existingWorkerShiftAssignmentCount == 0 {
		workerShiftAssignments := []struct {
			UserID  int
			ShiftID int
		}{
			{3, 3},
		}

		for _, assignment := range workerShiftAssignments {
			workerShiftAssignmentQuery := `
				INSERT INTO worker_shift_assignments (user_id, shift_id, assigned_at, assigned_by, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by)
				VALUES (?, ?, CURRENT_TIMESTAMP, ?, CURRENT_TIMESTAMP, ?, ?, ?, ?, ?)
			`
			_, err := db.Exec(workerShiftAssignmentQuery, assignment.UserID, assignment.ShiftID, "admin@example.com", "system", nil, nil, nil, nil)
			if err != nil {
				cfg.Logger().Error("Error inserting worker shift assignment:", zap.Error(err))
				return err
			}
			cfg.Logger().Info("Successfully inserted worker shift assignment.")
		}
	} else {
		cfg.Logger().Info("Worker shift assignments already exist.")
	}

	var existingWorkerShiftAvailabilityCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM worker_shift_availabilites`).Scan(&existingWorkerShiftAvailabilityCount)
	if err != nil {
		cfg.Logger().Error("Error checking for existing worker shift availabilities:", zap.Error(err))
		return err
	}

	if existingWorkerShiftAvailabilityCount == 0 {
		workerShiftAvailabilities := []struct {
			UserID   int
			FromDate string
			ToDate   string
		}{
			{1, "2025-05-01", "2025-05-02"},
			{2, "2025-05-03", "2025-05-04"},
			{3, "2025-05-05", "2025-05-06"},
		}

		for _, availability := range workerShiftAvailabilities {
			workerShiftAvailabilityQuery := `
				INSERT INTO worker_shift_availabilites (user_id, available_date_from, available_date_to, available_time_from, available_time_to)
				VALUES (?, ?, ?, ?, ?)
			`
			_, err := db.Exec(workerShiftAvailabilityQuery, availability.UserID, availability.FromDate, availability.ToDate, "08:00:00", "16:00:00")
			if err != nil {
				cfg.Logger().Error("Error inserting worker shift availability:", zap.Error(err))
				return err
			}
			cfg.Logger().Info("Successfully inserted worker shift availability.")
		}
	} else {
		cfg.Logger().Info("Worker shift availabilities already exist.")
	}

	return nil
}
