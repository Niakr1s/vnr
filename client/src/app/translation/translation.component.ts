import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { Sentence } from '../services/models/sentence';
import { Translation } from '../services/models/translation';
import { SentenceService } from '../services/sentence.service';
import { TranslationSettings } from '../services/translation-settings/translation-settings';
import { TranslationSettingsService } from '../services/translation-settings/translation-settings.service';

@Component({
  selector: 'app-translation',
  templateUrl: './translation.component.html',
  styleUrls: ['./translation.component.css'],
})
export class TranslationComponent implements OnInit, OnDestroy {
  sentence!: Sentence | null;
  translationSettings?: TranslationSettings;

  private subs: Subscription[] = [];

  constructor(
    private sentenceService: SentenceService,
    private translationSettingsService: TranslationSettingsService
  ) { }

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
      this.sentenceService.translationsUpdated$.subscribe({
        next: (sentence) => {
          if (sentence?.id !== this.sentence?.id) {
            this.sentence = sentence;
          }
        },
      })
    );
    this.subs.push(
      this.translationSettingsService.translationSettings$.subscribe({
        next: (t) => {
          this.translationSettings = t;
        },
      })
    );
  }

  getOrderedTranslations(): Translation[] {
    let res: Translation[] = []
    if (!this.translationSettings || !this.sentence) {
      return res;
    }

    for (let translator of this.translationSettings.translators) {
      for (let lang of translator.langs) {
        if (lang.selected) {
          let translation = this.sentence.translations[translator.name][lang.name];
          res.push(translation);
        }
      }
    }
    return res;
  }

  updateTranslation(translatorName: string, sentence: Sentence, to: string) {
    this.sentenceService.translate(translatorName, sentence, to, true);
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
  }
}
