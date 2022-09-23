import { Translation } from './translation';

export class Sentence {
  private static id = 0;

  id!: number;
  sentence!: string;
  translations: Record<string, Record<string, Translation>> = {};

  static create(sentence: string): Sentence {
    const s = new Sentence();
    s.sentence = sentence;
    s.id = Sentence.id++;
    return s;
  }

  getTranslation(
    translatorName: string,
    langName: string
  ): Translation | undefined {
    if (!(translatorName in this.translations)) {
      return;
    }
    return this.translations[translatorName][langName];
  }

  setTranslation(translator: string, translation: Translation): void {
    this.translations[translator] ||= {};
    this.translations[translator][translation.to] = translation;
  }

  hasTranslation(translator: string, to: string): boolean {
    let t = this.translations[translator];
    if (!t) return false;

    let l = t[to];
    if (!l) return false;

    return !l.lazy;
  }
}
