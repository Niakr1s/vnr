export class Lang {
  name!: string;
  description!: string;
  selected!: boolean;

  static create(name: string, description: string = ''): Lang {
    const lang = new Lang();
    lang.name = name;
    lang.description = description;
    lang.selected = false;
    return lang;
  }
}
