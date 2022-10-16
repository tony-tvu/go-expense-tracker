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

export default function MonthYearPicker({ availableYears }) {
  const [appState, setAppState] = useContext(AppStateContext)
  const bgColor = useColorModeValue('white', '#252526')

  function renderYearSelection() {
    let years = availableYears.sort()
    if (years.length === 0) years = [DateTime.now().year]
    return (
      <Select
        defaultValue={appState.selectedYear}
        mb={3}
        onChange={(event) =>
          setAppState({
            selectedMonth: appState.selectedMonth,
            selectedYear: Number(event.target.value),
          })
        }
      >
        {years.map((year) => {
          return <option value={year}>{year}</option>
        })}
      </Select>
    )
  }

  function renderMonthSelection() {
    const monthMap = {
      1: 'January',
      2: 'February',
      3: 'March',
      4: 'April',
      5: 'May',
      6: 'June',
      7: 'July',
      8: 'August',
      9: 'September',
      10: 'October',
      11: 'November',
      12: 'December',
    }

    return (
      <Select
        defaultValue={appState.selectedMonth}
        mb={3}
        onChange={(event) => {
          setAppState({
            selectedMonth: Number(event.target.value),
            selectedYear: appState.selectedYear,
          })
        }}
      >
        {appState.selectedYear === DateTime.now().year
          ? Object.keys(monthMap)
              .filter((key) => key <= DateTime.now().month)
              .map((key) => {
                return <option value={key}>{monthMap[key]}</option>
              })
          : Object.keys(monthMap).map((key) => {
              return <option value={key}>{monthMap[key]}</option>
            })}
      </Select>
    )
  }

  if (!availableYears) {
    return (
      <Center
        w={'100%'}
        minH={['50px', '50px', '145px', '145px']}
        bg={bgColor}
        borderRadius="7"
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
      minH={['50px', '50px', '145px', '145px']}
      mb={5}
      borderRadius="7"
    >
      <FormControl p={5}>
        {renderYearSelection()}
        {renderMonthSelection()}
      </FormControl>
    </Box>
  )
}
