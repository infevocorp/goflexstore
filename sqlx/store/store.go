// Package sqlxstore implements store.Store[Entity, ID] backed by jmoiron/sqlx.
package sqlxstore

import (
	"context"
	"database/sql"
	stderrs "errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/infevocorp/goflexstore/converter"
	"github.com/infevocorp/goflexstore/query"
	"github.com/infevocorp/goflexstore/store"
	sqlxopscope "github.com/infevocorp/goflexstore/sqlx/opscope"
	sqlxquery "github.com/infevocorp/goflexstore/sqlx/query"
	sqlxutils "github.com/infevocorp/goflexstore/sqlx/utils"
)

// ErrPreloadNotSupported is returned when a PreloadParam is passed to Get or List.
var ErrPreloadNotSupported = stderrs.New("preload is not supported in sqlx store")

// Store implements store.Store[Entity, ID] using raw SQL via jmoiron/sqlx.
//
// Entity is the clean domain model; DTO is the sqlx-tagged database struct. ID is the primary-key type.
type Store[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable] struct {
	OpScope      *sqlxopscope.TransactionScope
	Converter    converter.Converter[Entity, DTO, ID]
	QueryBuilder *sqlxquery.Builder
	Table        string
	BatchSize    int
	PKColumn     string
	Dialect      sqlxquery.Dialect
	ReturningID  bool
}

// New constructs a Store, applying options and supplying defaults for any
// unset fields.
func New[Entity store.Entity[ID], DTO store.Entity[ID], ID comparable](
	opScope *sqlxopscope.TransactionScope,
	options ...Option[Entity, DTO, ID],
) *Store[Entity, DTO, ID] {
	s := &Store[Entity, DTO, ID]{
		OpScope:   opScope,
		BatchSize: 50,
		PKColumn:  "id",
		Dialect:   sqlxquery.DialectMySQL,
	}

	for _, opt := range options {
		opt(s)
	}

	if s.Converter == nil {
		s.Converter = converter.NewReflect[Entity, DTO, ID](nil)
	}

	if s.QueryBuilder == nil {
		s.QueryBuilder = sqlxquery.NewBuilder(
			sqlxquery.WithFieldToColMap(sqlxutils.FieldToColMap(*new(DTO))),
		)
	}

	if s.Table == "" {
		s.Table = sqlxutils.TableName(*new(DTO))
	}

	return s
}

// ---- read methods -----------------------------------------------------------

// Get returns the first row matching params, mapping it to Entity.
// Returns store.ErrorNotFound when no row matches.
func (s *Store[Entity, DTO, ID]) Get(ctx context.Context, params ...query.Param) (Entity, error) {
	if err := checkPreload(params); err != nil {
		return *new(Entity), err
	}

	result := s.QueryBuilder.Build(query.NewParams(params...))
	sqlStr := s.buildSelectSQL(result, 1)

	sqlStr, args, err := sqlx.In(sqlStr, result.Args...)
	if err != nil {
		return *new(Entity), err
	}

	db := s.OpScope.Tx(ctx)
	sqlStr = db.Rebind(sqlStr)

	var dto DTO
	if err := sqlx.GetContext(ctx, db, &dto, sqlStr, args...); err != nil {
		if stderrs.Is(err, sql.ErrNoRows) {
			return *new(Entity), store.ErrorNotFound
		}
		return *new(Entity), err
	}

	return s.Converter.ToEntity(dto), nil
}

// List returns all rows matching params.
func (s *Store[Entity, DTO, ID]) List(ctx context.Context, params ...query.Param) ([]Entity, error) {
	if err := checkPreload(params); err != nil {
		return nil, err
	}

	result := s.QueryBuilder.Build(query.NewParams(params...))
	sqlStr := s.buildSelectSQL(result, 0)

	sqlStr, args, err := sqlx.In(sqlStr, result.Args...)
	if err != nil {
		return nil, err
	}

	db := s.OpScope.Tx(ctx)
	sqlStr = db.Rebind(sqlStr)

	var dtos []DTO
	if err := sqlx.SelectContext(ctx, db, &dtos, sqlStr, args...); err != nil {
		return nil, err
	}

	return converter.ToMany(dtos, s.Converter.ToEntity), nil
}

// Count returns the number of rows matching params.
func (s *Store[Entity, DTO, ID]) Count(ctx context.Context, params ...query.Param) (int64, error) {
	result := s.QueryBuilder.Build(query.NewParams(params...))

	sqlStr := "SELECT COUNT(*) FROM " + s.Table
	if result.Where != "" {
		sqlStr += " WHERE " + result.Where
	}

	sqlStr, args, err := sqlx.In(sqlStr, result.Args...)
	if err != nil {
		return 0, err
	}

	db := s.OpScope.Tx(ctx)
	sqlStr = db.Rebind(sqlStr)

	var count int64
	if err := db.QueryRowxContext(ctx, sqlStr, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// Exists returns true when at least one row matches params.
func (s *Store[Entity, DTO, ID]) Exists(ctx context.Context, params ...query.Param) (bool, error) {
	result := s.QueryBuilder.Build(query.NewParams(params...))

	inner := "SELECT 1 FROM " + s.Table
	if result.Where != "" {
		inner += " WHERE " + result.Where
	}
	sqlStr := "SELECT EXISTS(" + inner + ")"

	sqlStr, args, err := sqlx.In(sqlStr, result.Args...)
	if err != nil {
		return false, err
	}

	db := s.OpScope.Tx(ctx)
	sqlStr = db.Rebind(sqlStr)

	var exists bool
	if err := db.QueryRowxContext(ctx, sqlStr, args...).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

// ---- write methods ----------------------------------------------------------

// Create inserts entity and returns its ID.
// When the PK field is zero it is omitted from the INSERT and the database
// auto-increment value is retrieved via LastInsertId (MySQL / SQLite) or a
// RETURNING clause (Postgres, requires WithReturningID(true)).
func (s *Store[Entity, DTO, ID]) Create(ctx context.Context, entity Entity) (ID, error) {
	dto := s.Converter.ToDTO(entity)
	isZeroID := dto.GetID() == *new(ID)

	excludeCols := map[string]bool{}
	if isZeroID {
		excludeCols[s.PKColumn] = true
	}

	cols, vals := getStructColVals(dto, excludeCols, false)
	if len(cols) == 0 {
		return *new(ID), stderrs.New("no columns to insert")
	}

	sqlStr := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		s.Table,
		strings.Join(cols, ", "),
		strings.Join(buildPlaceholders(len(cols)), ", "),
	)

	db := s.OpScope.Tx(ctx)

	if s.ReturningID {
		sqlStr = db.Rebind(sqlStr + " RETURNING " + s.PKColumn)
		var id ID
		if err := db.QueryRowxContext(ctx, sqlStr, vals...).Scan(&id); err != nil {
			return *new(ID), err
		}
		return id, nil
	}

	sqlStr = db.Rebind(sqlStr)
	res, err := db.ExecContext(ctx, sqlStr, vals...)
	if err != nil {
		return *new(ID), err
	}

	if isZeroID {
		if lastID, err2 := res.LastInsertId(); err2 == nil && lastID != 0 {
			setPKField(&dto, s.PKColumn, lastID)
		}
	}

	return dto.GetID(), nil
}

// CreateMany batch-inserts entities in groups of BatchSize.
func (s *Store[Entity, DTO, ID]) CreateMany(ctx context.Context, entities []Entity) error {
	if len(entities) == 0 {
		return nil
	}

	dtos := converter.ToMany(entities, s.Converter.ToDTO)

	isZeroID := dtos[0].GetID() == *new(ID)
	excludeCols := map[string]bool{}
	if isZeroID {
		excludeCols[s.PKColumn] = true
	}

	firstCols, _ := getStructColVals(dtos[0], excludeCols, false)
	if len(firstCols) == 0 {
		return stderrs.New("no columns to insert")
	}

	batchSize := defaultValue(s.BatchSize, 50)
	db := s.OpScope.Tx(ctx)

	for i := 0; i < len(dtos); i += batchSize {
		end := i + batchSize
		if end > len(dtos) {
			end = len(dtos)
		}
		batch := dtos[i:end]

		rowPHs := make([]string, len(batch))
		var allVals []any

		for j, dto := range batch {
			dtoCopy := dto
			_, rowVals := getStructColVals(dtoCopy, excludeCols, false)
			rowPHs[j] = "(" + strings.Join(buildPlaceholders(len(firstCols)), ", ") + ")"
			allVals = append(allVals, rowVals...)
		}

		sqlStr := db.Rebind(fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES %s",
			s.Table,
			strings.Join(firstCols, ", "),
			strings.Join(rowPHs, ", "),
		))

		if _, err := db.ExecContext(ctx, sqlStr, allVals...); err != nil {
			return err
		}
	}

	return nil
}

// Update sets all columns (including zero values) for rows matched by params.
// Falls back to WHERE pk = entity.GetID() when params is empty.
func (s *Store[Entity, DTO, ID]) Update(ctx context.Context, entity Entity, params ...query.Param) error {
	dto := s.Converter.ToDTO(entity)
	id := dto.GetID()

	if id == *new(ID) && len(params) == 0 {
		return stderrs.New("id or query params required for Update")
	}

	setCols, setVals := getStructColVals(dto, map[string]bool{s.PKColumn: true}, false)
	if len(setCols) == 0 {
		return stderrs.New("no columns to update")
	}

	setParts := make([]string, len(setCols))
	for i, col := range setCols {
		setParts[i] = col + " = ?"
	}

	whereStr, whereArgs := s.buildWhere(id, params)

	allArgs := append(setVals, whereArgs...)

	sqlStr := "UPDATE " + s.Table + " SET " + strings.Join(setParts, ", ")
	if whereStr != "" {
		sqlStr += " WHERE " + whereStr
	}

	sqlStr, allArgs, err := sqlx.In(sqlStr, allArgs...)
	if err != nil {
		return err
	}

	db := s.OpScope.Tx(ctx)
	sqlStr = db.Rebind(sqlStr)

	_, err = db.ExecContext(ctx, sqlStr, allArgs...)
	return err
}

// PartialUpdate sets only non-zero fields for rows matched by params.
// Falls back to WHERE pk = entity.GetID() when params is empty.
func (s *Store[Entity, DTO, ID]) PartialUpdate(ctx context.Context, entity Entity, params ...query.Param) error {
	dto := s.Converter.ToDTO(entity)
	id := dto.GetID()

	setCols, setVals := getStructColVals(dto, map[string]bool{s.PKColumn: true}, true)
	if len(setCols) == 0 {
		return stderrs.New("no non-zero columns to update")
	}

	setParts := make([]string, len(setCols))
	for i, col := range setCols {
		setParts[i] = col + " = ?"
	}

	if id == *new(ID) && len(params) == 0 {
		return stderrs.New("id or query params required for PartialUpdate")
	}

	whereStr, whereArgs := s.buildWhere(id, params)

	allArgs := append(setVals, whereArgs...)

	sqlStr := "UPDATE " + s.Table + " SET " + strings.Join(setParts, ", ")
	if whereStr != "" {
		sqlStr += " WHERE " + whereStr
	}

	sqlStr, allArgs, err := sqlx.In(sqlStr, allArgs...)
	if err != nil {
		return err
	}

	db := s.OpScope.Tx(ctx)
	sqlStr = db.Rebind(sqlStr)

	_, err = db.ExecContext(ctx, sqlStr, allArgs...)
	return err
}

// Delete removes rows matched by params.
func (s *Store[Entity, DTO, ID]) Delete(ctx context.Context, params ...query.Param) error {
	result := s.QueryBuilder.Build(query.NewParams(params...))

	sqlStr := "DELETE FROM " + s.Table
	if result.Where != "" {
		sqlStr += " WHERE " + result.Where
	}

	sqlStr, args, err := sqlx.In(sqlStr, result.Args...)
	if err != nil {
		return err
	}

	db := s.OpScope.Tx(ctx)
	sqlStr = db.Rebind(sqlStr)

	_, err = db.ExecContext(ctx, sqlStr, args...)
	return err
}

// Upsert inserts or updates entity according to the OnConflict strategy.
// The dialect-specific SQL is selected based on Store.Dialect.
func (s *Store[Entity, DTO, ID]) Upsert(ctx context.Context, entity Entity, onConflict store.OnConflict) (ID, error) {
	dto := s.Converter.ToDTO(entity)
	isZeroID := dto.GetID() == *new(ID)

	excludeCols := map[string]bool{}
	if isZeroID {
		excludeCols[s.PKColumn] = true
	}

	cols, insertVals := getStructColVals(dto, excludeCols, false)
	if len(cols) == 0 {
		return *new(ID), stderrs.New("no columns to upsert")
	}

	phs := buildPlaceholders(len(cols))
	allArgs := make([]any, len(insertVals))
	copy(allArgs, insertVals)

	var sqlStr string

	switch s.Dialect {
	case sqlxquery.DialectPostgres:
		sqlStr = buildPostgresUpsert(s.Table, cols, phs, onConflict, s.PKColumn, &allArgs)
	case sqlxquery.DialectSQLite:
		sqlStr = buildSQLiteUpsert(s.Table, cols, phs, onConflict, s.PKColumn, &allArgs)
	default:
		sqlStr = buildMySQLUpsert(s.Table, cols, phs, onConflict, s.PKColumn, &allArgs)
	}

	db := s.OpScope.Tx(ctx)

	if s.ReturningID && s.Dialect == sqlxquery.DialectPostgres {
		sqlStr = db.Rebind(sqlStr + " RETURNING " + s.PKColumn)
		var id ID
		if err := db.QueryRowxContext(ctx, sqlStr, allArgs...).Scan(&id); err != nil {
			return *new(ID), err
		}
		return id, nil
	}

	sqlStr = db.Rebind(sqlStr)
	res, err := db.ExecContext(ctx, sqlStr, allArgs...)
	if err != nil {
		return *new(ID), err
	}

	if isZeroID {
		if lastID, err2 := res.LastInsertId(); err2 == nil && lastID != 0 {
			setPKField(&dto, s.PKColumn, lastID)
		}
	}

	return dto.GetID(), nil
}

// ---- helpers ----------------------------------------------------------------

func (s *Store[Entity, DTO, ID]) buildSelectSQL(result sqlxquery.Result, forceLimit int) string {
	cols := "*"
	if len(result.Cols) > 0 {
		cols = strings.Join(result.Cols, ", ")
	}

	var sb strings.Builder
	if result.Hint != "" {
		sb.WriteString(result.Hint)
		sb.WriteRune(' ')
	}
	sb.WriteString(fmt.Sprintf("SELECT %s FROM %s", cols, s.Table))

	if result.Where != "" {
		sb.WriteString(" WHERE " + result.Where)
	}
	if result.GroupBy != "" {
		sb.WriteString(" GROUP BY " + result.GroupBy)
	}
	if result.Having != "" {
		sb.WriteString(" HAVING " + result.Having)
	}
	if result.OrderBy != "" {
		sb.WriteString(" ORDER BY " + result.OrderBy)
	}

	limit := forceLimit
	if limit == 0 {
		limit = result.Limit
	}
	if limit > 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", limit))
		// Only emit OFFSET when not in the forced-single-row Get path.
		if forceLimit == 0 && result.Offset > 0 {
			sb.WriteString(fmt.Sprintf(" OFFSET %d", result.Offset))
		}
	}

	if result.Suffix != "" {
		sb.WriteRune(' ')
		sb.WriteString(result.Suffix)
	}

	return sb.String()
}

// buildWhere returns the WHERE clause and args for Update / PartialUpdate.
// When params are provided they take precedence; otherwise WHERE pk = id.
func (s *Store[Entity, DTO, ID]) buildWhere(id ID, params []query.Param) (whereStr string, whereArgs []any) {
	if len(params) > 0 {
		r := s.QueryBuilder.Build(query.NewParams(params...))
		return r.Where, r.Args
	}
	return s.PKColumn + " = ?", []any{id}
}

// checkPreload returns ErrPreloadNotSupported if any param is a PreloadParam.
func checkPreload(params []query.Param) error {
	for _, p := range params {
		if p.ParamType() == query.TypePreload {
			return ErrPreloadNotSupported
		}
	}
	return nil
}

// ---- dialect-specific Upsert builders --------------------------------------

func buildMySQLUpsert(table string, cols, phs []string, oc store.OnConflict, pkCol string, extraArgs *[]any) string {
	if oc.DoNothing {
		return fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES (%s)",
			table, strings.Join(cols, ", "), strings.Join(phs, ", "))
	}

	base := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table, strings.Join(cols, ", "), strings.Join(phs, ", "))

	updateClause := buildMySQLUpdateClause(cols, oc, pkCol, extraArgs)
	if updateClause == "" {
		return base
	}
	return base + " ON DUPLICATE KEY UPDATE " + updateClause
}

func buildMySQLUpdateClause(cols []string, oc store.OnConflict, pkCol string, extraArgs *[]any) string {
	if len(oc.Updates) > 0 {
		parts := make([]string, 0, len(oc.Updates))
		for col, val := range oc.Updates {
			parts = append(parts, col+" = ?")
			*extraArgs = append(*extraArgs, val)
		}
		return strings.Join(parts, ", ")
	}
	if len(oc.UpdateColumns) > 0 {
		parts := make([]string, len(oc.UpdateColumns))
		for i, col := range oc.UpdateColumns {
			parts[i] = col + " = VALUES(" + col + ")"
		}
		return strings.Join(parts, ", ")
	}
	if oc.UpdateAll {
		var parts []string
		for _, col := range cols {
			if col == pkCol {
				continue
			}
			parts = append(parts, col+" = VALUES("+col+")")
		}
		return strings.Join(parts, ", ")
	}
	return ""
}

func buildPostgresUpsert(table string, cols, phs []string, oc store.OnConflict, pkCol string, extraArgs *[]any) string {
	base := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table, strings.Join(cols, ", "), strings.Join(phs, ", "))

	if oc.DoNothing {
		return base + " ON CONFLICT DO NOTHING"
	}

	var conflictClause string
	switch {
	case oc.OnConstraint != "":
		conflictClause = "ON CONFLICT ON CONSTRAINT " + oc.OnConstraint
	case len(oc.Columns) > 0:
		conflictClause = "ON CONFLICT (" + strings.Join(oc.Columns, ", ") + ")"
	default:
		conflictClause = "ON CONFLICT (" + pkCol + ")"
	}

	updateClause := buildExcludedUpdateClause(cols, oc, pkCol, extraArgs)
	if updateClause == "" {
		return base + " " + conflictClause + " DO NOTHING"
	}
	return base + " " + conflictClause + " DO UPDATE SET " + updateClause
}

func buildSQLiteUpsert(table string, cols, phs []string, oc store.OnConflict, pkCol string, extraArgs *[]any) string {
	if oc.DoNothing {
		return fmt.Sprintf("INSERT OR IGNORE INTO %s (%s) VALUES (%s)",
			table, strings.Join(cols, ", "), strings.Join(phs, ", "))
	}
	// INSERT OR REPLACE only when no specific conflict columns are given
	if oc.UpdateAll && len(oc.Columns) == 0 && oc.OnConstraint == "" {
		return fmt.Sprintf("INSERT OR REPLACE INTO %s (%s) VALUES (%s)",
			table, strings.Join(cols, ", "), strings.Join(phs, ", "))
	}

	base := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table, strings.Join(cols, ", "), strings.Join(phs, ", "))

	conflictCols := oc.Columns
	if len(conflictCols) == 0 {
		conflictCols = []string{pkCol}
	}
	conflictClause := "ON CONFLICT (" + strings.Join(conflictCols, ", ") + ")"

	updateClause := buildExcludedUpdateClause(cols, oc, pkCol, extraArgs)
	if updateClause == "" {
		return base + " " + conflictClause + " DO NOTHING"
	}
	return base + " " + conflictClause + " DO UPDATE SET " + updateClause
}

func buildExcludedUpdateClause(cols []string, oc store.OnConflict, pkCol string, extraArgs *[]any) string {
	if len(oc.Updates) > 0 {
		parts := make([]string, 0, len(oc.Updates))
		for col, val := range oc.Updates {
			parts = append(parts, col+" = ?")
			*extraArgs = append(*extraArgs, val)
		}
		return strings.Join(parts, ", ")
	}
	if len(oc.UpdateColumns) > 0 {
		parts := make([]string, len(oc.UpdateColumns))
		for i, col := range oc.UpdateColumns {
			parts[i] = col + " = EXCLUDED." + col
		}
		return strings.Join(parts, ", ")
	}
	if oc.UpdateAll {
		var parts []string
		for _, col := range cols {
			if col == pkCol {
				continue
			}
			parts = append(parts, col+" = EXCLUDED."+col)
		}
		return strings.Join(parts, ", ")
	}
	return ""
}
