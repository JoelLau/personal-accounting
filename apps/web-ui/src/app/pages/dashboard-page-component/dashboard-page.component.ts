import { CommonModule, formatCurrency } from '@angular/common';
import { Component, inject } from '@angular/core';
import { Store } from '@ngxs/store';
import { AgGridAngular } from 'ag-grid-angular';
import { GridOptions } from 'ag-grid-community';
import {
  BehaviorSubject,
  combineLatest,
  filter,
  map,
  Observable,
  take,
} from 'rxjs';
import { LedgerAccount } from '../../services/accounting-api-service/models';
import { AccountingService } from '../../services/accounting-api-service/services';
import { LedgerState } from '../../store/ledger.state';

export interface Row extends LedgerAccount {
  isFolded: boolean;
}

const CARET_RIGHT = '▸';
const CARET_DOWN = '▾';

@Component({
  selector: 'app-dashboard-page',
  imports: [CommonModule, AgGridAngular],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.scss',
})
export class DashboardPageComponent {
  private readonly accounting = inject(AccountingService);
  private readonly store = inject(Store);

  protected readonly ledgerAccounts$ = this.accounting
    .apiV1AccountingLedgerAccountsGet()
    .pipe(map((response) => response?.data ?? []));

  protected readonly postings$ = this.accounting
    .apiV1AccountingPostingsGet()
    .pipe(map((response) => response?.data ?? []));

  gridOptions: GridOptions<Row>;

  protected readonly foldState$ = new BehaviorSubject<Record<string, boolean>>(
    Object.keys(this.store.selectSnapshot(LedgerState.getAccounts)).reduce(
      (prev, key) => {
        prev[key] = true;
        return prev;
      },
      {} as Record<string, boolean>
    )
  );

  protected readonly accounts$ = this.store.select(LedgerState.getAccounts);

  protected readonly row$: Observable<Row[]> = combineLatest({
    foldState: this.foldState$,
    accounts: this.accounts$,
  }).pipe(
    map(({ foldState, accounts }) => {
      const allAccounts = Object.values(accounts);

      // 1. Helper to check if a node should be visible
      const isVisible = (account: LedgerAccount): boolean => {
        let current = account;
        while (current.parent_id) {
          if (foldState[current.parent_id]) return false; // Parent is folded
          current = accounts[current.parent_id];
          if (!current) break;
        }
        return true;
      };

      // 2. Helper to calculate depth
      const getDepth = (account: LedgerAccount): number => {
        let depth = 0;
        let current = account;
        while (current.parent_id && accounts[current.parent_id]) {
          depth++;
          current = accounts[current.parent_id];
        }
        return depth;
      };

      // 3. Recursive function to build the sorted list
      const buildTree = (parentId: string | null = null): Row[] => {
        return allAccounts
          .filter((a) => (a.parent_id || null) === parentId)
          .sort((a, b) =>
            (a.qualified_name ?? '').localeCompare(b.qualified_name ?? '')
          )
          .reduce((acc, curr) => {
            if (!isVisible(curr)) return acc;

            const row: Row & { depth: number } = {
              ...curr,
              isFolded: !!foldState[curr.id],
              depth: getDepth(curr),
            } as any;

            return [...acc, row, ...buildTree(curr.id)];
          }, [] as Row[]);
      };

      return buildTree(null);
    })
  );

  constructor() {
    this.accounts$
      .pipe(
        filter((accounts) => {
          console.log(accounts);
          return accounts && Object.keys(accounts).length > 0;
        }),
        take(1)
      )
      .subscribe((accounts) => {
        const initialState = Object.entries(accounts).reduce(
          (prev, [key, value]) => {
            prev[key] = !!value.parent_id;
            return prev;
          },
          {} as Record<string, boolean>
        );
        this.foldState$.next(initialState);
      });

    this.gridOptions = this.getGridOptions();
  }

  fetchLedgerAccountsByParentId(parentId?: string) {
    return this.store.select(LedgerState.getAccounts).pipe(
      map((data) => {
        return Object.values(data).filter((account) => {
          if (!parentId) {
            return true;
          }

          return account.parent_id == parentId;
        });
      })
    );
  }

  getGridOptions(): GridOptions<Row> {
    const foldState = this.foldState$;
    return {
      // ... existing options
      columnDefs: [
        {
          colId: 'fold',
          headerName: '',
          width: 50,
          valueGetter: (params) => params.data?.isFolded,
          valueFormatter: (params) => {
            return params.value ? CARET_RIGHT : CARET_DOWN;
          },
          onCellClicked: (event) => {
            const id = event.data?.id;
            if (!id) return;

            foldState.next({
              ...foldState.value,
              [id]: !event.data?.isFolded,
            });
          },
          // Apply dynamic padding to the first column or name column
          cellStyle: (params) => {
            const depth = (params.data as any)?.depth || 0;
            return {
              'padding-left': `${depth * 2}rem`,
              display: 'flex',
              'align-items': 'center',
              cursor: 'pointer',
            };
          },
        },
        {
          colId: 'name',
          field: 'qualified_name',
          headerName: 'Account Name',
          cellStyle: () => {
            return {
              'font-family': 'monospace',
            };
          },
        },
        {
          colId: 'description',
          field: 'description',
          headerName: 'Description',
        },
        {
          colId: 'total',
          headerName: 'Total',
          type: ['numericColumn'],
          valueGetter: function () {
            const randomNumber = Math.floor(Math.random() * 100_000) + 1;
            return formatCurrency(randomNumber, 'EN_US', '$', 'SGD', '0.2-2');
          },
        },
      ],
    };
  }
}
