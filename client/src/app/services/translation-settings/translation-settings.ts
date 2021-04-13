import { plainToClass } from 'class-transformer';
import { Lang } from '../models/lang';

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
    return findLang(this.langs, name);
  }

  getSelectedLangs(): string[] {
    return this.langs.filter((l) => l.selected).map((l) => l.name);
  }

  mergeWithSavedTranslator(savedTranslator: Translator): void {
    this.langs = mergeLangs(savedTranslator.langs, this.langs);
  }
}

export class TranslationSettings {
  private static saveKey = 'translationSettings';

  translators!: Translator[];

  static create(translatorList: Translator[]): TranslationSettings {
    const res = new TranslationSettings();
    res.translators = translatorList;
    return res;
  }

  findTranslator(name: string): Translator | undefined {
    return findTranslator(this.translators, name);
  }

  save(): void {
    localStorage.setItem(TranslationSettings.saveKey, JSON.stringify(this));
  }

  load(): void {
    const savedTranslatorJson = localStorage.getItem(
      TranslationSettings.saveKey
    );
    if (!savedTranslatorJson) {
      return;
    }
    const savedTranslator = plainToClass(
      TranslationSettings,
      JSON.parse(savedTranslatorJson)
    );
    this.mergeWithSavedSettings(savedTranslator);
  }

  private mergeWithSavedSettings(savedSettings: TranslationSettings): void {
    for (const savedTranslator of savedSettings.translators) {
      const translator = this.findTranslator(savedTranslator.name);
      translator?.mergeWithSavedTranslator(savedTranslator);
    }
  }

  forEachLang(fn: (translatorName: string, lang: Lang) => void): void {
    for (const t of this.translators) {
      for (const l of t.langs) {
        fn(t.name, l);
      }
    }
  }
}

const findTranslator = (
  translators: Translator[],
  name: string
): Translator | undefined => {
  return translators.find((t) => t.name === name);
};

const findLang = (langs: Lang[], name: string): Lang | undefined => {
  return langs.find((l) => l.name === name);
};

export const mergeLangs = (from: Lang[], to: Lang[]): Lang[] => {
  const res: Lang[] = [];

  for (const fromLang of from) {
    const langExistsInTo = findLang(to, fromLang.name);
    if (langExistsInTo) {
      res.push(fromLang);
    }
  }

  for (const toLang of to) {
    const langExistsInRes = findLang(res, toLang.name);
    if (!langExistsInRes) {
      res.push(toLang);
    }
  }

  return res;
};
