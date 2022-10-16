import React, { useContext } from 'react'
import {
  Box,
  Center,
  FormControl,
  Select,
  Spinner,
  useColorModeValue,
} from '@chakra-ui/react'
import { DateTime } from 'luxon'
import { AppStateContext } from '../hooks/AppStateProvider'

export default function MonthYearPicker({
  selectedMonth,
  selectedYear,
  availableYears,
}) {
  const [appState, setAppState] = useContext(AppStateContext)
  const bgColor = useColorModeValue('white', '#252526')

  function renderYearSelection() {
    let years = availableYears.sort().reverse()
    if (years.length === 0) years = [DateTime.now().year]
    return (
      <Select
        defaultValue={selectedYear}
        mb={3}
        onChange={(event) => setAppState({
          selectedMonth: appState.selectedMonth,
          selectedYear: Number(event.target.value)
        })}
      >
        {years.map((year) => {
          return <option value={year}>{year}</option>
        })}
      </Select>
    )
  }

  if (!availableYears) {
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
          onChange={(event) => setAppState({
            selectedMonth: Number(event.target.value),
            selectedYear: appState.selectedYear
          })}
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
