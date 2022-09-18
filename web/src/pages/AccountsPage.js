import React, { useCallback, useEffect, useState } from 'react'
import { useVerifyLogin } from '../hooks/useVerifyLogin'
import Navbar from '../components/Navbar'

import {
  Box,
  Stack,
  Grid,
  GridItem,
  Text,
  Spacer,
  Center,
  Spinner,
  HStack,
  useColorModeValue,
} from '@chakra-ui/react'
import EditAccountBtn from '../components/EditAccountBtn'
import AddAccountBtn from '../components/AddAccountBtn'
import logger from '../logger'

export default function Accounts() {
  useVerifyLogin()
  const [items, setItems] = useState(null)
  const [isEmpty, setIsEmpty] = useState(false)

  const stackBgColor = useColorModeValue('white', 'gray.900')

  const getItems = useCallback(async () => {
    await fetch(`${process.env.REACT_APP_API_URL}/items`, {
      method: 'GET',
      credentials: 'include',
    })
      .then(async (res) => {
        if (!res) return
        const data = await res.json().catch((err) => logger(err))
        if (data) {
          setItems(data)
          setIsEmpty(false)
          return
        }
        setItems([])
        setIsEmpty(true)
      })
      .catch((err) => {
        logger('error getting items', err)
      })
  }, [])

  useEffect(() => {
    getItems()
  }, [getItems])

  function renderItem(item) {
    return (
      <GridItem w="100%" key={item.id}>
        <HStack
          borderWidth="1px"
          borderRadius="lg"
          height={'150px'}
          bg={stackBgColor}
          boxShadow={'2xl'}
        >
          <Stack flex={1} alignItems="center">
            <Text fontSize="xl" as="b">
              {item.institution}
            </Text>
          </Stack>
          <Stack justifyContent="center" alignItems="center" p={5}>
            <EditAccountBtn item={item} onSuccess={getItems} />
          </Stack>
        </HStack>
      </GridItem>
    )
  }

  return (
    <>
      <Navbar />
      <Box pt={5} px={5} min={'100vh'}>
        <Stack direction={{ base: 'row', md: 'row' }} pb={5} alignItems="end">
          <Stack direction={{ base: 'row', md: 'row' }} alignItems="end">
            <Text fontSize="3xl" as="b" pl={1}>
              Accounts
            </Text>
          </Stack>
          <Spacer />
          <AddAccountBtn onSuccess={getItems} />
        </Stack>

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
        ) : isEmpty ? (
          <Text fontSize="l" pl={1}>
            You have not linked any accounts.
          </Text>
        ) : (
          <Grid templateColumns="repeat(2, 1fr)" gap={5}>
            {items.map((item) => {
              return renderItem(item)
            })}
          </Grid>
        )}
      </Box>
    </>
  )
}

{
  /* <Text fontSize="l" pl={1}>
You have not linked any accounts.
</Text> */
}
