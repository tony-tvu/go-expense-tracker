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
import EditAccountBtn from '../components/EditAccountBtn'
import AddAccountBtn from '../components/AddAccountBtn'
import logger from '../logger'
import { IoIosWarning } from 'react-icons/io'

export default function LinkedAccounts() {
  const [data, setData] = useState([])
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)

  const stackBgColor = useColorModeValue('white', 'gray.900')
  const tooltipBg = useColorModeValue('white', 'gray.900')
  const tooltipColor = useColorModeValue('black', 'white')

  useEffect(() => {
    if (loading) {
      fetch(`${process.env.REACT_APP_API_URL}/items/${page}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then(async (res) => {
          if (!res) return
          const resData = await res.json().catch((err) => logger(err))
          if (res.status === 200 && resData.items) {
            setData((curr) => [...curr, ...resData.items])
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
          logger('error getting items', err)
        })
    }
  }, [page, loading])

  function onSuccess() {
    setPage(1)
    setData([])
    setLoading(true)
  }

  function renderAccounts() {
    if (data.length === 0 && !loading) {
      return (
        <Text fontSize="l" pl={1}>
          You have not linked any accounts.
        </Text>
      )
    }

    return data.map((item) => {
      return (
        <HStack
          key={item.id}
          width={'100%'}
          borderWidth="1px"
          borderRadius="lg"
          bg={stackBgColor}
          boxShadow={'2xl'}
          p={3}
          mb={5}
        >
          <Text fontSize="xl" as="b">
            {item.institution}
          </Text>
          <Spacer />

          {(item.item_login_required) && (
            <>
              <Tooltip
                label={`This account is unable to connect to your financial instituion. Please remove this account and add it back to repair the connection.`}
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

          <EditAccountBtn item={item} onSuccess={onSuccess} />
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
