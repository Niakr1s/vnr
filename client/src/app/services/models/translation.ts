export class Translation {
  to!: string;

  translation?: string;
  pending = false;
  error?: string;

  static create(to: string, translation: string): Translation {
    const t = new Translation();
    t.to = to;
    t.translation = translation;
    return t;
  }

  static createPending(to: string): Translation {
    const t = new Translation();
    t.to = to;
    t.pending = true;
    return t;
  }

  static createError(to: string, error: string): Translation {
    const t = new Translation();
    t.to = to;
    t.error = error;
    return t;
  }
}
