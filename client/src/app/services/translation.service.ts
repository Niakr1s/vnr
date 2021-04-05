import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Sentence } from './models/sentence';
import { Translation } from './models/translation';

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
  constructor(private http: HttpClient) {}

  async translate(sentence: Sentence, to: string): Promise<Translation> {
    try {
      const res = await this.http
        .get<TranslationResponse>(
          `api/translate?sentence=${sentence.sentence}&to=${to}`
        )
        .toPromise();
      return Translation.create(to, res.translation);
    } catch (e) {
      return Translation.createError(to, e.message);
    }
  }
}
