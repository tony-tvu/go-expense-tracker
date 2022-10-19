import React, { useEffect, useState } from 'react'
import { IconButton, useColorModeValue } from '@chakra-ui/react'
import { CSVLink } from 'react-csv'
import { DateTime } from 'luxon'
import { FaDownload } from 'react-icons/fa'

export default function DownloadCsvBtn({ transactionsData }) {
  const [data, setData] = useState([])
  const bgColor = useColorModeValue('white', '#252526')

  useEffect(() => {
    if (transactionsData) {
      const filteredData = []
      transactionsData.forEach((t) => {
        if (t.category !== 'ignore') {
          filteredData.push({
            date: DateTime.fromISO(t.date, {
              zone: 'utc',
            }).toFormat('LL/dd/yyyy'),
            name: t.name,
            category: t.category,
            amount: t.amount < 0 ? -1 * t.amount : t.amount,
          })
        }
      })
      setData(filteredData)
    }
  }, [transactionsData])

  if (data.length === 0) return null

  return (
    <>
      <CSVLink
        data={data}
        filename={`${DateTime.now().year}_${DateTime.now().month}_${
          DateTime.now().day
        }_transactions.csv`}
      >
        <IconButton type="button" bg={bgColor} icon={<FaDownload />} />
      </CSVLink>
    </>
  )
}
