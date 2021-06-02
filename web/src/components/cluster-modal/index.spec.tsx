import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import ErrorModal from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    handleClose: jest.fn(),
    open: true,
    status: 401,
    statusText: 'unauthorized',
    message: 'another message'
  }

  const wrapper = shallow(<ErrorModal {...props} />)

  return {
    props,
    wrapper
  }
}

describe('ErrorModal should', () => {
  test('render open', () => {
    const { wrapper } = setup();

    expect(wrapper.find('ModalHeader').findWhere(n => n.text() === '401 - unauthorized')).toBeDefined();
    expect(wrapper.find('ModalBody').findWhere(n => n.text() === 'another message')).toBeDefined();
    expect(wrapper.find('Button').findWhere(n => n.text() === 'Close').length).toBe(1);
  })
})