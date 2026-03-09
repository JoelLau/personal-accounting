import { UpdateAccountingEntry$Params } from '../services/accounting-api-service/functions';
import { Entry } from '../services/accounting-api-service/models';

export class FetchLedgerData {
  static readonly type = '[App Startup] Fetch Ledger Data';
}

export class UpdateEntry {
  static readonly type = '[Entry] Edit';

  constructor(public entryId: string, public body: Entry) {}
}
