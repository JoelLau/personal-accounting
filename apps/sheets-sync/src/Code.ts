// ---------------------------------------------------------------------------
// personal-accounting → Google Sheets sync
//
// Setup:
//   1. Set these Script Properties (Extensions > Apps Script > Project Settings):
//      - DB_URL      e.g. jdbc:postgresql://your-host:5432/your-db
//      - DB_USER
//      - DB_PASSWORD
//   2. Deploy a time-driven trigger on `syncAll` (e.g. every hour).
// ---------------------------------------------------------------------------

const MICRO_SGD = 1_000_000;

function getConn(): GoogleAppsScript.JDBC.JdbcConnection {
  const props = PropertiesService.getScriptProperties();
  return Jdbc.getConnection(
    props.getProperty("DB_URL")!,
    props.getProperty("DB_USER")!,
    props.getProperty("DB_PASSWORD")!,
  );
}

function microSgdToSgd(micro: number): number {
  return micro / MICRO_SGD;
}

// ---------------------------------------------------------------------------
// Sheet helpers
// ---------------------------------------------------------------------------

function getOrCreateSheet(
  ss: GoogleAppsScript.Spreadsheet.Spreadsheet,
  name: string,
): GoogleAppsScript.Spreadsheet.Sheet {
  return ss.getSheetByName(name) ?? ss.insertSheet(name);
}

function writeSheet(
  ss: GoogleAppsScript.Spreadsheet.Spreadsheet,
  name: string,
  headers: string[],
  rows: unknown[][],
): void {
  const sheet = getOrCreateSheet(ss, name);
  sheet.clearContents();

  const all = [headers, ...rows];
  sheet.getRange(1, 1, all.length, headers.length).setValues(all);

  // Bold the header row and freeze it
  sheet.getRange(1, 1, 1, headers.length).setFontWeight("bold");
  sheet.setFrozenRows(1);

  // Mark the sheet as read-only for editors (owner retains full access)
  sheet.protect().setDescription(`${name} — managed by sheets-sync`).setWarningOnly(true);
}

// ---------------------------------------------------------------------------
// Sync: ledger_accounts
// ---------------------------------------------------------------------------

function syncLedgerAccounts(
  conn: GoogleAppsScript.JDBC.JdbcConnection,
  ss: GoogleAppsScript.Spreadsheet.Spreadsheet,
): void {
  const sql = `
    SELECT
      a.id,
      a.name,
      a.description,
      p.name AS parent_name
    FROM ledger_accounts a
    LEFT JOIN ledger_accounts p ON p.id = a.parent_id
    ORDER BY a.id
  `;
  const rs = conn.createStatement().executeQuery(sql);

  const rows: unknown[][] = [];
  while (rs.next()) {
    rows.push([
      rs.getLong(1),    // id
      rs.getString(2),  // name
      rs.getString(3),  // description
      rs.getString(4),  // parent_name
    ]);
  }
  rs.close();

  writeSheet(ss, "Accounts", ["ID", "Name", "Description", "Parent Account"], rows);
}

// ---------------------------------------------------------------------------
// Sync: postings (journal entry headers)
// ---------------------------------------------------------------------------

function syncPostings(
  conn: GoogleAppsScript.JDBC.JdbcConnection,
  ss: GoogleAppsScript.Spreadsheet.Spreadsheet,
): void {
  const sql = `
    SELECT
      p.id,
      p.transacted_at AT TIME ZONE 'Asia/Singapore' AS transacted_at_sgt,
      p.description,
      p.system_notes,
      p.source_hash
    FROM postings p
    ORDER BY p.transacted_at DESC
  `;
  const rs = conn.createStatement().executeQuery(sql);

  const rows: unknown[][] = [];
  while (rs.next()) {
    rows.push([
      rs.getLong(1),
      rs.getString(2),  // timestamp string
      rs.getString(3),
      rs.getString(4),
      rs.getString(5),
    ]);
  }
  rs.close();

  writeSheet(
    ss,
    "Postings",
    ["ID", "Transacted At (SGT)", "Description", "System Notes", "Source Hash"],
    rows,
  );
}

// ---------------------------------------------------------------------------
// Sync: entries (line items, with account name and SGD amounts)
// ---------------------------------------------------------------------------

function syncEntries(
  conn: GoogleAppsScript.JDBC.JdbcConnection,
  ss: GoogleAppsScript.Spreadsheet.Spreadsheet,
): void {
  const sql = `
    SELECT
      e.id,
      e.postings_id,
      p.transacted_at AT TIME ZONE 'Asia/Singapore' AS transacted_at_sgt,
      a.name                          AS account,
      e.description,
      e.debit_microsgd  / 1000000.0  AS debit_sgd,
      e.credit_microsgd / 1000000.0  AS credit_sgd,
      e.system_notes
    FROM entries e
    JOIN postings        p ON p.id = e.postings_id
    JOIN ledger_accounts a ON a.id = e.ledger_accounts_id
    ORDER BY p.transacted_at DESC, e.id
  `;
  const rs = conn.createStatement().executeQuery(sql);

  const rows: unknown[][] = [];
  while (rs.next()) {
    rows.push([
      rs.getLong(1),    // entry id
      rs.getLong(2),    // posting id
      rs.getString(3),  // transacted_at
      rs.getString(4),  // account name
      rs.getString(5),  // description
      rs.getDouble(6),  // debit SGD
      rs.getDouble(7),  // credit SGD
      rs.getString(8),  // system_notes
    ]);
  }
  rs.close();

  writeSheet(
    ss,
    "Entries",
    ["Entry ID", "Posting ID", "Transacted At (SGT)", "Account", "Description", "Debit (SGD)", "Credit (SGD)", "System Notes"],
    rows,
  );
}

// ---------------------------------------------------------------------------
// Entry point — set a time-driven trigger on this function
// ---------------------------------------------------------------------------

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function syncAll(): void {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const conn = getConn();

  try {
    syncLedgerAccounts(conn, ss);
    syncPostings(conn, ss);
    syncEntries(conn, ss);
    Logger.log("Sync complete: %s", new Date().toISOString());
  } finally {
    conn.close();
  }
}
