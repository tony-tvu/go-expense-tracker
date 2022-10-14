import React, { useEffect, useState } from 'react'

import {
  Text,
  Spacer,
  VStack,
  Center,
  HStack,
  Spinner,
  Container,
  useColorModeValue,
  Tooltip,
} from '@chakra-ui/react'
import DeleteAccountBtn from '../components/DeleteAccountBtn'
import AddAccountBtn from '../components/AddAccountBtn'
import { IoIosWarning } from 'react-icons/io'
import logger from '../logger'

export default function LinkedAccounts() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)

  const stackBgColor = useColorModeValue('white', 'gray.900')
  const bgColor = useColorModeValue('white', '#252526')
  const tooltipBg = useColorModeValue('white', 'gray.900')
  const tooltipColor = useColorModeValue('black', 'white')

  function onSuccess() {
    setData([])
    setLoading(true)
  }

  useEffect(() => {
    document.title = 'Accounts'
    if (loading) {
      fetch(`${process.env.REACT_APP_API_URL}/enrollments`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then(async (res) => {
          if (!res) return
          const resData = await res.json().catch((err) => logger(err))
          if (res.status === 200 && resData.enrollments) {
            setData(resData.enrollments)
          }
          setLoading(false)
        })
        .catch((err) => {
          logger('error getting enrollments', err)
        })
    }
  }, [loading])

  function renderAccounts() {
    if (data.length === 0 && !loading) {
      return null
    }

    return data.map((enrollment) => {
      return (
        <HStack
          key={enrollment.id}
          width={'100%'}
          borderWidth="1px"
          borderRadius="lg"
          bg={stackBgColor}
          boxShadow={'2xl'}
          p={3}
          mb={5}
        >
          <Text fontSize="xl" as="b">
            {enrollment.institution}
          </Text>
          <Spacer />

          {enrollment.disconnected && (
            <>
              <Tooltip
                label={`This account is unable to connect to your financial instituion. To resolve this issue, remove this account and add it again (your existing transactions will not be deleted)`}
                fontSize="md"
                bg={tooltipBg}
                color={tooltipColor}
                borderWidth="1px"
                boxShadow={'2xl'}
                borderRadius="lg"
                p={5}
              >
                <span>
                  <IoIosWarning size={'40px'} color={'red'} />
                </span>
              </Tooltip>
            </>
          )}

          <DeleteAccountBtn enrollment={enrollment} onSuccess={onSuccess} />
        </HStack>
      )
    })
  }

  return (
    <VStack>
      <Container maxW="container.md" centerContent bg={bgColor} mb={5}>
        <HStack alignItems="end" width="100%" mt={5} mb={5}>
          <Text fontSize="3xl" as="b" pl={1}>
            Accounts
          </Text>
          <Spacer />
          <AddAccountBtn onSuccess={onSuccess} />
        </HStack>
        <VStack alignItems="start" width="100%" mb={10}>
          <Text fontSize="md" pl={1} mb={5}>
            Add an account to begin tracking your expenses. This application
            pulls in new transactions every hour. Checking and savings account
            balances are refreshed every 12 hours.
          </Text>
          <HStack>
            <IoIosWarning size={'150px'} color={'red'} />
            <Text>
              When you delete an account, your existing transactions will not be
              deleted. This is to prevent you from needing to re-categorize the
              transactions you've already categorized. Whenever your account
              gets disconnected from this application, you'll need to remove
              that account and add it back again to resolve the issue. To remove
              existing transactions, you will need to completely wipe all
              transactions.
            </Text>
          </HStack>
        </VStack>
      </Container>

      <Container maxW="container.md" p={0} centerContent minH={'100px'}>
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
          renderAccounts()
        )}
      </Container>
    </VStack>
  )
}
