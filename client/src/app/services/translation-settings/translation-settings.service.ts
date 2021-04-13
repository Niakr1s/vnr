import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { Lang } from '../models/lang';
import { TranslationSettings, Translator } from './translation-settings';

@Injectable({
  providedIn: 'root',
})
export class TranslationSettingsService {
  private _translationSettings?: TranslationSettings;

  private translationSettingsSubject = new BehaviorSubject<
    TranslationSettings | undefined
  >(undefined);

  get translationSettings$(): Observable<TranslationSettings | undefined> {
    return this.translationSettingsSubject.asObservable();
  }

  constructor(private http: HttpClient) {
    this.getTranslationSettings();
    this.translationSettings$.subscribe({
      next: (t) => {
        this._translationSettings = t;
        t?.save();
      },
    });
  }

  moveTranslatorUp(translatorName: string): void {
    this._translationSettings?.moveUp(translatorName);
    this.translationSettingsSubject.next(this._translationSettings);
  }

  moveTranslatorDown(translatorName: string): void {
    this._translationSettings?.moveDown(translatorName);
    this.translationSettingsSubject.next(this._translationSettings);
  }

  moveLangUp(translatorName: string, langName: string): void {
    this._translationSettings?.findTranslator(translatorName)?.moveUp(langName);
    this.translationSettingsSubject.next(this._translationSettings);
  }

  moveLangDown(translatorName: string, langName: string): void {
    this._translationSettings
      ?.findTranslator(translatorName)
      ?.moveDown(langName);
    this.translationSettingsSubject.next(this._translationSettings);
  }

  toggleLanguage(translatorName: string, langName: string): void {
    const lang = this._translationSettings
      ?.findTranslator(translatorName)
      ?.findLang(langName);
    if (!lang) {
      return;
    }
    lang.selected = !lang.selected;
    this.translationSettingsSubject.next(this._translationSettings);
  }

  private async getTranslationSettings(): Promise<void> {
    try {
      const translatorNames = await this.http
        .get<string[]>(`api/knownTranslators`)
        .toPromise();
      console.log('got known translators', translatorNames);

      const translatorList = await Promise.all(
        translatorNames.map(async (name) => {
          const langs: string[] = await this.http
            .get<string[]>(`api/langs/${name}`)
            .toPromise();
          const translator = Translator.create(
            name,
            langs.map((n) => Lang.create(n))
          );
          return translator;
        })
      );
      const translationSettings = TranslationSettings.create(translatorList);
      translationSettings.load();
      console.log('got translationSettings', translationSettings);
      this.translationSettingsSubject.next(translationSettings);
    } catch (e) {
      console.error(`couldn't get known translators`, e);
    }
  }

  setLangs(translatorName: string, langs: Lang[]): void {
    const translator = this._translationSettings?.findTranslator(
      translatorName
    );
    if (!translator) {
      return;
    }
    translator.langs = langs;
    this.translationSettingsSubject.next(this._translationSettings);
  }
}
