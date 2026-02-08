import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ApiService } from '../../services/accounting-api-service/services';

@Component({
  selector: 'app-dashboard-page',
  imports: [CommonModule],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.scss',
})
export class DashboardPageComponent implements OnInit {
  protected readonly healthyz = signal<unknown | null>(null);

  private apiService = inject(ApiService);

  async ngOnInit() {
    this.healthyz.set(await this.apiService.apiReadyzGet());
  }
}
