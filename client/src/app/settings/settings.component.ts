import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { plainToClass } from 'class-transformer';
import { Lang } from '../services/models/lang';
import { SettingsService } from '../services/settings.service';
import { TranslationSettings } from '../services/translation-settings/translation-settings';
import { TranslationSettingsService } from '../services/translation-settings/translation-settings.service';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.css'],
})
export class SettingsComponent implements OnInit {
  @ViewChild('outer') outer!: ElementRef;

  translationSettings?: TranslationSettings;

  constructor(
    private settingsService: SettingsService,
    public translationSettingsService: TranslationSettingsService
  ) {}

  ngOnInit(): void {
    this.translationSettingsService.translationSettings$.subscribe({
      next: (s) => (this.translationSettings = s),
    });
  }

  onOuterClick(ev: MouseEvent): void {
    if (ev.target !== this.outer.nativeElement) {
      return;
    }
    this.onClose();
  }

  onClose(): void {
    console.log('click');
    this.settingsService.visible = false;
  }

  onToggle(translatorName: string, langName: string): void {
    this.translationSettingsService.toggleLanguage(translatorName, langName);
  }

  onLangsSorted(translatorName: string, langs: Lang[]): void {
    langs = plainToClass(Lang, langs);
    this.translationSettingsService.setLangs(translatorName, langs);
  }
}
