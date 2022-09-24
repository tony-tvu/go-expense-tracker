import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  Center,
  Container,
  HStack,
  Spacer,
  Spinner,
  Text,
  VStack,
} from '@chakra-ui/react'
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

  function getCashTotal() {
    let total = 0

    data.forEach((acc) => {
      total += acc.current_balance
    })

    return total
  }

  return (
    <Container maxW="container.md" centerContent >
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
        <Accordion allowMultiple width="100%">
          <AccordionItem>
            <h2>
              <AccordionButton>

                <HStack flex="1">
                  <Text fontSize='2xl' as='b'>Checking & Savings</Text>
                  <Spacer/>
                  <Text fontSize='2xl' as='b' pr={5}>${getCashTotal()}</Text>
                </HStack>
                <AccordionIcon />

              </AccordionButton>
            </h2>
            <AccordionPanel pb={4}>
              Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
              eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut
              enim ad minim veniam, quis nostrud exercitation ullamco laboris
              nisi ut aliquip ex ea commodo consequat.
            </AccordionPanel>
          </AccordionItem>
        </Accordion>
      )}
    </Container>
  )
}
