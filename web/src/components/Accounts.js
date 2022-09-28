import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Center,
  Container,
  Divider,
  HStack,
  Spacer,
  Spinner,
  Text,
  VStack,
} from '@chakra-ui/react'
import React, { useEffect, useState } from 'react'
import { currency, timeSince } from '../commons'
import logger from '../logger'

export default function CashAccounts() {
  const [data, setData] = useState([])
  const [cashTotal, setCashTotal] = useState(0)
  const [creditTotal, setCreditTotal] = useState(0)
  const [netWorth, setNetWorth] = useState(0)
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
            let cashTotal = 0
            let creditTotal = 0
            resData.accounts.forEach((acc) => {
              if (acc.type === 'checking' || acc.type === 'savings') {
                cashTotal += acc.current_balance
              } else if (acc.type === 'credit card') {
                creditTotal += acc.current_balance
              }
            })
            setNetWorth(cashTotal - creditTotal)

            setCashTotal(cashTotal)
            creditTotal = -1 * creditTotal
            setCreditTotal(creditTotal)
            setLoading(false)
          }
        })
        .catch((err) => {
          logger('error getting transactions', err)
        })
    }
  }, [data, loading])

  function renderAccounts(type) {
    return data
      .filter((acc) => acc.type === type)
      .map((acc) => {
        return (
          <AccordionPanel key={acc.id} pb={4}>
            <VStack>
              <HStack width="100%">
                <Text fontSize={['sm', 'md', 'lg', 'xl']}>{acc.name}</Text>
                <Spacer />
                <Text fontSize={['sm', 'md', 'lg', 'xl']}>
                  {currency.format(acc.current_balance)}
                </Text>
              </HStack>
              <HStack width="100%">
                <Text fontSize={['sm', 'md', 'lg', 'xl']} color={'gray.500'}>
                  {acc.institution}
                </Text>
                <Spacer />
                <Text fontSize={['sm', 'md', 'lg', 'xl']} color={'gray.500'}>
                  {timeSince(Date.parse(acc.updated_at))}
                </Text>
              </HStack>
            </VStack>
          </AccordionPanel>
        )
      })
  }

  return (
    <Container maxW="container.md" centerContent>
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
        <VStack width="100%" mt={5}>
          <HStack width="100%" m={0} p={0}>
            <Text
              mb={-3}
              fontSize={{
                base: '12px',
                sm: '12px',
                md: '16px',
                lg: '16px',
              }}
            >
              Net Worth
            </Text>
            <Spacer />
          </HStack>
          <HStack width="100%">
            <Text
              fontSize={{
                base: '26px',
                sm: '26px',
                md: '30px',
                lg: '34px',
              }}
              as="b"
            >
              {currency.format(netWorth)}
            </Text>
            <Spacer />
          </HStack>
          <Accordion allowMultiple defaultIndex={[]} width="100%">
            <AccordionItem>
              <AccordionButton>
                <HStack flex="1" pt={3} pb={3}>
                  <Text fontSize={['sm', 'md', 'lg', 'xl']} as="b">
                    Checking & Savings
                  </Text>
                  <Spacer />
                  <Text
                    fontSize={{
                      base: '20px',
                      sm: '20px',
                      md: '22px',
                      lg: '24px',
                    }}
                    as="b"
                    pr={5}
                  >
                    {currency.format(cashTotal)}
                  </Text>
                </HStack>
                <AccordionIcon />
              </AccordionButton>
              <Divider />
              {renderAccounts('checking')}
              {renderAccounts('savings')}
            </AccordionItem>
            <AccordionItem>
              <AccordionButton>
                <HStack flex="1" pt={3} pb={3}>
                  <Text fontSize={['sm', 'md', 'lg', 'xl']} as="b">
                    Credit Cards
                  </Text>
                  <Spacer />
                  <Text
                    fontSize={{
                      base: '20px',
                      sm: '20px',
                      md: '22px',
                      lg: '24px',
                    }}
                    as="b"
                    pr={5}
                  >
                    {currency.format(creditTotal)}
                  </Text>
                </HStack>
                <AccordionIcon />
              </AccordionButton>
              <Divider />
              {renderAccounts('credit card')}
            </AccordionItem>
          </Accordion>
        </VStack>
      )}
    </Container>
  )
}
