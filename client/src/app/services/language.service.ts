import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

export type Languages = Record<string, boolean>;

@Injectable({
  providedIn: 'root',
})
export class LanguageService {
  private languages: Languages = {
    en: true,
    ru: true,
  };

  private languagesSubject = new BehaviorSubject<Languages>(this.languages);
  get languages$(): Observable<Languages> {
    return this.languagesSubject.asObservable();
  }

  constructor() {}

  enabledLanguages(): string[] {
    return Object.entries(this.languages)
      .filter((l) => l[1])
      .map((l) => l[0]);
  }

  toggleLanguage(lang: string): boolean {
    if (!(lang in this.languages)) {
      return false;
    }

    this.languages[lang] = !this.languages[lang];
    this.languagesSubject.next(this.languages);

    return true;
  }
}
