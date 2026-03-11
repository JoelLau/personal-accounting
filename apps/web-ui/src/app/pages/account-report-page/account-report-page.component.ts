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

const ASSETS = 1;
const LIABILITIES = 2;
const INCOME = 3;
const EXPENSES = 4;
const EQUITY = 5;
const LIABILITIES_CREDITCARD = 2001;
const INCOME_UNCATEGORIZED = 3000;
const EXPENSES_UNCATEGORIZED = 4000;
const OBLIGATIONS = 4100;
const HOME = 4200;
const LIFESTYLE = 4300;
const FOOD = 4400;
const TRANSPORT = 4500;
const HEALTHnGROWTH = 4600;
const INSURANCE = 4101;
const TAX = 4102;
const PARENTS = 4103;
const MAINTENANCE = 4201;
const CLEANER = 4202;
const DOG = 4203;
const PHONE = 4204;
const ENTERTAINMENT = 4301;
const GIFTS = 4302;
const HOLIDAY = 4303;
const SHOPPINGJ = 4304;
const SHOPPINGK = 4305;
const HANDBAGS = 4399;
const GROCERIES = 4401;
const EATINGOUT = 4402;
const TRAIN = 4501;
const TAXI = 4502;
const DOCTOR = 4601;
const EDUCATION = 4602;


@Component({
  selector: 'app-account-report-page',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './account-report-page.component.html',
})
export class AccountReportPageComponent {
  private readonly accounting = inject(AccountingService);

  readonly budget: { [accountId: number]: number } = {
    [EXPENSES]: 146_068,
    [INSURANCE]: 22_000,
    [PARENTS]: 9_600,
    [TAX]: 35_892,
    [PHONE]: 360,
    [TRAIN]: 2_920,
    [TAXI]: 1_000,
    [HOLIDAY]: 12_000,
    [EATINGOUT]: 14_000,
    [GROCERIES]: 3_000,
    [SHOPPINGJ]: 5_000,
    [SHOPPINGK]: 11_000,
    [ENTERTAINMENT]: 500,
    [MAINTENANCE]: 5_000,
    [CLEANER]: 3_796,
    [DOG]: 0,
    [GIFTS]: 7_500,
  };

  formGroup = new FormGroup({
    transactionsType: new FormControl('expenses', { nonNullable: true }),
    month: new FormControl(
      formatDate(FIRST_DAY_OF_LAST_MONTH, 'yyyy-MM', 'en-US'),
      { nonNullable: true }
    ),
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
      }).map((node): ReportRow => {
        return {
          ...node,
          account_id: parseInt(node.ledger_account_id),
          total_credit_num: parseFloat(node.total_credit),
          total_debit_num: parseFloat(node.total_debit)
        }
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

interface ReportRow extends AccountBalanceNode {
  account_id: number
  total_debit_num: number
  total_credit_num: number
}