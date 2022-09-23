export class Translation {
  translatorName!: string;
  to!: string;

  translation?: string;

  // lazy used for on-demand translations
  lazy = false;

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

  isError(): boolean {
    return !!this.error;
  }

  static createLazy(translatorName: string, to: string): Translation {
    const t = new Translation();
    t.translatorName = translatorName;
    t.to = to;
    t.lazy = true;
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
