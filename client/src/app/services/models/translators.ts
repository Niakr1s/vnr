import { Lang } from './lang';

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
