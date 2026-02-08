import { Component, inject, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';
import { Store } from '@ngxs/store';
import { FetchLedgerData } from './store/ledger.actions';

@Component({
  imports: [RouterModule],
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent implements OnInit {
  private store = inject(Store);

  ngOnInit(): void {
    this.store.dispatch(new FetchLedgerData());
  }
}
