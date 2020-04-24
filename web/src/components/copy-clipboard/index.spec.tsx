import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import CopyClipboard from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    size: 12,
    value: 'some value',
    labelText: 'label',
    bottom: 6
  }

  const wrapper = shallow(<CopyClipboard {...props} />)

  return {
    props,
    wrapper
  }
}

describe('CopyClipboard should', () => {
  test('render', () => {
    const { wrapper } = setup();

    expect(wrapper.find('div').length).toBe(1);
    expect(wrapper.find('div').findWhere(n => n.text() === 'label').length).toBe(1);
    expect(wrapper.find('img').length).toBe(1);
  })
})