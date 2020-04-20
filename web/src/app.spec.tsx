import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import App from './app'

Enzyme.configure({ adapter: new Adapter() })

function setup() {
  const props = {
    history: {},
    authClient: {}
  }

  const enzymeWrapper = shallow(<App {...props} />)

  return {
    props,
    enzymeWrapper
  }
}

describe('app.tsx should', () => {
  test('have props', () => {
    const { enzymeWrapper } = setup();

    expect(enzymeWrapper.find({ history: {} }).length).toBe(3);
    expect(enzymeWrapper.find({ authClient: {} }).length).toBe(2);
  })
})