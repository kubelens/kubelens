import React from 'react';
import renderer from 'react-test-renderer';

import Home from './index';

test('Navbar Exists', () => {
  let component = renderer.create(Home);

  let tree = component.toJSON();

  expect(tree).toMatchSnapshot();
});