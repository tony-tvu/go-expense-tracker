import React, { useState } from 'react'
import {
  Box,
  FormControl,
  FormHelperText,
  FormLabel,
  Input,
  Select,
  useColorModeValue,
  VStack,
} from '@chakra-ui/react'
import { DateTime } from 'luxon'

export default function MonthYearPicker({
  selectedMonth,
  setSelectedMonth,
  selectedYear,
  setSelectedYear,
}) {
  const [loading, setLoading] = useState('')

  const bgColor = useColorModeValue('white', '#252526')
  const textColor = useColorModeValue('black', '#DCDCE2')

  function renderYearSelection() {
    const years = ['2022', '2021', '2020']

    return (
      <Select
        defaultValue={selectedYear}
        mb={3}
        onChange={(event) => setSelectedYear(event.target.value)}
      >
        {years.map((year) => {
          return <option value={year}>{year}</option>
        })}
      </Select>
    )
  }

  return (
    <Box
      bg={bgColor}
      w={'100%'}
      minH={['90px', '120px', '120px', '130px']}
      mb={5}
    >
      <FormControl p={5}>
        {renderYearSelection()}
        <Select
          defaultValue={selectedMonth}
          onChange={(event) => setSelectedMonth(event.target.value)}
        >
          <option value="1">January</option>
          <option value="2">February</option>
          <option value="3">March</option>
          <option value="4">April</option>
          <option value="5">May</option>
          <option value="6">June</option>
          <option value="7">July</option>
          <option value="8">August</option>
          <option value="9">September</option>
          <option value="10">October</option>
          <option value="11">November</option>
          <option value="12">December</option>
        </Select>
      </FormControl>
    </Box>
  )
}
