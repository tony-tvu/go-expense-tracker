import { Center, Spinner } from '@chakra-ui/react'
import React, { useEffect, useState } from 'react'
import logger from '../logger'

export default function CashAccounts() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (loading) {
      fetch(`${process.env.REACT_APP_API_URL}/cash_accounts`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then(async (res) => {
          if (!res) return
          const resData = await res.json().catch((err) => logger(err))
          if (res.status === 200 && resData.accounts) {
            setData(resData.accounts)
            setLoading(false)
          }
        })
        .catch((err) => {
          logger('error getting transactions', err)
        })
    }
  }, [loading])

  if (loading && !data) {
    return (
      <Center pt={10}>
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

  return <div>CashAccounts</div>
}
