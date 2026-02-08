import { inject, Injectable } from '@angular/core';
import { Action, Selector, State, StateContext } from '@ngxs/store';
import { forkJoin, tap } from 'rxjs';
import { AccountingService } from '../services/accounting-api-service/services';
import {
  Entry,
  LedgerAccount,
  Posting,
} from '../services/accounting-api-service/models';
import { FetchLedgerData } from './ledger.actions';

export interface LedgerStateModel {
  accounts: { [id: string]: LedgerAccount };
  postings: { [id: string]: Posting };
  entries: { [id: string]: Entry };
}

export interface LedgerAccountWithChildren extends LedgerAccount {
  children: LedgerAccountWithChildren[];
}

@State<LedgerStateModel>({
  name: 'ledger',
  defaults: {
    accounts: {},
    postings: {},
    entries: {},
  },
})
@Injectable()
export class LedgerState {
  private readonly accounting = inject(AccountingService);

  @Selector()
  static getAccounts({ accounts }: LedgerStateModel) {
    return accounts;
  }

  @Selector()
  static getPostings({ postings }: LedgerStateModel) {
    return postings;
  }

  @Selector()
  static GetEntries({ entries }: LedgerStateModel) {
    return entries;
  }

  @Action(FetchLedgerData)
  fetchAll(ctx: StateContext<LedgerStateModel>) {
    return forkJoin({
      accounts: this.accounting.apiV1AccountingLedgerAccountsGet(),
      postings: this.accounting.apiV1AccountingPostingsGet(),
      entries: this.accounting.apiV1AccountingEntriesGet(),
    }).pipe(
      tap((data) => {
        console.log(data);
      }),
      tap(({ accounts, postings, entries }) => {
        ctx.patchState({
          accounts: accounts.data?.reduce((prev, curr) => {
            prev[curr.id] = curr;
            return prev;
          }, {} as Record<string, LedgerAccount>),
          postings: postings.data?.reduce((prev, curr) => {
            prev[curr.id] = curr;
            return prev;
          }, {} as Record<string, Posting>),
          entries: entries.data?.reduce((prev, curr) => {
            prev[curr.id] = curr;
            return prev;
          }, {} as Record<string, Entry>),
        });
      })
    );
  }
}
