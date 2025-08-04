// oracle_database_system_management_functions.go
// Package oracle_database_system_management_functions: helpers for CDB/PDB lifecycle and verification.
package oracle_database_system_management_functions

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/godror/godror"
)

// -------------------------------
// Path helpers (CDB root & PDB$SEED)
// -------------------------------

// Get_root_datafile_directory returns the directory where CDB$ROOT's SYSTEM01.DBF resides,
// normalized to Windows backslash with a trailing backslash.
func Get_root_datafile_directory(ctx context.Context, db *sql.DB) (string, error) {
	const q = `
SELECT DISTINCT
       SUBSTR(name, 1, REGEXP_INSTR(name, 'SYSTEM01\.DBF', 1, 1, 0, 'i') - 1)
FROM   v$datafile
WHERE  REGEXP_LIKE(name, 'SYSTEM01\.DBF', 'i')
  AND  NOT REGEXP_LIKE(name, '[\\/]{1}PDB[^\\/]*', 'i')
`
	var dir string
	if err := db.QueryRowContext(ctx, q).Scan(&dir); err != nil {
		return "", fmt.Errorf("fetching CDB$ROOT datafile directory: %w", err)
	}
	return normalize_windows_dir(dir), nil
}

// Get_pdbseed_datafile_directory returns the directory where PDB$SEED's SYSTEM01.DBF resides,
// normalized to Windows backslash with a trailing backslash.
func Get_pdbseed_datafile_directory(ctx context.Context, db *sql.DB) (string, error) {
	const q = `
SELECT DISTINCT
       SUBSTR(name, 1, REGEXP_INSTR(name, 'SYSTEM01\.DBF', 1, 1, 0, 'i') - 1)
FROM   v$datafile
WHERE  REGEXP_LIKE(name, '[\\/]{1}PDBSEED[\\/]{1}SYSTEM01\.DBF', 'i')
`
	var dir string
	if err := db.QueryRowContext(ctx, q).Scan(&dir); err != nil {
		return "", fmt.Errorf("fetching PDB$SEED datafile directory: %w", err)
	}
	return normalize_windows_dir(dir), nil
}

// Verify_pdbseed_directory_matches_expected checks that root\PDBSEED\ equals actual PDB$SEED location.
func Verify_pdbseed_directory_matches_expected(ctx context.Context, db *sql.DB) error {
	rootDir, err := Get_root_datafile_directory(ctx, db)
	if err != nil {
		return err
	}
	expected := normalize_for_compare(rootDir + "PDBSEED\\")
	actualDir, err := Get_pdbseed_datafile_directory(ctx, db)
	if err != nil {
		return err
	}
	actual := normalize_for_compare(actualDir)
	if expected != actual {
		return fmt.Errorf("expected PDBSEED at %s but found %s", expected, actual)
	}
	return nil
}

// -------------------------------
// CDB container guard
// -------------------------------

// Ensure_connected_to_cdb_root returns error unless current container is CDB$ROOT.
func Ensure_connected_to_cdb_root(ctx context.Context, db *sql.DB) error {
	var con string
	if err := db.QueryRowContext(ctx, `SELECT SYS_CONTEXT('USERENV','CON_NAME') FROM dual`).Scan(&con); err != nil {
		return fmt.Errorf("cannot determine current container: %w", err)
	}
	if !strings.EqualFold(con, "CDB$ROOT") {
		return fmt.Errorf("not connected to CDB$ROOT (current: %s)", con)
	}
	return nil
}

// -------------------------------
// PDB lifecycle helpers
// -------------------------------

// Create_pluggable_database_from_seed creates a new PDB using FILE_NAME_CONVERT from PDB$SEED.
// It verifies the PDB$SEED path matches the expected path (root + PDBSEED\) before creating.
// admin_user/admin_password are used for the PDB admin account.
// Returns the computed destination directory (root\pdb_name\).
func Create_pluggable_database_from_seed(ctx context.Context, db *sql.DB, pdb_name, admin_user, admin_password string) (string, error) {
	if err := Ensure_connected_to_cdb_root(ctx, db); err != nil {
		return "", err
	}
	if err := Verify_pdbseed_directory_matches_expected(ctx, db); err != nil {
		return "", err
	}

	// check not exists
	var cnt int
	existsSQL := fmt.Sprintf(`SELECT COUNT(*) FROM DBA_PDBS WHERE PDB_NAME = UPPER('%s')`, pdb_name)
	if err := db.QueryRowContext(ctx, existsSQL).Scan(&cnt); err != nil {
		return "", fmt.Errorf("checking if PDB exists: %w", err)
	}
	if cnt > 0 {
		return "", fmt.Errorf("PDB %s already exists", pdb_name)
	}

	rootDir, err := Get_root_datafile_directory(ctx, db)
	if err != nil {
		return "", err
	}
	seedDir := normalize_windows_dir(rootDir + "PDBSEED\\")
	destDir := normalize_windows_dir(rootDir + pdb_name + `\`)

	createSQL := fmt.Sprintf(
		"CREATE PLUGGABLE DATABASE %s ADMIN USER %s IDENTIFIED BY %s FILE_NAME_CONVERT = ('%s', '%s')",
		pdb_name,
		admin_user,
		admin_password,
		escape_single_quotes(seedDir),
		escape_single_quotes(destDir),
	)

	log.Println("▶ Executing:", createSQL)
	if _, err := db.ExecContext(ctx, createSQL); err != nil {
		return "", fmt.Errorf("CREATE PLUGGABLE DATABASE failed: %w", err)
	}
	return destDir, nil
}

// Open_pluggable_database_read_write opens the given PDB READ WRITE.
func Open_pluggable_database_read_write(ctx context.Context, db *sql.DB, pdb_name string) error {
	sqlText := fmt.Sprintf("ALTER PLUGGABLE DATABASE %s OPEN READ WRITE", pdb_name)
	_, err := db.ExecContext(ctx, sqlText)
	return err
}

// Save_pluggable_database_state saves the auto-open state for the given PDB.
func Save_pluggable_database_state(ctx context.Context, db *sql.DB, pdb_name string) error {
	sqlText := fmt.Sprintf("ALTER PLUGGABLE DATABASE %s SAVE STATE", pdb_name)
	_, err := db.ExecContext(ctx, sqlText)
	return err
}

// Get_pdb_status returns the OPEN_MODE from V$PDBS for the given PDB name.
func Get_pdb_status(ctx context.Context, db *sql.DB, pdb_name string) (string, error) {
	sqlText := fmt.Sprintf("SELECT OPEN_MODE FROM V$PDBS WHERE NAME = UPPER('%s')", pdb_name)
	var openMode string
	if err := db.QueryRowContext(ctx, sqlText).Scan(&openMode); err != nil {
		return "", err
	}
	return openMode, nil
}

// Get_saved_state_info returns (STATE, RESTRICTED) from DBA_PDB_SAVED_STATES for the PDB (by CON_NAME).
// If the row is not present, it returns ("","") and nil error.
func Get_saved_state_info(ctx context.Context, db *sql.DB, pdb_name string) (string, string, error) {
	sqlText := fmt.Sprintf(`
SELECT state, restricted
FROM   dba_pdb_saved_states
WHERE  con_name = UPPER('%s')`, pdb_name)
	var state, restricted string
	err := db.QueryRowContext(ctx, sqlText).Scan(&state, &restricted)
	if err == sql.ErrNoRows {
		return "", "", nil
	}
	return state, restricted, err
}

// Close_pluggable_database_immediate closes the PDB immediately.
// If instances_all is true, issues "INSTANCES=ALL" (for RAC).
func Close_pluggable_database_immediate(ctx context.Context, db *sql.DB, pdb_name string, instances_all bool) error {
	sqlText := fmt.Sprintf("ALTER PLUGGABLE DATABASE %s CLOSE IMMEDIATE", pdb_name)
	if instances_all {
		sqlText += " INSTANCES=ALL"
	}
	_, err := db.ExecContext(ctx, sqlText)
	return err
}

// Discard_pluggable_database_state removes any saved state for the PDB,
// preventing auto-open on next CDB restart.
func Discard_pluggable_database_state(ctx context.Context, db *sql.DB, pdb_name string) error {
	sqlText := fmt.Sprintf("ALTER PLUGGABLE DATABASE %s DISCARD STATE", pdb_name)
	_, err := db.ExecContext(ctx, sqlText)
	return err
}

// Drop_pluggable_database_including_datafiles drops the PDB and removes its datafiles.
// The PDB must be closed on all instances (in RAC) before this succeeds.
func Drop_pluggable_database_including_datafiles(ctx context.Context, db *sql.DB, pdb_name string) error {
	sqlText := fmt.Sprintf("DROP PLUGGABLE DATABASE %s INCLUDING DATAFILES", pdb_name)
	_, err := db.ExecContext(ctx, sqlText)
	return err
}

// Verify_pluggable_database_dropped returns true if DBA_PDBS no longer has the PDB.
func Verify_pluggable_database_dropped(ctx context.Context, db *sql.DB, pdb_name string) (bool, error) {
	var cnt int
	sqlText := fmt.Sprintf("SELECT COUNT(*) FROM DBA_PDBS WHERE PDB_NAME = UPPER('%s')", pdb_name)
	if err := db.QueryRowContext(ctx, sqlText).Scan(&cnt); err != nil {
		return false, err
	}
	return cnt == 0, nil
}

// Kill_user_sessions_in_pdb terminates all USER sessions attached to the given PDB.
func Kill_user_sessions_in_pdb(ctx context.Context, db *sql.DB, pdb_name string) error {
	// gather sessions to kill
	const q = `
SELECT s.inst_id, s.sid, s.serial#
FROM   gv$session s
WHERE  s.type = 'USER'
AND    s.con_id = (SELECT con_id FROM v$pdbs WHERE name = UPPER(:1))
`
	rows, err := db.QueryContext(ctx, q, pdb_name)
	if err != nil {
		return fmt.Errorf("querying sessions: %w", err)
	}
	defer rows.Close()

	type sess struct {
		instID int
		sid    int
		serial int
	}
	var sessions []sess
	for rows.Next() {
		var r sess
		if err := rows.Scan(&r.instID, &r.sid, &r.serial); err != nil {
			return err
		}
		sessions = append(sessions, r)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// kill sessions
	for _, s := range sessions {
		kill := fmt.Sprintf("ALTER SYSTEM KILL SESSION '%d,%d' IMMEDIATE", s.sid, s.serial)
		// In RAC, optionally: "ALTER SYSTEM KILL SESSION 'sid,serial#,@inst_id' IMMEDIATE"
		// if you need to target specific instance. Uncomment and modify if required.
		if _, err := db.ExecContext(ctx, kill); err != nil {
			log.Printf("⚠️ Failed to kill session sid=%d serial=%d: %v", s.sid, s.serial, err)
		}
	}
	return nil
}

// -------------------------------
// Utility
// -------------------------------

func normalize_windows_dir(p string) string {
	p = strings.ReplaceAll(p, "/", `\`)
	if !strings.HasSuffix(p, `\`) {
		p += `\`
	}
	return p
}

func normalize_for_compare(p string) string {
	p = strings.ReplaceAll(p, "/", `\`)
	if !strings.HasSuffix(p, `\`) {
		p += `\`
	}
	return strings.ToUpper(p)
}

func escape_single_quotes(s string) string {
	return strings.ReplaceAll(s, `'`, `''`)
}
