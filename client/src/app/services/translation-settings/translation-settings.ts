import { plainToClass } from 'class-transformer';
import { Lang } from '../models/lang';

export class Translator {
  name!: string;
  translateAlways!: boolean;
  langs!: Lang[];

  static create(name: string, langs: Lang[], translateAlways: boolean): Translator {
    const translator = new Translator();
    translator.name = name;
    translator.langs = langs;
    translator.translateAlways = translateAlways;
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
    this.translateAlways = savedTranslator.translateAlways;
  }

  moveUp(langName: string): void {
    if (this.isFirst(langName)) {
      return;
    }
    const index = this.langs.findIndex((l) => l.name === langName);
    if (index === -1) {
      return;
    }
    swap(this.langs, index, index - 1);
  }

  moveDown(langName: string): void {
    if (this.isLast(langName)) {
      return;
    }
    const index = this.langs.findIndex((l) => l.name === langName);
    if (index === -1) {
      return;
    }
    swap(this.langs, index, index + 1);
  }

  isFirst(langName: string): boolean {
    const index = this.langs.findIndex((l) => l.name === langName);
    return index === 0;
  }

  isLast(langName: string): boolean {
    const index = this.langs.findIndex((l) => l.name === langName);
    if (index === -1) {
      return false;
    }
    return index === this.langs.length - 1;
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

  moveUp(translatorName: string): void {
    if (this.isFirst(translatorName)) {
      return;
    }
    const index = this.translators.findIndex((l) => l.name === translatorName);
    if (index === -1) {
      return;
    }
    swap(this.translators, index, index - 1);
  }

  moveDown(translatorName: string): void {
    if (this.isLast(translatorName)) {
      return;
    }
    const index = this.translators.findIndex((l) => l.name === translatorName);
    if (index === -1) {
      return;
    }
    swap(this.translators, index, index + 1);
  }

  isFirst(translatorName: string): boolean {
    const index = this.translators.findIndex((t) => t.name === translatorName);
    return index === 0;
  }

  isLast(translatorName: string): boolean {
    const index = this.translators.findIndex((t) => t.name === translatorName);
    if (index === -1) {
      return false;
    }
    return index === this.translators.length - 1;
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

const swap = <T>(arr: T[], index1: number, index2: number): void => {
  [arr[index1], arr[index2]] = [arr[index2], arr[index1]];
};
