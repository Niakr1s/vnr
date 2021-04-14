import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Sentence } from './models/sentence';
import { Translation } from './models/translation';
import { TranslationSettings } from './translation-settings/translation-settings';
import { TranslationSettingsService } from './translation-settings/translation-settings.service';

interface TranslationResponse {
  from: string;
  to: string;
  sentence: string;
  translation: string;
}

@Injectable({
  providedIn: 'root',
})
export class TranslationService {
  private _translationSettings?: TranslationSettings;

  constructor(
    private http: HttpClient,
    private translationSettingsService: TranslationSettingsService
  ) {
    this.translationSettingsService.translationSettings$.subscribe({
      next: (t) => (this._translationSettings = t),
    });
  }

  async translate(
    translatorName: string,
    sentence: Sentence,
    to: string
  ): Promise<Translation> {
    if (!this._translationSettings) {
      throw new Error('no translation settings');
    }
    const translator = this._translationSettings.findTranslator(translatorName);
    if (!translator) {
      throw new Error(`no translator with name ${translatorName}`);
    }
    const lang = translator.findLang(to);
    if (!lang) {
      throw new Error(`no lang ${to}`);
    }
    try {
      const res = await this.http
        .get<TranslationResponse>(
          `api/translate/${translator.name}?sentence=${sentence.sentence}&to=${lang.name}&from=ja`
        )
        .toPromise();
      return Translation.create(translatorName, to, res.translation);
    } catch (e) {
      return Translation.createError(translatorName, to, e.message);
    }
  }
}
