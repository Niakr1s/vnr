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
