import {
  Center,
  Container,
  HStack,
  Spinner,
  Text,
  VStack,
} from '@chakra-ui/react'
import React, { useCallback, useEffect, useState } from 'react'
import logger from '../logger'

export default function Transactions() {
  const [data, setData] = useState([])
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [expensesTotal, setExpensesTotal] = useState(null)

  useEffect(() => {
    if (loading) {
      fetch(`${process.env.REACT_APP_API_URL}/transactions/${page}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then(async (res) => {
          if (!res) return
          const resData = await res.json().catch((err) => logger(err))
          if (res.status === 200 && resData.transactions) {
            setData((curr) => [...curr, ...resData.transactions])
            if (page !== Number(resData.page_info.totalPage)) {
              setPage(resData.page_info.next)
            } else {
              setLoading(false)
            }
          } else {
            setLoading(false)
          }
        })
        .catch((err) => {
          logger('error getting transactions', err)
        })
    }
  }, [page, loading])

  function show() {
    if (data.length === 0 && !loading) {
      return (
        <Text fontSize="l" pl={1}>
          You have no transactions
        </Text>
      )
    } else {
      
      return (
        <Text fontSize="l" pl={1}>
          {calculateTotals()}
        </Text>
      )
    }
  }

  function calculateTotals() {
    let expTotal = 0

    data.forEach((transaction) => {
      expTotal += transaction.amount
    })

    return expTotal
  }

  return (
    <VStack>
      <Container maxW="container.md" centerContent>
        <HStack
          alignItems="end"
          justifyContent={'center'}
          width="100%"
          mt={5}
          mb={5}
        >
          <Text fontSize="3xl" as="b" pl={1}>
            Overview
          </Text>
        </HStack>

        {loading ? (
          <Center pt={10}>
            <Spinner
              thickness="4px"
              speed="0.65s"
              emptyColor="gray.200"
              color="blue.500"
              size="xl"
            />
          </Center>
        ) : (
          show()
        )}
      </Container>
    </VStack>
  )
}
