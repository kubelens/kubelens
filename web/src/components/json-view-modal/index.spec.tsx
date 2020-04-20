import React from 'react'
import Enzyme, { shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import JsonViewModal from '.'

Enzyme.configure({ adapter: new Adapter() })

const setup = () => {
  const props = {
    handleClose: jest.fn(),
    show: true,
    body: {
      some: {
        nested: 'json'
      }
    },
    title: 'modal title',
    className: 'text-center'
  }

  const wrapper = shallow(<JsonViewModal {...props} />)

  return {
    props,
    wrapper
  }
}

describe('JsonViewModal should', () => {
  test('render open', () => {
    const { wrapper } = setup();

    expect(wrapper.find('Modal').hasClass('text-center')).toBe(true);
    expect(wrapper.find('ModalHeader').findWhere(n => n.text() === 'modal title')).toBeDefined();
    expect(wrapper.find('JsonView')).toBeDefined();
    expect(wrapper.find('Button').findWhere(n => n.text() === 'Close')).toBeDefined();
  })
})