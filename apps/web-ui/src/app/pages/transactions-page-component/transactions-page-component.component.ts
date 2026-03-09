import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import {
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { Store } from '@ngxs/store';
import { BehaviorSubject, combineLatest, map, startWith } from 'rxjs';
import { Entry, Posting } from '../../services/accounting-api-service/models';
import { UpdateEntry } from '../../store/ledger.actions';
import { LedgerState } from '../../store/ledger.state';

@Component({
  selector: 'app-transactions-page-component',
  imports: [CommonModule, FormsModule, ReactiveFormsModule],
  templateUrl: './transactions-page-component.component.html',
})
export class TransactionsPageComponentComponent {
  private readonly state = inject(Store);

  private _rowStates = new BehaviorSubject<{ [posting_id: string]: RowState }>(
    {}
  );
  private rowState$ = this._rowStates.asObservable();

  accountsById$ = this.state.select(LedgerState.getAccounts);

  accounts$ = this.accountsById$.pipe(
    map((accountsById) => {
      return Object.values(accountsById);
    })
  );

  formGroup = new FormGroup({
    transactionsType: new FormControl('all', { nonNullable: true }),
  });

  rows$ = combineLatest({
    accounts: this.accountsById$,
    rowState: this.rowState$,
    transactions: this.state.select(LedgerState.getTransactions),
    tab: this.formGroup.controls.transactionsType.valueChanges.pipe(
      startWith('all')
    ),
  }).pipe(
    map(({ accounts, rowState, transactions, tab }): TableRow[] => {
      return Object.values(transactions).reduce((prev, curr) => {
        const entries = Object.values(curr.entries).filter((entry) => {
          if (tab == 'all') {
            return true;
          }
          return accounts[entry.ledger_accounts_id].qualified_name
            .toLocaleLowerCase()
            .startsWith(`${tab}:`);
        });

        return [
          ...prev,
          {
            ...curr,
            isExpanded: (rowState[curr.id] ?? { isExpanded: false }).isExpanded,
            entries: entries,
            debit: Object.values(entries).reduce((prev, curr) => {
              return prev + parseFloat(curr.debit_amount);
            }, 0),
            credit: Object.values(entries).reduce((prev, curr) => {
              return prev + parseFloat(curr.credit_amount);
            }, 0),
            total: Object.values(entries).reduce((prev, curr) => {
              return (
                prev +
                parseFloat(curr.debit_amount) -
                parseFloat(curr.credit_amount)
              );
            }, 0),
          },
        ];
      }, [] as TableRow[]);
    })
  );

  expandAll() {
    this._rowStates.next(
      Object.values(this.state.selectSnapshot(LedgerState.getPostings)).reduce(
        (prev, curr) => {
          return {
            ...prev,
            [curr.id]: {
              ...prev[curr.id],
              isExpanded: true,
            },
          };
        },
        {} as Record<string, RowState>
      )
    );
  }

  foldAll() {
    this._rowStates.next(
      Object.values(this.state.selectSnapshot(LedgerState.getPostings)).reduce(
        (prev, curr) => {
          return {
            ...prev,
            [curr.id]: {
              ...prev[curr.id],
              isExpanded: false,
            },
          };
        },
        {} as Record<string, RowState>
      )
    );
  }

  trackByIdField(_: number, data: { id: string }): string {
    return data.id;
  }

  toggleRow(posting_id: string) {
    console.log(posting_id);

    this._rowStates.next({
      ...this._rowStates.value,
      [posting_id]: {
        ...this._rowStates.value[posting_id],
        isExpanded: !(
          this._rowStates.value[posting_id] ?? { isExpanded: false }
        ).isExpanded,
      },
    });
  }

  onEntryChange(entryId: string, field: keyof Entry, newValue: string) {
    const entry: Entry = {
      ...this.state.selectSnapshot(LedgerState.getEntries)[entryId],
      [field]: `${newValue}`, // WARN: temporary workaround
    };

    this.state.dispatch(new UpdateEntry(entryId, entry));
  }
}

interface RowState {
  isExpanded: boolean;
}

interface TableRow extends Posting, RowState {
  entries: Entry[];
  debit: number;
  credit: number;
  total: number;
}
