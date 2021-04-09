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

  setTranslation(translator: string, translation: Translation): void {
    this.translations[translator] ||= {};
    this.translations[translator][translation.to] = translation;
  }
}
