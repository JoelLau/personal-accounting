import { inject, Injectable } from '@angular/core';
import { Action, Selector, State, StateContext } from '@ngxs/store';
import { forkJoin, switchMap, tap } from 'rxjs';
import { AccountingService } from '../services/accounting-api-service/services';
import {
  Entry,
  LedgerAccount,
  Posting,
} from '../services/accounting-api-service/models';
import { FetchLedgerData, UpdateEntry } from './ledger.actions';

export interface LedgerStateModel {
  accounts: { [id: string]: LedgerAccount };
  postings: { [id: string]: Posting };
  entries: { [id: string]: Entry };
}

export interface LedgerAccountWithChildren extends LedgerAccount {
  children: LedgerAccountWithChildren[];
}

export interface Transaction extends Posting {
  entries: { [id: string]: Entry };
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
  static getEntries({ entries }: LedgerStateModel) {
    console.log(entries);
    return entries;
  }

  @Selector()
  static getTransactions({
    postings,
    entries,
  }: LedgerStateModel): Record<string, Transaction> {
    const entriesByPostingId = Object.values(entries).reduce((prev, curr) => {
      if (!prev[curr.postings_id]) {
        prev[curr.postings_id] = {};
      }
      prev[curr.postings_id][curr.id] = curr;
      return prev;
    }, {} as { [parentId: string]: Record<string, Entry> });

    return Object.entries(postings).reduce((prev, [postingID, posting]) => {
      prev[postingID] = {
        ...posting,
        entries: entriesByPostingId[postingID],
      };
      return prev;
    }, {} as { [id: string]: Transaction });
  }

  @Action(FetchLedgerData)
  fetchAll(ctx: StateContext<LedgerStateModel>) {
    return forkJoin({
      accounts: this.accounting.apiV1AccountingLedgerAccountsGet(),
      postings: this.accounting.apiV1AccountingPostingsGet(),
      entries: this.accounting.getEntries(),
    }).pipe(
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

  @Action(UpdateEntry)
  updateEntry(ctx: StateContext<LedgerStateModel>, param: UpdateEntry) {
    return this.accounting
      .updateEntry({
        entry_id: param.entryId,
        body: param.body,
      })
      .pipe(
        tap((response) => console.log(response)),
        switchMap(() => {
          return ctx.dispatch(FetchLedgerData);
        })
      );
  }
}
