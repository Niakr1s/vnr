import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

export class Lang {
  name!: string;
  selected!: boolean;

  static create(name: string): Lang {
    const lang = new Lang();
    lang.name = name;
    lang.selected = false;
    return lang;
  }
}

export class Translator {
  name!: string;
  langs!: Lang[];

  static create(name: string, langs: Lang[]): Translator {
    const translator = new Translator();
    translator.name = name;
    translator.langs = langs;
    return translator;
  }

  findLang(name: string): Lang | undefined {
    return this.langs.find((l) => l.name === name);
  }

  getSelectedLangs(): string[] {
    return this.langs.filter((l) => l.selected).map((l) => l.name);
  }
}

export class Translators {
  translators!: Translator[];
  private selectedTranslator!: string;

  static create(translatorList: Translator[]): Translators {
    const res = new Translators();
    res.translators = translatorList;
    if (res.translators.length > 0) {
      res.selectedTranslator = res.translators[0].name;
    }
    return res;
  }

  getSelectedTranslator(): Translator | undefined {
    return this.findTranslator(this.selectedTranslator);
  }

  findTranslator(name: string): Translator | undefined {
    return this.translators.find((t) => t.name === name);
  }
}

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
      next: (t) => (this._translators = t),
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
      this.translatorsSubject.next(translators);
    } catch (e) {
      console.error(`couldn't get known translators`, e);
    }
  }
}
