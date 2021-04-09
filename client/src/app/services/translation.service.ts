import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Sentence } from './models/sentence';
import { Translation } from './models/translation';
import {
  Translators,
  TranslatorsRepoService,
} from './translators-repo.service';

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
  private _translators?: Translators;

  constructor(
    private http: HttpClient,
    private translatorsRepo: TranslatorsRepoService
  ) {
    this.translatorsRepo.translators$.subscribe({
      next: (t) => (this._translators = t),
    });
  }

  async translate(
    translatorName: string,
    sentence: Sentence,
    to: string
  ): Promise<Translation> {
    if (!this._translators) {
      throw new Error('no translators');
    }
    const translator = this._translators.findTranslator(translatorName);
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
          `api/translate/${translator.name}?sentence=${sentence.sentence}&to=${lang.name}`
        )
        .toPromise();
      return Translation.create(to, res.translation);
    } catch (e) {
      return Translation.createError(to, e.message);
    }
  }
}
