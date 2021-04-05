import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { Languages, LanguageService } from '../services/language.service';

@Component({
  selector: 'app-language',
  templateUrl: './language.component.html',
  styleUrls: ['./language.component.css'],
})
export class LanguageComponent implements OnInit, OnDestroy {
  languages!: Languages;

  subs: Subscription[] = [];

  constructor(private languageService: LanguageService) {}

  ngOnInit(): void {
    this.subs.push(
      this.languageService.languages$.subscribe({
        next: (l) => {
          console.log(l);
          this.languages = l;
        },
      })
    );
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
  }

  onToggle(lang: string): void {
    this.languageService.toggleLanguage(lang);
  }
}
