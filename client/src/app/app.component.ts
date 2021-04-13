import { Component, OnInit } from '@angular/core';
import { SettingsService } from './services/settings.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent implements OnInit {
  title = 'client';

  settingsVisible!: boolean;

  constructor(private settingsService: SettingsService) {
    this.settingsVisible = settingsService.visible;
  }

  ngOnInit(): void {
    this.settingsService.visible$.subscribe({
      next: (v) => {
        this.settingsVisible = v;
      },
    });
  }

  onSettingsClick(): void {
    this.settingsService.visible = true;
  }
}
