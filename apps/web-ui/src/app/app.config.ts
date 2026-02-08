import { provideHttpClient } from '@angular/common/http';
import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideStore } from '@ngxs/store';
import { environment } from '../environments/environment';
import { appRoutes } from './app.routes';
import { provideApiConfiguration } from './services/accounting-api-service/api-configuration';
import { LedgerState } from './store/ledger.state';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(appRoutes),
    provideHttpClient(),
    provideApiConfiguration(environment.apiUrl),
    provideStore([LedgerState]),
  ],
};
