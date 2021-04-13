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
      this.translationSettingsService.translationSettings$.subscribe({
        next: (t) => {
          this.translationSettings = t;
        },
      })
    );
  }

  getTranslations(): Translation[] {
    const res: Translation[] = [];

    this.translationSettings?.forEachLang((name, lang) => {
      if (!lang.selected) {
        return;
      }
      const translation = this.sentence?.getTranslation(name, lang.name);
      if (translation) {
        res.push(translation);
      }
    });

    return res;
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
  }
}
