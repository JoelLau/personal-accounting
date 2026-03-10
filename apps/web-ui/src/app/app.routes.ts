import { Route } from '@angular/router';
import { StackedLayoutComponent } from './layouts/stacked-layout-component/stacked-layout.component';
import { TransactionsPageComponentComponent } from './pages/transactions-page-component/transactions-page-component.component';
import { AccountReportPageComponent } from './pages/account-report-page/account-report-page.component';

export const appRoutes: Route[] = [
  {
    path: '',
    component: StackedLayoutComponent,
    children: [
      {
        path: 'transactions',
        component: TransactionsPageComponentComponent,
      },
      {
        path: 'account-report',
        component: AccountReportPageComponent,
      },
      {
        path: '**',
        redirectTo: 'transactions',
      },
    ],
  },
  {
    path: '**',
    redirectTo: 'transactions',
  },
];
