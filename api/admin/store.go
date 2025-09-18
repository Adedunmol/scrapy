package admin

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const UniqueViolationCode = "23505"

type Store interface {
	CreateAdmin(ctx context.Context, body *CreateAdminBody) (Admin, error)
	FindAdminByEmail(ctx context.Context, email string) (Admin, error)
	ComparePasswords(password, candidatePassword string) bool
	GetAdmins(ctx context.Context) ([]Admin, error)
	CreateRole(ctx context.Context, body *CreateRoleBody) (Role, error)
	CreatePermission(ctx context.Context, body *CreatePermissionBody) (Permission, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetPermissions(ctx context.Context, roleID string) ([]Permission, error)
	BatchCreatePermissions(ctx context.Context, perms []CreatePermissionBody) error
	BatchCreateRoles(ctx context.Context, roles []CreateRoleBody) error
	GetRolePermissions(ctx context.Context, roleID string) ([]Permission, error)
}

type AdminStore struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewAdminStore(db *pgxpool.Pool, queryTimeout time.Duration) *AdminStore {

	return &AdminStore{db: db, queryTimeout: queryTimeout}
}

func (s *AdminStore) CreateAdmin(ctx context.Context, body *CreateAdminBody) (Admin, error) {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "INSERT INTO admins (email, first_name, last_name, password) VALUES (@email, @firstName, @lastName, @password) RETURNING id, email, first_name, last_name;"
	args := pgx.NamedArgs{
		"email":     body.Email,
		"firstName": body.FirstName,
		"lastName":  body.LastName,
		"password":  body.Password,
	}

	var admin Admin

	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&admin.ID, &admin.Email, &admin.FirstName, &admin.LastName)

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == UniqueViolationCode {
			return Admin{}, helpers.ErrConflict
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Admin{}, fmt.Errorf("error scanning row (create admin): %w", err)
	}

	return admin, nil
}

func (s *AdminStore) FindAdminByEmail(ctx context.Context, email string) (Admin, error) {

	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "SELECT id, first_name, last_name, email, password FROM admins WHERE email = @email;"
	args := pgx.NamedArgs{
		"email": email,
	}

	var admin Admin

	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&admin.ID, &admin.Email, &admin.FirstName, &admin.LastName, &admin.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Admin{}, helpers.ErrNotFound
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Admin{}, fmt.Errorf("error scanning row (find user by email): %w", err)
	}

	return admin, nil
}

func (s *AdminStore) ComparePasswords(password, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(candidatePassword))

	if err != nil {
		return false
	}
	return true
}

func (s *AdminStore) GetAdmins(ctx context.Context) ([]Admin, error) {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "SELECT id, first_name, last_name, email FROM admins;"

	var admins []Admin

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return admins, fmt.Errorf("error fetching admins: %v", err)
	}

	for rows.Next() {
		var admin Admin
		if err := rows.Scan(&admin.ID, &admin.FirstName, &admin.LastName, &admin.Email); err != nil {
			return nil, fmt.Errorf("error scanning admins: %v", err)
		}
		admins = append(admins, admin)
	}

	return admins, nil
}

func (s *AdminStore) CreateRole(ctx context.Context, body *CreateRoleBody) (Role, error) {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "INSERT INTO roles (name, description, slug, created_by) VALUES (@name, @description, @slug, @createdBy) RETURNING id, name, description, slug, created_by, created_at;"
	args := pgx.NamedArgs{
		"name":        body.Name,
		"description": body.Description,
		"slug":        body.Slug,
		"createdBy":   body.CreatedBy,
	}

	var role Role

	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&role.ID, &role.Name, &role.Description, &role.Slug, &role.CreatedBy)

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == UniqueViolationCode {
			return Role{}, helpers.ErrConflict
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Role{}, fmt.Errorf("error scanning row (create role): %w", err)
	}

	return role, nil
}

func (s *AdminStore) CreatePermission(ctx context.Context, body *CreatePermissionBody) (Permission, error) {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "INSERT INTO permissions (name, description, slug, created_by) VALUES (@name, @description, @slug, @createdBy) RETURNING id, name, description, slug, created_by, created_at;"
	args := pgx.NamedArgs{
		"name":        body.Name,
		"description": body.Description,
		"slug":        body.Slug,
		"createdBy":   body.CreatedBy,
	}

	var permission Permission

	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.Slug, &permission.CreatedBy)

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == UniqueViolationCode {
			return Permission{}, helpers.ErrConflict
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Permission{}, fmt.Errorf("error scanning row (create permission): %w", err)
	}

	return permission, nil
}

func (s *AdminStore) GetRoles(ctx context.Context) ([]Role, error) {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	var roles []Role

	query := "SELECT id, name, description, slug, created_by, created_at FROM roles;"

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var role Role

		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.Slug, &role.CreatedBy, &role.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning roles: %v", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (s *AdminStore) GetPermissions(ctx context.Context, roleID string) ([]Permission, error) {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	var permissions []Permission

	query := "SELECT id, name, description, slug, created_by, created_at FROM permissions;"

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var perm Permission

		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Description, &perm.Slug, &perm.CreatedBy, &perm.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning permissions: %v", err)
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (s *AdminStore) BatchCreatePermissions(ctx context.Context, perms []CreatePermissionBody) error {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := `
		INSERT INTO permissions (
			name,
			description,
			slug,
			created_by,
			created_at,
			updated_at
		) VALUES (
			@name,
			@description,
			@slug,
			@createdBy,
			NOW(),
			NOW()
		);`

	batch := &pgx.Batch{}

	for _, perm := range perms {
		args := pgx.NamedArgs{
			"name":        perm.Name,
			"description": perm.Description,
			"slug":        perm.Slug,
			"createdBy":   perm.CreatedBy,
		}
		batch.Queue(query, args)
	}

	results := s.db.SendBatch(ctx, batch)
	defer results.Close()

	for range perms {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("error while batch creating permissions: %w", err)
		}
	}

	return results.Close()
}

func (s *AdminStore) BatchCreateRoles(ctx context.Context, roles []CreateRoleBody) error {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := `
		INSERT INTO roles (
			name,
			description,
			slug,
			created_by,
			created_at,
			updated_at
		) VALUES (
			@name,
			@description,
			@slug,
			@createdBy,
			NOW(),
			NOW()
		);`

	batch := &pgx.Batch{}

	for _, role := range roles {
		args := pgx.NamedArgs{
			"name":        role.Name,
			"description": role.Description,
			"slug":        role.Slug,
			"createdBy":   role.CreatedBy,
		}
		batch.Queue(query, args)
	}

	results := s.db.SendBatch(ctx, batch)
	defer results.Close()

	for range roles {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("error while batch creating roles: %w", err)
		}
	}

	return results.Close()
}

func (s *AdminStore) GetRolePermissions(ctx context.Context, roleID string) ([]Permission, error) {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := `
			SELECT p.id, p.name, p.description, p.slug, p.created_by, p.created_at
			FROM permissions p
			JOIN roles_perms_table rp ON rp.perm_id = p.id
			WHERE rp.role_id = @roleID;
			`
	args := pgx.NamedArgs{
		"roleID": roleID,
	}

	rows, err := s.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	var permissions []Permission

	for rows.Next() {
		var perm Permission

		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Description, &perm.Slug, &perm.CreatedBy, &perm.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning role permissions: %v", err)
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (s *AdminStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, s.queryTimeout)
}
