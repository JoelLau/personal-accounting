import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Store } from '@ngxs/store';
import { LedgerState } from '../../../store/ledger.state';

@Component({
  selector: 'app-accounts-page',
  imports: [CommonModule],
  templateUrl: './accounts-page.component.html',
  styleUrl: './accounts-page.component.scss',
})
export class AccountsPageComponent {
  private readonly store = inject(Store);

  protected readonly account$ = this.store.select(LedgerState.getAccounts);
}
