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
  private currentIndex = -1;

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
}
