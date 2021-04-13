import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { TranslationSettings } from '../services/translation-settings/translation-settings';
import { TranslationSettingsService } from '../services/translation-settings/translation-settings.service';

@Component({
  selector: 'app-language',
  templateUrl: './language.component.html',
  styleUrls: ['./language.component.css'],
})
export class LanguageComponent implements OnInit, OnDestroy {
  translationSettings?: TranslationSettings;

  subs: Subscription[] = [];

  constructor(private translationSettingsService: TranslationSettingsService) {}

  ngOnInit(): void {
    this.subs.push(
      this.translationSettingsService.translationSettings$.subscribe({
        next: (t) => {
          this.translationSettings = t;
        },
      })
    );
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
  }

  onToggle(translatorName: string, lang: string): void {
    this.translationSettingsService.toggleLanguage(translatorName, lang);
  }
}
