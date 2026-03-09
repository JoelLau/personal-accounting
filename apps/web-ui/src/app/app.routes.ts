import { Route } from '@angular/router';
import { StackedLayoutComponent } from './layouts/stacked-layout-component/stacked-layout.component';
import { TransactionsPageComponentComponent } from './pages/transactions-page-component/transactions-page-component.component';

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
