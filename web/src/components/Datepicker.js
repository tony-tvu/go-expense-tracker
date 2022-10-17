import React, { forwardRef } from 'react'
import ReactDatePicker from 'react-datepicker'
import { useColorModeValue } from '@chakra-ui/react'
import { InputGroup, Input, InputRightElement } from '@chakra-ui/react'
import { CalendarIcon } from '@chakra-ui/icons'

const customDateInput = ({ value, onClick, onChange }, ref) => (
  <Input
    autoComplete="off"
    value={value}
    ref={ref}
    onClick={onClick}
    onChange={onChange}
  />
)
customDateInput.displayName = 'DateInput'
const CustomInput = forwardRef(customDateInput)

const DatePicker = ({ selectedDate, onChange, ...props }) => {
  const bgColor = useColorModeValue('white', '#252526')

  return (
    <>
      <InputGroup className={useColorModeValue('light-theme', 'dark-theme')}>
        <ReactDatePicker
          selected={selectedDate}
          onChange={onChange}
          className="react-datapicker__input-text"
          customInput={<CustomInput bg={bgColor} />}
          dateFormat="MM/dd/yyyy"
          {...props}
        />
        <InputRightElement
          color="gray.500"
          children={<CalendarIcon fontSize="sm" />}
        />
      </InputGroup>
    </>
  )
}

export default DatePicker
