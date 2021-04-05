import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { ClipboardService } from './clipboard.service';
import { Sentence } from './models/sentence';
import { Translation } from './models/translation';
import { TranslationService } from './translation.service';

@Injectable({
  providedIn: 'root',
})
export class SentenceService {
  private sentences: Sentence[] = [];

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

  get currentSentence(): Sentence | null {
    return this.sentences[this.currentIndex];
  }

  private currentSencenceSubject = new BehaviorSubject<Sentence | null>(null);
  get currentSentence$(): Observable<Sentence | null> {
    return this.currentSencenceSubject.asObservable();
  }

  constructor(
    clipboardService: ClipboardService,
    private translationService: TranslationService
  ) {
    clipboardService.clipboard.subscribe({
      next: async (s) => {
        const sentence = Sentence.create(s);
        this.pushSentence(sentence);

        await this.translate(sentence, ['en', 'ru']);
      },
    });
  }

  private pushSentence(sentence: Sentence): void {
    this.sentences.push(sentence);
    this.currentIndex = this.sentences.length - 1;
    this.currentSencenceSubject.next(sentence);
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

    // if we deleted last sentence - adjusting currentIndex
    if (
      index < this.currentIndex ||
      this.sentences.length === 0 ||
      this.currentIndex === this.sentences.length
    ) {
      this.currentIndex--;
    }

    this.currentSencenceSubject.next(this.currentSentence);
    this.totalSentencesSubject.next(this.sentences.length);
  }

  deleteCurrentSentence(): void {
    this.deleteSentenceAt(this.currentIndex);
  }

  private async translate(sentence: Sentence, tos: string[]): Promise<void> {
    const translations = await Promise.all(
      tos.map((to) => this.translationService.translate(sentence, to))
    );
    this.setTranslation(sentence.id, translations);
  }

  setTranslation(id: number, translations: Translation[]): void {
    const sentence = this.sentences.find((s) => s.id === id);
    if (!sentence) {
      return;
    }

    for (const translation of translations) {
      sentence.translations[translation.to] = translation;
    }

    if (this.isCurrent(sentence)) {
      this.currentSencenceSubject.next(sentence);
    }
  }

  isCurrent(sentence: Sentence): boolean {
    const currentSentence = this.currentSentence;
    return currentSentence != null && currentSentence === sentence;
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
