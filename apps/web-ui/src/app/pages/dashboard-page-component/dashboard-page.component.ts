import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { map } from 'rxjs';
import { AccountingService } from '../../services/accounting-api-service/services';

@Component({
  selector: 'app-dashboard-page',
  imports: [CommonModule],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.scss',
})
export class DashboardPageComponent {
  private readonly accounting = inject(AccountingService);

  protected readonly ledgerAccounts$ = this.accounting
    .apiV1AccountingLedgerAccountsGet()
    .pipe(map((response) => response?.data ?? []));

  protected readonly postings$ = this.accounting
    .apiV1AccountingPostingsGet()
    .pipe(map((response) => response?.data ?? []));
}
