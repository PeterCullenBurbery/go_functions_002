// oracle_database_system_management_functions.go
// Package oracle_database_system_management_functions: helpers for CDB/PDB lifecycle and verification.
package oracle_database_system_management_functions

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/godror/godror"
	"github.com/PeterCullenBurbery/go_functions_002/v5/date_time_functions"
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
	root_dir, err := Get_root_datafile_directory(ctx, db)
	if err != nil {
		return err
	}
	expected := normalize_for_compare(root_dir + "PDBSEED\\")
	actual_dir, err := Get_pdbseed_datafile_directory(ctx, db)
	if err != nil {
		return err
	}
	actual := normalize_for_compare(actual_dir)
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
	exists_sql := fmt.Sprintf(`SELECT COUNT(*) FROM DBA_PDBS WHERE PDB_NAME = UPPER('%s')`, pdb_name)
	if err := db.QueryRowContext(ctx, exists_sql).Scan(&cnt); err != nil {
		return "", fmt.Errorf("checking if PDB exists: %w", err)
	}
	if cnt > 0 {
		return "", fmt.Errorf("PDB %s already exists", pdb_name)
	}

	root_dir, err := Get_root_datafile_directory(ctx, db)
	if err != nil {
		return "", err
	}
	seed_dir := normalize_windows_dir(root_dir + "PDBSEED\\")
	dest_dir := normalize_windows_dir(root_dir + pdb_name + `\`)

	create_sql := fmt.Sprintf(
		"CREATE PLUGGABLE DATABASE %s ADMIN USER %s IDENTIFIED BY %s FILE_NAME_CONVERT = ('%s', '%s')",
		pdb_name,
		admin_user,
		admin_password,
		escape_single_quotes(seed_dir),
		escape_single_quotes(dest_dir),
	)

	log.Println("▶ Executing:", create_sql)
	if _, err := db.ExecContext(ctx, create_sql); err != nil {
		return "", fmt.Errorf("CREATE PLUGGABLE DATABASE failed: %w", err)
	}
	return dest_dir, nil
}

// Open_pluggable_database_read_write opens the given PDB READ WRITE.
func Open_pluggable_database_read_write(ctx context.Context, db *sql.DB, pdb_name string) error {
	sql_text := fmt.Sprintf("ALTER PLUGGABLE DATABASE %s OPEN READ WRITE", pdb_name)
	_, err := db.ExecContext(ctx, sql_text)
	return err
}

// Save_pluggable_database_state saves the auto-open state for the given PDB.
func Save_pluggable_database_state(ctx context.Context, db *sql.DB, pdb_name string) error {
	sql_text := fmt.Sprintf("ALTER PLUGGABLE DATABASE %s SAVE STATE", pdb_name)
	_, err := db.ExecContext(ctx, sql_text)
	return err
}

// Get_pdb_status returns the OPEN_MODE from V$PDBS for the given PDB name.
func Get_pdb_status(ctx context.Context, db *sql.DB, pdb_name string) (string, error) {
	sql_text := fmt.Sprintf("SELECT OPEN_MODE FROM V$PDBS WHERE NAME = UPPER('%s')", pdb_name)
	var open_mode string
	if err := db.QueryRowContext(ctx, sql_text).Scan(&open_mode); err != nil {
		return "", err
	}
	return open_mode, nil
}

// Get_saved_state_info returns (STATE, RESTRICTED) from DBA_PDB_SAVED_STATES for the PDB (by CON_NAME).
// If the row is not present, it returns ("","") and nil error.
func Get_saved_state_info(ctx context.Context, db *sql.DB, pdb_name string) (string, string, error) {
	sql_text := fmt.Sprintf(`
SELECT state, restricted
FROM   dba_pdb_saved_states
WHERE  con_name = UPPER('%s')`, pdb_name)
	var state, restricted string
	err := db.QueryRowContext(ctx, sql_text).Scan(&state, &restricted)
	if err == sql.ErrNoRows {
		return "", "", nil
	}
	return state, restricted, err
}

// Close_pluggable_database_immediate closes the PDB immediately.
// If instances_all is true, issues "INSTANCES=ALL" (for RAC).
func Close_pluggable_database_immediate(ctx context.Context, db *sql.DB, pdb_name string, instances_all bool) error {
	sql_text := fmt.Sprintf("ALTER PLUGGABLE DATABASE %s CLOSE IMMEDIATE", pdb_name)
	if instances_all {
		sql_text += " INSTANCES=ALL"
	}
	_, err := db.ExecContext(ctx, sql_text)
	return err
}

// Discard_pluggable_database_state removes any saved state for the PDB,
// preventing auto-open on next CDB restart.
func Discard_pluggable_database_state(ctx context.Context, db *sql.DB, pdb_name string) error {
	sql_text := fmt.Sprintf("ALTER PLUGGABLE DATABASE %s DISCARD STATE", pdb_name)
	_, err := db.ExecContext(ctx, sql_text)
	return err
}

// Drop_pluggable_database_including_datafiles drops the PDB and removes its datafiles.
// The PDB must be closed on all instances (in RAC) before this succeeds.
func Drop_pluggable_database_including_datafiles(ctx context.Context, db *sql.DB, pdb_name string) error {
	sql_text := fmt.Sprintf("DROP PLUGGABLE DATABASE %s INCLUDING DATAFILES", pdb_name)
	_, err := db.ExecContext(ctx, sql_text)
	return err
}

// Verify_pluggable_database_dropped returns true if DBA_PDBS no longer has the PDB.
func Verify_pluggable_database_dropped(ctx context.Context, db *sql.DB, pdb_name string) (bool, error) {
	var cnt int
	sql_text := fmt.Sprintf("SELECT COUNT(*) FROM DBA_PDBS WHERE PDB_NAME = UPPER('%s')", pdb_name)
	if err := db.QueryRowContext(ctx, sql_text).Scan(&cnt); err != nil {
		return false, err
	}
	return cnt == 0, nil
}

// user_session represents a single USER session in a PDB.
type user_session struct {
	inst_id int
	sid     int
	serial  int
	username string
	status   string
	machine  string
	program  string
	module   string
}

// Get_user_sessions returns all USER sessions attached to the given PDB.
// Note: the returned struct and its fields are unexported (snake_case). If you
// want callers in other packages to access fields directly, we can provide an
// exported DTO or accessor helpers.
func Get_user_sessions(ctx context.Context, db *sql.DB, pdb_name string) ([]user_session, error) {
	const q = `
SELECT s.inst_id,
       s.sid,
       s.serial#,
       NVL(s.username, ' '),
       NVL(s.status,   ' '),
       NVL(s.machine,  ' '),
       NVL(s.program,  ' '),
       NVL(s.module,   ' ')
FROM   gv$session s
WHERE  s.type = 'USER'
AND    s.con_id = (SELECT con_id FROM v$pdbs WHERE name = UPPER(:1))
ORDER  BY s.inst_id, s.sid`

	rows, err := db.QueryContext(ctx, q, pdb_name)
	if err != nil {
		return nil, fmt.Errorf("querying user sessions for %s: %w", pdb_name, err)
	}
	defer rows.Close()

	var sessions []user_session
	for rows.Next() {
		var us user_session
		if err := rows.Scan(
			&us.inst_id,
			&us.sid,
			&us.serial,
			&us.username,
			&us.status,
			&us.machine,
			&us.program,
			&us.module,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, us)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

// Kill_user_sessions_in_pdb terminates all USER sessions attached to the given PDB (single pass).
func Kill_user_sessions_in_pdb(ctx context.Context, db *sql.DB, pdb_name string) error {
    sessions, err := Get_user_sessions(ctx, db, pdb_name)
    if err != nil {
        return fmt.Errorf("querying user sessions for %s: %w", pdb_name, err)
    }
    if len(sessions) == 0 {
        log.Printf("✓ No USER sessions found in %s.", pdb_name)
        return nil
    }

    for _, us := range sessions {
        kill := fmt.Sprintf("ALTER SYSTEM KILL SESSION '%d,%d' IMMEDIATE", us.sid, us.serial)
        // RAC note: to target a specific instance use 'sid,serial#,@inst_id'
        if _, err := db.ExecContext(ctx, kill); err != nil {
            log.Printf("⚠️ Failed to kill session sid=%d serial=%d (inst=%d): %v",
                us.sid, us.serial, us.inst_id, err)
        }
    }
    return nil
}


// Kill_user_sessions_in_pdb_until_gone keeps killing USER sessions in the PDB until none remain
// or until max_attempts (default: 100) is reached. It rechecks with Get_user_sessions between attempts.
// wait_between_attempts controls the sleep duration between attempts.
func Kill_user_sessions_in_pdb_until_gone(
    ctx context.Context,
    db *sql.DB,
    pdb_name string,
    max_attempts int,
    wait_between_attempts time.Duration,
) error {

    if max_attempts <= 0 {
        max_attempts = 100
    }
    if wait_between_attempts <= 0 {
        wait_between_attempts = 300 * time.Millisecond
    }

    for attempt := 1; attempt <= max_attempts; attempt++ {
        sessions, err := Get_user_sessions(ctx, db, pdb_name)
        if err != nil {
            return fmt.Errorf("attempt %d: querying user sessions for %s: %w", attempt, pdb_name, err)
        }

        if len(sessions) == 0 {
            log.Printf("✓ No USER sessions remain in %s (after %d attempt(s)).", pdb_name, attempt-1)
            return nil
        }

        log.Printf("🔎 Attempt %d/%d: %d USER session(s) found in %s. Killing...",
            attempt, max_attempts, len(sessions), pdb_name)

        // Kill all visible sessions this round
        for _, us := range sessions {
            kill := fmt.Sprintf("ALTER SYSTEM KILL SESSION '%d,%d' IMMEDIATE", us.sid, us.serial)
            if _, err := db.ExecContext(ctx, kill); err != nil {
                log.Printf("⚠️ Failed to kill session sid=%d serial=%d (inst=%d): %v",
                    us.sid, us.serial, us.inst_id, err)
            }
        }

        // Short pause before re-checking
        time.Sleep(wait_between_attempts)
    }

    // Final check and error
    remaining, err := Get_user_sessions(ctx, db, pdb_name)
    if err != nil {
        return fmt.Errorf("final check: querying user sessions for %s: %w", pdb_name, err)
    }
    if len(remaining) == 0 {
        log.Printf("✓ No USER sessions remain in %s (after %d attempts).", pdb_name, max_attempts)
        return nil
    }
    return fmt.Errorf("after %d attempts, %d USER session(s) still remain in %s",
        max_attempts, len(remaining), pdb_name)
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

// Create_open_save_state_pdb_from_seed generates a PDB name using
// date_time_functions.Generate_pdb_name_from_timestamp(), creates the PDB from PDB$SEED,
// opens it READ WRITE, saves state, verifies, and returns (pdb_name, dest_dir).
func Create_open_save_state_pdb_from_seed(
	ctx context.Context,
	db *sql.DB,
	admin_user string,
	admin_password string,
) (string, string, error) {

	// Ensure we're in CDB$ROOT and PDBSEED path is sane
	if err := Ensure_connected_to_cdb_root(ctx, db); err != nil {
		return "", "", err
	}
	if err := Verify_pdbseed_directory_matches_expected(ctx, db); err != nil {
		return "", "", err
	}

	// Generate name (pdb_YYYY_MMM_DDD_HHH_MMM_SSS)
	pdb_name, err := date_time_functions.Generate_pdb_name_from_timestamp()
	if err != nil {
		return "", "", fmt.Errorf("generate pdb name: %w", err)
	}

	// Create from seed (returns destination directory)
	dest_dir, err := Create_pluggable_database_from_seed(ctx, db, pdb_name, admin_user, admin_password)
	if err != nil {
		return "", "", err
	}

	// Open READ WRITE
	if err := Open_pluggable_database_read_write(ctx, db, pdb_name); err != nil {
		return "", "", fmt.Errorf("open pdb read write: %w", err)
	}

	// Verify open mode
	if open_mode, err := Get_pdb_status(ctx, db, pdb_name); err == nil {
		log.Println("🔎 PDB open mode:", open_mode)
	} else {
		log.Println("ℹ️ Could not read PDB open mode:", err)
	}

	// Save state
	if err := Save_pluggable_database_state(ctx, db, pdb_name); err != nil {
		return "", "", fmt.Errorf("save pdb state: %w", err)
	}

	// Verify saved state info (DBA_PDB_SAVED_STATES)
	if state, restricted, err := Get_saved_state_info(ctx, db, pdb_name); err == nil {
		if state != "" {
			log.Printf("💾 Saved state recorded: STATE=%s, RESTRICTED=%s\n", state, restricted)
		} else {
			log.Println("ℹ️ No saved state record found for PDB (view present but row missing).")
		}
	} else {
		log.Println("ℹ️ Could not read DBA_PDB_SAVED_STATES (view may be unavailable):", err)
	}

	return pdb_name, dest_dir, nil
}

// Teardown_drop_pdb closes, discards state, and drops a PDB INCLUDING DATAFILES.
// If kill_sessions is true, it will attempt to terminate USER sessions in that PDB first
// (retrying until none remain, up to 100 attempts, with 300ms between attempts).
// If instances_all is true, CLOSE IMMEDIATE will be issued with INSTANCES=ALL (for RAC).
// The function also verifies the drop via DBA_PDBS and inspects saved state after DISCARD STATE.
func Teardown_drop_pdb(
	ctx context.Context,
	db *sql.DB,
	pdb_name string,
	instances_all bool,
	kill_sessions bool,
) error {

	// Must be in CDB$ROOT to manage PDBs
	if err := Ensure_connected_to_cdb_root(ctx, db); err != nil {
		return err
	}

	// Optionally kill active USER sessions to avoid close failures (retry until gone)
	if kill_sessions {
		if err := Kill_user_sessions_in_pdb_until_gone(ctx, db, pdb_name, 100, 300*time.Millisecond); err != nil {
			return fmt.Errorf("kill sessions in %s: %w", pdb_name, err)
		}
	}

	// Try to close; if it fails due to sessions, try killing and retry once
	if err := Close_pluggable_database_immediate(ctx, db, pdb_name, instances_all); err != nil {
		// one retry path: attempt to kill sessions and close again
		if !kill_sessions {
			_ = Kill_user_sessions_in_pdb_until_gone(ctx, db, pdb_name, 100, 300*time.Millisecond) // best effort
		}
		if err_second_close := Close_pluggable_database_immediate(ctx, db, pdb_name, instances_all); err_second_close != nil {
			return fmt.Errorf("close pdb %s: %w", pdb_name, err_second_close)
		}
	}

	// Discard auto-open state
	if err := Discard_pluggable_database_state(ctx, db, pdb_name); err != nil {
		// non-fatal; continue drop anyway
		log.Printf("⚠️ DISCARD STATE for %s returned: %v (continuing)", pdb_name, err)
	} else {
		// Inspect saved state after DISCARD STATE (view may be unavailable in some editions)
		state, restricted, err_state := Get_saved_state_info(ctx, db, pdb_name)
		if err_state != nil {
			log.Println("ℹ️ Could not read DBA_PDB_SAVED_STATES (view may be unavailable):", err_state)
		} else if state == "" {
			log.Printf("💾 Saved state for %s is absent after DISCARD STATE (as expected).", pdb_name)
		} else {
			log.Printf("⚠️ Saved state still present after DISCARD STATE for %s: STATE=%s, RESTRICTED=%s", pdb_name, state, restricted)
		}
	}

	// Drop INCLUDING DATAFILES
	if err := Drop_pluggable_database_including_datafiles(ctx, db, pdb_name); err != nil {
		return fmt.Errorf("drop pdb %s INCLUDING DATAFILES: %w", pdb_name, err)
	}

	// Verify gone
	gone, err := Verify_pluggable_database_dropped(ctx, db, pdb_name)
	if err != nil {
		return fmt.Errorf("verify drop for %s: %w", pdb_name, err)
	}
	if !gone {
		return fmt.Errorf("drop not confirmed: %s still appears in DBA_PDBS", pdb_name)
	}

	return nil
}