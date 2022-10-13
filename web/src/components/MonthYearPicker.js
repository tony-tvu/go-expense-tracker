import React from 'react'
import {
  Box,
  Center,
  FormControl,
  Select,
  Spinner,
  useColorModeValue,
} from '@chakra-ui/react'
import { DateTime } from 'luxon'

export default function MonthYearPicker({
  selectedMonth,
  setSelectedMonth,
  selectedYear,
  setSelectedYear,
  transactionsData,
}) {
  const bgColor = useColorModeValue('white', '#252526')

  function getYears(transactions) {
    const years = [selectedYear]
    transactions.forEach((transaction) => {
      const year = DateTime.fromISO(transaction.date).year
      if (!years.includes(year)) {
        years.push(year)
      }
    })

    years.sort().reverse()
    return years
  }

  function renderYearSelection() {
    const years = getYears(transactionsData)
    return (
      <Select
        defaultValue={selectedYear}
        mb={3}
        onChange={(event) => setSelectedYear(Number(event.target.value))}
      >
        {years.map((year) => {
          return <option value={year}>{year}</option>
        })}
      </Select>
    )
  }

  if (!transactionsData) {
    return (
      <Center
        w={'100%'}
        minH={['90px', '120px', '120px', '130px']}
        bg={bgColor}
      >
        <Spinner
          thickness="4px"
          speed="0.65s"
          emptyColor="gray.200"
          color="blue.500"
          size="xl"
        />
      </Center>
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
          onChange={(event) => setSelectedMonth(Number(event.target.value))}
        >
          <option value={1}>January</option>
          <option value={2}>February</option>
          <option value={3}>March</option>
          <option value={4}>April</option>
          <option value={5}>May</option>
          <option value={6}>June</option>
          <option value={7}>July</option>
          <option value={8}>August</option>
          <option value={9}>September</option>
          <option value={10}>October</option>
          <option value={11}>November</option>
          <option value={12}>December</option>
        </Select>
      </FormControl>
    </Box>
  )
}
