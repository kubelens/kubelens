import store from './store';

describe('store should', () => {
  test('send "test" dispatch successfully', () => {
    const s = store();
    const d = s.dispatch({type: "test"});
    expect(d.type).toEqual("test")
  });
});