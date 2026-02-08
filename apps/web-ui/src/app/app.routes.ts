import { Route } from '@angular/router';
import { StackedLayoutComponent } from './web-ui/layouts/stacked-layout/stacked-layout.component';
import { AccountsPageComponent } from './web-ui/pages/accounts-page/accounts-page.component';

export const appRoutes: Route[] = [
  {
    path: '',
    component: StackedLayoutComponent,
    children: [
      {
        path: '',
        component: AccountsPageComponent,
      },
      {
        path: '**',
        redirectTo: '',
      },
    ],
  },
  {
    path: '**',
    redirectTo: '',
  },
];
