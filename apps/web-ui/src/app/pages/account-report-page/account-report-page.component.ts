import { Component, inject } from '@angular/core';
import { CommonModule, formatDate } from '@angular/common';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { combineLatest, map, startWith, Subject, switchMap } from 'rxjs';
import { AccountingService } from '../../services/accounting-api-service/services';
import { AccountBalanceNode } from '../../services/accounting-api-service/models';

const now = new Date();
const thisYear = now.getFullYear();
const thisMonth = now.getMonth();

const FIRST_DAY_OF_LAST_MONTH = new Date(thisYear, thisMonth - 1, 1);

@Component({
  selector: 'app-account-report-page',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './account-report-page.component.html',
})
export class AccountReportPageComponent {
  private readonly accounting = inject(AccountingService);

  formGroup = new FormGroup({
    transactionsType: new FormControl('expenses', { nonNullable: true }),
    month: new FormControl(
      formatDate(FIRST_DAY_OF_LAST_MONTH, 'yyyy-MM', 'en-US'),
      { nonNullable: true }
    ),
    hideParents: new FormControl(false, { nonNullable: true }),
  });

  refreshReports$ = new Subject<void>();

  report$ = combineLatest({
    monthFilter: this.formGroup.controls.month.valueChanges.pipe(
      startWith(this.formGroup.controls.month.value)
    ),
  }).pipe(
    switchMap(({ monthFilter }) => {
      const [year, month] = monthFilter.split('-').map(Number);
      const firstDay = new Date(year, month - 1, 1);
      const lastDay = new Date(year, month, 0);

      return this.accounting.getAccountBalances({
        start_date: formatDate(firstDay, 'yyyy-MM-dd', 'en-US'),
        end_date: formatDate(lastDay, 'yyyy-MM-dd', 'en-US'),
      });
    }),
    map(({ data }) => data)
  );

  accounts$ = combineLatest({
    report: this.report$,
    transactionType: this.formGroup.controls.transactionsType.valueChanges.pipe(
      startWith(this.formGroup.controls.transactionsType.value)
    ),
  }).pipe(
    map(({ report, transactionType }) => {
      return (report ?? []).filter((node) => {
        if (transactionType == 'all') {
          return true;
        }

        return node.name.toLocaleLowerCase().startsWith(transactionType);
      });
    })
  );

  trackBalanceNode(_: number, row: AccountBalanceNode) {
    return row.ledger_account_id;
  }

  onPreviousButtonClick() {
    const formControl = this.formGroup.controls.month;
    const [year, month] = formControl.value.split('-').map(Number);

    // NOTE: month is starts at index 0 (jan is 0)
    const date = new Date(year, month - 1, 1);
    date.setMonth(date.getMonth() - 1); // this line performs the actual subtraction

    this.formGroup.controls.month.setValue(
      formatDate(date, 'yyyy-MM', 'en-US')
    );
  }

  onNextButtonClick() {
    const formControl = this.formGroup.controls.month;
    const [year, month] = formControl.value.split('-').map(Number);

    // NOTE: month is starts at index 0 (jan is 0)
    const date = new Date(year, month - 1, 1);
    date.setMonth(date.getMonth() + 1); // this line performs the actual addition

    this.formGroup.controls.month.setValue(
      formatDate(date, 'yyyy-MM', 'en-US')
    );
  }
}
