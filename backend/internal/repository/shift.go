package repository

import (
	"database/sql"
	"github.com/andibalo/payd-test/backend/internal/constants"
	"github.com/andibalo/payd-test/backend/internal/model"
	"github.com/andibalo/payd-test/backend/internal/response"
	"github.com/andibalo/payd-test/backend/pkg/httpresp"
	"time"
)

type shiftRepository struct {
	db *sql.DB
}

func NewShiftRepository(db *sql.DB) ShiftRepository {
	return &shiftRepository{
		db: db,
	}
}

type GetShiftListFilter struct {
	ShowOnlyUnassigned bool `json:"show_only_unassigned"`
	Limit              int  `json:"limit"`
	Offset             int  `json:"offset"`
}

type GetShiftRequestListFilter struct {
	UserID int64  `json:"user_id"`
	Status string `json:"status"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type GetShiftAssignmentListFilter struct {
	UserID int64 `json:"user_id"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

func (r *shiftRepository) Save(shift *model.Shift) error {
	query := `
		INSERT INTO shifts (date, start_time, end_time, role_id, location, is_active, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	_, err := r.db.Exec(query, shift.Date, shift.StartTime, shift.EndTime, shift.RoleID, shift.Location, shift.IsActive, shift.CreatedBy)
	return err
}

func (r *shiftRepository) GetByID(id int64) (*model.Shift, error) {
	shift := &model.Shift{}

	query := `
		SELECT id, date, start_time, end_time, role_id, location, created_by, created_at, 
			updated_by, updated_at, deleted_by, deleted_at
		FROM shifts 
		WHERE id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	err := r.db.QueryRow(query, id).Scan(
		&shift.ID, &shift.Date, &shift.StartTime, &shift.EndTime, &shift.RoleID, &shift.Location,
		&shift.CreatedBy, &shift.CreatedAt, &shift.UpdatedBy, &shift.UpdatedAt, &shift.DeletedBy, &shift.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return shift, nil
}

func (r *shiftRepository) GetList(filter GetShiftListFilter) ([]response.GetShiftListData, *httpresp.Pagination, error) {
	shifts := []response.GetShiftListData{}

	query := `
		SELECT 
			shifts.id, 
			shifts.date, 
			shifts.start_time, 
			shifts.end_time, 
			shifts.role_id, 
			shift_role_enum.role_name, 
			shifts.location, 
			shifts.created_by, 
			shifts.created_at, 
			shifts.updated_by, 
			shifts.updated_at, 
			shifts.deleted_by, 
			shifts.deleted_at
		FROM shifts
		JOIN shift_role_enum ON shifts.role_id = shift_role_enum.id
		WHERE shifts.deleted_at IS NULL
	`
	var args []interface{}

	if filter.ShowOnlyUnassigned {
		query += `
			AND shifts.id NOT IN (
				SELECT shift_id 
				FROM worker_shift_assignments 
				WHERE worker_shift_assignments.deleted_at IS NULL
			)
		`
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	query += " ORDER BY shifts.created_at DESC"

	query += " LIMIT ? OFFSET ?"

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var shift response.GetShiftListData
		err := rows.Scan(
			&shift.ID, &shift.Date, &shift.StartTime, &shift.EndTime, &shift.RoleID, &shift.RoleName, &shift.Location,
			&shift.CreatedBy, &shift.CreatedAt, &shift.UpdatedBy, &shift.UpdatedAt, &shift.DeletedBy, &shift.DeletedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		shifts = append(shifts, shift)
	}

	countQuery := `
		SELECT COUNT(*) 
		FROM shifts
		JOIN shift_role_enum ON shifts.role_id = shift_role_enum.id
		WHERE shifts.deleted_at IS NULL
	`

	if filter.ShowOnlyUnassigned {
		countQuery += `
			AND shifts.id NOT IN (
				SELECT shift_id 
				FROM worker_shift_assignments 
				WHERE worker_shift_assignments.deleted_at IS NULL
			)
		`
	}

	var totalCount int64
	countRow := r.db.QueryRow(countQuery, args...)
	err = countRow.Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (totalCount + int64(filter.Limit-1)) / int64(filter.Limit)

	pagination := &httpresp.Pagination{
		CurrentPage:   int64(filter.Offset/filter.Limit + 1),
		TotalPages:    totalPages,
		TotalElements: totalCount,
		SortBy:        "created_at",
	}

	return shifts, pagination, nil
}

func (r *shiftRepository) UpdateByID(id int64, shift *model.Shift) error {
	query := `
		UPDATE shifts 
		SET date = ?, start_time = ?, end_time = ?, role_id = ?, location = ?, 
			updated_by = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, shift.Date, shift.StartTime, shift.EndTime, shift.RoleID, shift.Location, shift.UpdatedBy, id)
	return err
}

func (r *shiftRepository) DeleteByID(id int64, deletedBy string) error {
	query := `
		UPDATE shifts 
		SET deleted_at = CURRENT_TIMESTAMP, deleted_by = ? 
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, deletedBy, id)
	return err
}

func (r *shiftRepository) GetShiftRequestList(filter GetShiftRequestListFilter) ([]response.GetShiftRequestListData, *httpresp.Pagination, error) {
	shiftRequests := []response.GetShiftRequestListData{}

	query := `
		SELECT sr.id, sr.user_id, sr.shift_id, s.date AS shift_date, s.start_time AS shift_start_time, 
			s.end_time AS shift_end_time, s.role_id AS shift_role_id, sre.role_name AS shift_role_name,
			sr.status, sr.requested_by, sr.admin_actor, sr.rejection_reason, 
			sr.created_at, sr.created_by, sr.updated_at, sr.updated_by, sr.deleted_at, sr.deleted_by
		FROM shift_requests sr
		JOIN shifts s ON sr.shift_id = s.id
		JOIN shift_role_enum sre ON s.role_id = sre.id
		WHERE sr.deleted_at IS NULL
	`

	var args []interface{}

	if filter.Status != "" {
		query += " AND sr.status = ?"
		args = append(args, filter.Status)
	}

	if filter.UserID != 0 {
		query += " AND sr.user_id = ?"
		args = append(args, filter.UserID)
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	query += " ORDER BY sr.created_at DESC"

	query += " LIMIT ? OFFSET ?"

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var shiftRequest response.GetShiftRequestListData

		err := rows.Scan(
			&shiftRequest.ID, &shiftRequest.UserID, &shiftRequest.ShiftID, &shiftRequest.ShiftDate,
			&shiftRequest.ShiftStartTime, &shiftRequest.ShiftEndTime, &shiftRequest.ShiftRoleID,
			&shiftRequest.ShiftRoleName, &shiftRequest.Status, &shiftRequest.RequestedBy,
			&shiftRequest.AdminActor, &shiftRequest.RejectionReason, &shiftRequest.CreatedAt,
			&shiftRequest.CreatedBy, &shiftRequest.UpdatedAt, &shiftRequest.UpdatedBy,
			&shiftRequest.DeletedAt, &shiftRequest.DeletedBy,
		)
		if err != nil {
			return nil, nil, err
		}

		shiftRequests = append(shiftRequests, shiftRequest)
	}

	countQuery := `
		SELECT COUNT(*) FROM shift_requests sr
		JOIN shifts s ON sr.shift_id = s.id
		JOIN shift_role_enum sre ON s.role_id = sre.id
		WHERE sr.deleted_at IS NULL
	`

	var totalCount int64
	countRows := r.db.QueryRow(countQuery, args...)
	err = countRows.Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (totalCount + int64(filter.Limit-1)) / int64(filter.Limit)

	pagination := &httpresp.Pagination{
		CurrentPage:   int64(filter.Offset/filter.Limit + 1),
		TotalPages:    totalPages,
		TotalElements: totalCount,
		SortBy:        "created_at",
	}

	return shiftRequests, pagination, nil
}

func (r *shiftRepository) GetShiftRequestByID(id int64) (*model.ShiftRequest, error) {
	sr := &model.ShiftRequest{}

	// Raw SQL query to fetch the shift request by ID
	query := `
		SELECT id, user_id, shift_id, status, requested_by, admin_actor, rejection_reason, 
			created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM shift_requests 
		WHERE id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	err := r.db.QueryRow(query, id).Scan(
		&sr.ID, &sr.UserID, &sr.ShiftID, &sr.Status,
		&sr.RequestedBy, &sr.AdminActor, &sr.RejectionReason,
		&sr.CreatedAt, &sr.CreatedBy, &sr.UpdatedAt,
		&sr.UpdatedBy, &sr.DeletedAt, &sr.DeletedBy,
	)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

func (r *shiftRepository) SaveShiftRequest(shiftRequest *model.ShiftRequest) error {
	query := `
		INSERT INTO shift_requests (
			user_id, shift_id, status, requested_by, admin_actor, rejection_reason, created_at, created_by, 
			updated_at, updated_by, deleted_at, deleted_by
		) 
		VALUES (?, ?, ?, ?, NULL, NULL, CURRENT_TIMESTAMP, ?, NULL, NULL, NULL, NULL)
	`

	_, err := r.db.Exec(query, shiftRequest.UserID, shiftRequest.ShiftID, shiftRequest.Status, shiftRequest.RequestedBy, shiftRequest.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (r *shiftRepository) SaveWorkerShift(workerShift *model.WorkerShift) error {
	query := `
		INSERT INTO worker_shift_assignments (
			user_id, shift_id, assigned_at, assigned_by, created_at, created_by, 
			updated_at, updated_by, deleted_at, deleted_by
		) 
		VALUES (?, ?, CURRENT_TIMESTAMP, ?, CURRENT_TIMESTAMP, ?, NULL, NULL, NULL, NULL)
	`

	_, err := r.db.Exec(query, workerShift.UserID, workerShift.ShiftID, workerShift.AssignedBy, workerShift.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (r *shiftRepository) UpdateShiftRequestByID(id int64, sr *model.ShiftRequest) error {
	query := `
		UPDATE shift_requests 
		SET status = ?, admin_actor = ?, rejection_reason = ?,
			updated_by = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, sr.Status, sr.AdminActor, sr.RejectionReason, sr.UpdatedBy, id)
	return err
}

func (r *shiftRepository) CheckUserAssignedShiftExistsByDate(userID int64, shiftDate time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM worker_shift_assignments wsa
		JOIN shifts s ON wsa.shift_id = s.id
		WHERE wsa.user_id = ? 
		AND s.date = ? 
		AND wsa.deleted_at IS NULL
	`
	var assignedShiftCount int
	err := r.db.QueryRow(query, userID, shiftDate).Scan(&assignedShiftCount)
	if err != nil {
		return false, err
	}

	if assignedShiftCount > 0 {
		return true, nil
	}

	return false, nil
}

func (r *shiftRepository) CheckIfShiftIsAlreadyAssigned(shiftID int64) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM shifts s
		JOIN worker_shift_assignments wsa ON wsa.shift_id = s.id
		WHERE s.id = ?
		AND   s.deleted_at IS NULL
	`

	var assignedShiftCount int
	err := r.db.QueryRow(query, shiftID).Scan(&assignedShiftCount)
	if err != nil {
		return false, err
	}

	if assignedShiftCount > 0 {
		return true, nil
	}

	return false, nil
}

func (r *shiftRepository) CheckIfShiftRequestTimeOverlaps(userID int64, shiftDate, requestedStartTime, requestedEndTime time.Time) (bool, error) {
	checkOverlapQuery := `
		SELECT COUNT(*) 
		FROM shift_requests sr
		JOIN shifts s ON sr.shift_id = s.id
		WHERE sr.user_id = ? 
		AND sr.status IN (?) 
		AND s.date = ?
		AND sr.deleted_at IS NULL
		AND s.starttime <= ?
    	AND s.endtime >= ?
	`

	var overlapCount int
	err := r.db.QueryRow(
		checkOverlapQuery,
		userID,
		constants.SHIFT_REQUEST_STATUS_PENDING,
		shiftDate,
		requestedEndTime,
		requestedStartTime).Scan(&overlapCount)
	if err != nil {
		return false, err
	}

	if overlapCount > 0 {
		return true, nil
	}

	return false, nil
}

func (r *shiftRepository) GetUserWeeklyAssignedShiftCountByDate(userID int64, shiftDate time.Time) (int, error) {
	checkWeeklyShiftsQuery := `
		SELECT COUNT(*) 
		FROM shift_requests sr
		JOIN shifts s ON sr.shift_id = s.id
		WHERE sr.user_id = ? 
		AND sr.status IN (?) 
		AND strftime('%W', s.date) = strftime('%W', ?) 
		AND sr.deleted_at IS NULL
	`

	var weeklyShiftCount int
	err := r.db.QueryRow(checkWeeklyShiftsQuery, userID, constants.SHIFT_REQUEST_STATUS_APPROVED, shiftDate).Scan(&weeklyShiftCount)
	if err != nil {
		return 0, err
	}

	return weeklyShiftCount, nil
}

func (r *shiftRepository) GetShiftAssignmentList(filter GetShiftAssignmentListFilter) ([]response.GetShiftAssignmentListData, *httpresp.Pagination, error) {
	assignments := []response.GetShiftAssignmentListData{}

	query := `
		SELECT 
			worker_shift_assignments.id, 
			worker_shift_assignments.user_id, 
			users.first_name, 
			users.last_name, 
			users.email, 
			worker_shift_assignments.shift_id, 
			shifts.date, 
			shifts.start_time, 
			shifts.end_time, 
			shift_role_enum.role_name, 
			worker_shift_assignments.assigned_at, 
			worker_shift_assignments.assigned_by, 
			worker_shift_assignments.created_at, 
			worker_shift_assignments.created_by, 
			worker_shift_assignments.updated_at, 
			worker_shift_assignments.updated_by, 
			worker_shift_assignments.deleted_at, 
			worker_shift_assignments.deleted_by
		FROM worker_shift_assignments
		JOIN users ON worker_shift_assignments.user_id = users.id
		JOIN shifts ON worker_shift_assignments.shift_id = shifts.id
		JOIN shift_role_enum ON shifts.role_id = shift_role_enum.id
		WHERE worker_shift_assignments.deleted_at IS NULL
	`

	var args []interface{}

	if filter.UserID > 0 {
		query += ` AND worker_shift_assignments.user_id = ?`
		args = append(args, filter.UserID)
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	query += " ORDER BY worker_shift_assignments.created_at DESC"

	query += " LIMIT ? OFFSET ?"

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var assignment response.GetShiftAssignmentListData
		err := rows.Scan(
			&assignment.ID, &assignment.UserID, &assignment.FirstName, &assignment.LastName, &assignment.Email,
			&assignment.ShiftID, &assignment.ShiftDate, &assignment.ShiftStartTime, &assignment.ShiftEndTime,
			&assignment.ShiftRoleName, &assignment.AssignedAt, &assignment.AssignedBy,
			&assignment.CreatedAt, &assignment.CreatedBy, &assignment.UpdatedAt, &assignment.UpdatedBy,
			&assignment.DeletedAt, &assignment.DeletedBy,
		)
		if err != nil {
			return nil, nil, err
		}
		assignments = append(assignments, assignment)
	}

	countQuery := `
		SELECT COUNT(*) 
		FROM worker_shift_assignments
		JOIN users ON worker_shift_assignments.user_id = users.id
		JOIN shifts ON worker_shift_assignments.shift_id = shifts.id
		JOIN shift_role_enum ON shifts.role_id = shift_role_enum.id
		WHERE worker_shift_assignments.deleted_at IS NULL
	`

	if filter.UserID > 0 {
		countQuery += ` AND worker_shift_assignments.user_id = ?`
		args = append(args, filter.UserID)
	}

	var totalCount int64
	countRows := r.db.QueryRow(countQuery, args...)
	err = countRows.Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (totalCount + int64(filter.Limit-1)) / int64(filter.Limit)

	pagination := &httpresp.Pagination{
		CurrentPage:   int64(filter.Offset/filter.Limit + 1),
		TotalPages:    totalPages,
		TotalElements: totalCount,
		SortBy:        "created_at",
	}

	return assignments, pagination, nil
}
