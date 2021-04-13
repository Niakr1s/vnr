import { Lang } from '../models/lang';
import { mergeLangs } from './translation-settings';

describe('mergeLangs', () => {
  it('should work', () => {
    const from: Lang[] = [
      Lang.create('ru'),
      Lang.create('en'),
      Lang.create('jp'),
    ];

    from[1].selected = true;

    const to: Lang[] = [
      Lang.create('en'),
      Lang.create('de'),
      Lang.create('ru'),
    ];

    const expected: Lang[] = [
      Lang.create('ru'),
      Lang.create('en'),
      Lang.create('de'),
    ];
    expected[1].selected = true;

    const got = mergeLangs(from, to);

    expect(got).toEqual(expected);
  });
});
