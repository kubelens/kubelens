import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import Text from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    label: 'text label',
    value: 'text value'
  }

  const wrapper = shallow(<Text {...props} />)

  return {
    props,
    wrapper
  }
}

describe('Text should', () => {
  test('render', () => {
    const { wrapper } = setup();

    expect(wrapper.find('CardText').length).toBe(2);
  })
})