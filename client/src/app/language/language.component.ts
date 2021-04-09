import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { Translator } from '../services/models/translators';
import { TranslatorsRepoService } from '../services/translators-repo.service';

@Component({
  selector: 'app-language',
  templateUrl: './language.component.html',
  styleUrls: ['./language.component.css'],
})
export class LanguageComponent implements OnInit, OnDestroy {
  translator?: Translator;

  subs: Subscription[] = [];

  constructor(private translatorsRepo: TranslatorsRepoService) {}

  ngOnInit(): void {
    this.subs.push(
      this.translatorsRepo.translators$.subscribe({
        next: (t) => {
          this.translator = t?.getSelectedTranslator();
        },
      })
    );
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
  }

  onToggle(lang: string): void {
    if (!this.translator) {
      return;
    }
    this.translatorsRepo.toggleLanguage(this.translator.name, lang);
  }
}
