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
          logger('error getting items', err)
        })
    }
  }, [loading])

  function renderAccounts() {
    if (data.length === 0 && !loading) {
      return (
        <Text fontSize="l" pl={1}>
          You have not linked any accounts.
        </Text>
      )
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
      <Container maxW="container.md" centerContent>
        <HStack alignItems="end" width="100%" mt={5} mb={5}>
          <Text fontSize="3xl" as="b" pl={1}>
            Accounts
          </Text>
          <Spacer />
          <AddAccountBtn onSuccess={onSuccess} />
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
          renderAccounts()
        )}
      </Container>
    </VStack>
  )
}
