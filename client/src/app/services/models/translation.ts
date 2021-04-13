export class Translation {
  translatorName!: string;
  to!: string;

  translation?: string;
  pending = false;
  error?: string;

  static create(
    translatorName: string,
    to: string,
    translation: string
  ): Translation {
    const t = new Translation();
    t.translatorName = translatorName;
    t.to = to;
    t.translation = translation;
    return t;
  }

  static createPending(translatorName: string, to: string): Translation {
    const t = new Translation();
    t.translatorName = translatorName;
    t.to = to;
    t.pending = true;
    return t;
  }

  static createError(
    translatorName: string,
    to: string,
    error: string
  ): Translation {
    const t = new Translation();
    t.translatorName = translatorName;
    t.to = to;
    t.error = error;
    return t;
  }
}
