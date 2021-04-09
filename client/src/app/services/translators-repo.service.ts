import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { Lang } from './models/lang';
import { Translator, Translators } from './models/translators';

@Injectable({
  providedIn: 'root',
})
export class TranslatorsRepoService {
  private _translators?: Translators;

  private translatorsSubject = new BehaviorSubject<Translators | undefined>(
    undefined
  );

  get translators$(): Observable<Translators | undefined> {
    return this.translatorsSubject.asObservable();
  }

  constructor(private http: HttpClient) {
    this.getKnownTranslators();
    this.translators$.subscribe({
      next: (t) => {
        this._translators = t;
        t?.save();
      },
    });
  }

  toggleLanguage(translatorName: string, langName: string): void {
    const lang = this._translators
      ?.findTranslator(translatorName)
      ?.findLang(langName);
    if (!lang) {
      return;
    }
    lang.selected = !lang.selected;
    this.translatorsSubject.next(this._translators);
  }

  private async getKnownTranslators(): Promise<void> {
    try {
      const translatorNames = await this.http
        .get<string[]>(`api/knownTranslators`)
        .toPromise();
      console.log('got known translators', translatorNames);

      const translatorList = await Promise.all(
        translatorNames.map(async (name) => {
          const langs: string[] = await this.http
            .get<string[]>(`api/langs/${name}`)
            .toPromise();
          const translator = Translator.create(
            name,
            langs.map((n) => Lang.create(n))
          );
          return translator;
        })
      );
      const translators = Translators.create(translatorList);
      console.log('got translators', translators);
      translators.load();
      this.translatorsSubject.next(translators);
    } catch (e) {
      console.error(`couldn't get known translators`, e);
    }
  }
}
