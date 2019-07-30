import * as Exports from './index';

describe('index', () => {
  it('error boundary', () => {
    expect(Object.keys(Exports).length).toBeGreaterThanOrEqual(1);
  });
});
