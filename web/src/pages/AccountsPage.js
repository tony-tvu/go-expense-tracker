import React, { useCallback, useEffect, useState } from 'react'

import {
  Text,
  Spacer,
  VStack,
  Center,
  HStack,
  Spinner,
  Container,
  useColorModeValue,
} from '@chakra-ui/react'
import EditAccountBtn from '../components/EditAccountBtn'
import AddAccountBtn from '../components/AddAccountBtn'
import logger from '../logger'

export default function Accounts() {
  const [items, setItems] = useState(null)

  const stackBgColor = useColorModeValue('white', 'gray.900')

  const getItems = useCallback(async () => {
    await fetch(`${process.env.REACT_APP_API_URL}/items`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    })
      .then(async (res) => {
        if (!res) return
        const data = await res.json().catch((err) => logger(err))
        if (res.status === 200 && data.items) {
          setItems(data.items)
          return
        }
        setItems([])
      })
      .catch((err) => {
        logger('error getting items', err)
      })
  }, [])

  useEffect(() => {
    getItems()
  }, [getItems])

  function renderItems() {
    if (items.length === 0) {
      return (
        <Text fontSize="l" pl={1}>
          You have not linked any accounts.
        </Text>
      )
    }

    return items.map((item) => {
      return (
        <HStack
          width={'100%'}
          key={item.id}
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
          <EditAccountBtn item={item} onSuccess={getItems} />
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
          <AddAccountBtn onSuccess={getItems} />
        </HStack>

        {!items ? (
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
          renderItems()
        )}
      </Container>
    </VStack>
  )
}
