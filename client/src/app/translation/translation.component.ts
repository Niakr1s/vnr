import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { Sentence } from '../services/models/sentence';
import { Translation } from '../services/models/translation';
import { SentenceService } from '../services/sentence.service';
import {
  Translator,
  TranslatorsRepoService,
} from '../services/translators-repo.service';

@Component({
  selector: 'app-translation',
  templateUrl: './translation.component.html',
  styleUrls: ['./translation.component.css'],
})
export class TranslationComponent implements OnInit, OnDestroy {
  sentence!: Sentence | null;
  translator?: Translator;

  private subs: Subscription[] = [];

  constructor(
    private sentenceService: SentenceService,
    private translatorsRepo: TranslatorsRepoService // private languageService: LanguageService
  ) {}

  ngOnInit(): void {
    // this.sentence?.translations['a'].
    this.subs.push(
      this.sentenceService.currentSentence$.subscribe({
        next: (sentence) => {
          this.sentence = sentence;
        },
      })
    );
    this.subs.push(
      this.translatorsRepo.translators$.subscribe({
        next: (t) => {
          this.translator = t?.getSelectedTranslator();
        },
      })
    );
  }

  getTranslation(lang: string): Translation | undefined {
    if (!this.translator) {
      return;
    }
    const translations = this.sentence?.translations[this.translator.name];
    if (!translations) {
      return;
    }
    return translations[lang];
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
  }
}
