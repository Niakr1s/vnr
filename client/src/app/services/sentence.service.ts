import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { ClipboardService } from './clipboard.service';
import { Sentence } from './models/sentence';
import { Translation } from './models/translation';
import { Translator } from './models/translators';
import { TranslationService } from './translation.service';
import { TranslatorsRepoService } from './translators-repo.service';

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

  private currentTranslator?: Translator;

  constructor(
    clipboardService: ClipboardService,
    private translationService: TranslationService,
    private translatorRepo: TranslatorsRepoService
  ) {
    this.translatorRepo.translators$.subscribe({
      next: (translators) => {
        this.currentTranslator = translators?.getSelectedTranslator();
        this.translateMissedLanguages();
      },
    });
    this.currentSentence$.subscribe({
      next: () => {
        this.translateMissedLanguages();
      },
    });
    clipboardService.clipboard.subscribe({
      next: async (s) => {
        const sentence = Sentence.create(s);
        this.pushSentence(sentence);
      },
    });
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

  private translateMissedLanguages(): void {
    if (!this.currentTranslator || !this.currentSentence) {
      return;
    }
    const { langs, name } = this.currentTranslator;

    const tos: string[] = [];
    for (const lang of langs) {
      if (
        lang.selected &&
        !this.currentSentence.hasTranslation(name, lang.name)
      ) {
        tos.push(lang.name);
      }
    }
    this.translate(name, this.currentSentence, tos);
  }

  private translate(
    translator: string,
    sentence: Sentence,
    tos: string[]
  ): void {
    tos.forEach((to) => {
      {
        this.setTranslation(
          translator,
          sentence.id,
          Translation.createPending(to)
        );
        this.translationService
          .translate(translator, sentence, to)
          .then((translation) => {
            this.setTranslation(translator, sentence.id, translation);
          });
      }
    });
  }

  private setTranslation(
    translator: string,
    id: number,
    translation: Translation
  ): void {
    const sentence = this.sentences.find((s) => s.id === id);
    if (!sentence) {
      return;
    }

    sentence.setTranslation(translator, translation);

    if (this.isCurrent(sentence)) {
      this.currentSencenceSubject.next(sentence);
    }
  }

  private isCurrent(sentence: Sentence): boolean {
    const currentSentence = this.currentSentence;
    return currentSentence != null && currentSentence === sentence;
  }

  deleteCurrentSentence(): void {}

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
