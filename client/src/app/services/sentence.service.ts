import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { ClipboardService } from './clipboard.service';
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
export class SentenceService {
  private sentences: Sentence[] = [];
  private maxSentences = 99;
  private sentenceCounter = 0;

  private _currentIndex = -1;
  private get currentIndex(): number {
    return this._currentIndex;
  }
  private set currentIndex(index: number) {
    if (index >= this.sentences.length) {
      index = this.sentences.length - 1;
    }
    if (index < -1) {
      index = -1;
    }
    this._currentIndex = index;
    this.currentIndexSubject.next(index);
    this.currentSencenceSubject.next(this.currentSentence);
  }

  private currentIndexSubject = new BehaviorSubject<number>(-1);
  get currentIndex$(): Observable<number> {
    return this.currentIndexSubject.asObservable();
  }

  private totalSentencesSubject = new BehaviorSubject<number>(0);
  get totalSentences$(): Observable<number> {
    return this.totalSentencesSubject.asObservable();
  }

  private sentenceCounterSubject = new BehaviorSubject<number>(
    this.sentenceCounter
  );
  get sentenceCounter$(): Observable<number> {
    return this.sentenceCounterSubject.asObservable();
  }

  private get currentSentence(): Sentence | null {
    return this.sentences[this.currentIndex];
  }

  private currentSencenceSubject = new BehaviorSubject<Sentence | null>(null);
  get currentSentence$(): Observable<Sentence | null> {
    return this.currentSencenceSubject.asObservable();
  }

  private translationsUpdatedSubject = new BehaviorSubject<Sentence | null>(null);
  get translationsUpdated$(): Observable<Sentence | null> {
    return this.translationsUpdatedSubject.asObservable();
  }

  private translationSettings?: TranslationSettings;

  constructor(
    private http: HttpClient,
    clipboardService: ClipboardService,
    private translationSettingsService: TranslationSettingsService
  ) {
    this.translationSettingsService.translationSettings$.subscribe({
      next: (translationSettings) => {
        this.translationSettings = translationSettings;
        if (this.currentSentence) {
          this.translateMissedLanguages(this.currentSentence);
        }
      },
    });
    this.currentSentence$.subscribe({
      next: (s) => {
        if (s) {
          this.translateMissedLanguages(s);
        }
      },
    });
    clipboardService.clipboard.subscribe({
      next: async (s) => {
        await this.onNewSencence(s);
      },
    });
  }

  private async onNewSencence(s: string) {
    const sentence = Sentence.create(s);
    this.pushSentence(sentence);

    this.translateMissedLanguages(sentence);
  }

  private pushSentence(sentence: Sentence): void {
    this.sentences.push(sentence);

    this.sentenceCounter++;
    this.sentenceCounterSubject.next(this.sentenceCounter);

    this.currentIndex = this.sentences.length - 1;

    if (this.sentences.length > this.maxSentences) {
      this.deleteSentenceAt(0);
    }

    this.totalSentencesSubject.next(this.sentences.length);
  }

  private deleteSentenceAt(index: number): void {
    if (
      index < 0 ||
      index >= this.sentences.length ||
      this.sentences.length === 0
    ) {
      return;
    }

    this.sentences.splice(index, 1);

    this.currentIndex =
      index > 0 || (index === 0 && this.sentences.length === 0)
        ? this.currentIndex - 1
        : this.currentIndex; // nessesary, to emit correct data

    this.totalSentencesSubject.next(this.sentences.length);
  }

  private translateMissedLanguages(sentence: Sentence): void {
    if (!this.translationSettings) return;

    for (const t of this.translationSettings.translators) {
      for (const lang of t.langs) {
        if (lang.selected && !sentence.hasTranslation(t.name, lang.name)) {
          this.translate(t.name, sentence, lang.name, t.translateAlways);
        }
      }
    }
  }

  translate(
    translatorName: string,
    sentence: Sentence,
    to: string,
    _doActualTranslate = true,
  ) {
    if (!_doActualTranslate) {
      this.setTranslation(
        translatorName,
        sentence.id,
        Translation.createLazy(translatorName, to)
      );
      return;
    }

    this.setTranslation(
      translatorName,
      sentence.id,
      Translation.createPending(translatorName, to)
    );

    this.actualTranslate(translatorName, sentence, to)
      .then((translation) => {
        console.log("got translation: ", translation);
        this.setTranslation(translatorName, sentence.id, translation);
      })
      .catch((e) => {
        console.error(e);
      });
  }

  private async actualTranslate(
    translatorName: string,
    sentence: Sentence,
    to: string
  ): Promise<Translation> {
    if (!this.translationSettings) {
      throw new Error('no translation settings');
    }
    const translator = this.translationSettings.findTranslator(translatorName);
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
      return Translation.create(translatorName, to, res.translation);
    } catch (e: any) {
      console.error(e);
      return Translation.createError(translatorName, to, e.message);
    }
  }

  private setTranslation(
    translatorName: string,
    id: number,
    translation: Translation
  ): void {
    const sentence = this.sentences.find((s) => s.id === id);
    if (!sentence) {
      return;
    }

    sentence.setTranslation(translatorName, translation);
    this.translationsUpdatedSubject.next(sentence);
  }

  private isCurrent(sentence: Sentence): boolean {
    const currentSentence = this.currentSentence;
    return currentSentence != null && currentSentence === sentence;
  }

  deleteCurrentSentence(): void {
    this.deleteSentenceAt(this.currentIndex);
  }

  hasPrev(): boolean {
    if (this.sentences.length === 0) {
      return false;
    }
    return this.currentIndex > 0;
  }

  prev(): void {
    if (!this.hasPrev()) {
      return;
    }
    this.currentIndex--;
  }

  hasNext(): boolean {
    if (this.sentences.length === 0) {
      return false;
    }
    return this.currentIndex < this.sentences.length - 1;
  }

  next(): void {
    if (!this.hasNext()) {
      return;
    }
    this.currentIndex++;
  }

  last(): void {
    if (this.currentIndex === this.sentences.length - 1) {
      return;
    }
    this.currentIndex = this.sentences.length - 1;
  }
}
