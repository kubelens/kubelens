import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import { Logo } from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = { }

  const wrapper = shallow(<Logo {...props} />)

  return {
    props,
    wrapper
  }
}

describe('Logo should', () => {
  test('render', () => {
    const { wrapper } = setup();

    expect(wrapper.find('img').length).toBe(2);
  })
})