import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { Languages, LanguageService } from '../services/language.service';
import { Sentence } from '../services/models/sentence';
import { SentenceService } from '../services/sentence.service';

@Component({
  selector: 'app-translation',
  templateUrl: './translation.component.html',
  styleUrls: ['./translation.component.css'],
})
export class TranslationComponent implements OnInit, OnDestroy {
  sentence!: Sentence | null;
  languages!: Languages;

  private subs: Subscription[] = [];

  constructor(
    private sentenceService: SentenceService,
    private languageService: LanguageService
  ) {}

  ngOnInit(): void {
    this.subs.push(
      this.sentenceService.currentSentence$.subscribe({
        next: (sentence) => {
          this.sentence = sentence;
        },
      })
    );
    this.subs.push(
      this.languageService.languages$.subscribe({
        next: (l) => {
          this.languages = l;
        },
      })
    );
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
  }
}
