import React, { useEffect } from "react"
import { useVerifyLogin } from "../hooks/useVerifyLogin"
import Navbar from "../components/Navbar"
import { useLazyQuery, gql } from "@apollo/client"

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
} from "@chakra-ui/react"
import EditAccountBtn from "../components/EditAccountBtn"
import AddAccountBtn from "../components/AddAccountBtn"

const query = gql`
  query {
    items {
      id
      userID
      institution
      createdAt
      updatedAt
    }
  }
`

export default function Accounts() {
  useVerifyLogin()

  const [getItems, { data, loading }] = useLazyQuery(query, {
    fetchPolicy: "no-cache",
  })

  const stackBgColor = useColorModeValue("white", "gray.900")

  useEffect(() => {
    getItems()
  }, [getItems])

  function renderItem(item) {
    return (
      <GridItem w="100%" key={item.id}>
        <HStack
          borderWidth="1px"
          borderRadius="lg"
          height={"150px"}
          bg={stackBgColor}
          boxShadow={"2xl"}
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
      <Box pt={5} px={5} min={"100vh"}>
        <Stack direction={{ base: "row", md: "row" }} pb={5} alignItems="end">
          <Stack direction={{ base: "row", md: "row" }} alignItems="end">
            <Text fontSize="3xl" as="b" pl={1}>
              Accounts
            </Text>
          </Stack>
          <Spacer />
          <AddAccountBtn onSuccess={getItems} />
        </Stack>

        {!data || loading ? (
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
          <Grid templateColumns="repeat(2, 1fr)" gap={5}>
            {data.items.map((item) => {
              return renderItem(item)
            })}
          </Grid>
        )}
      </Box>
    </>
  )
}
