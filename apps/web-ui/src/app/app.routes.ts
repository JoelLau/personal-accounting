import { Route } from '@angular/router';
import { DashboardPageComponent } from './pages/dashboard-page-component/dashboard-page.component';

export const appRoutes: Route[] = [
  {
    path: '**',
    component: DashboardPageComponent,
  },
];
