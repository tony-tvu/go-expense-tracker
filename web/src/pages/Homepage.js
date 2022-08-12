import { useState } from 'react'
import {
  Box,
  Grid,
  GridItem,
  Flex,
  Text,
  Center,
  Square,
} from '@chakra-ui/react'
import GoogleLoginBtn from '../components/GoogleLoginBtn'

function Homepage() {
  return (
    <div>
      <GoogleLoginBtn />
      <Flex color="white" minH="85vh">
        <Center w="500px" bg="green.500">
          <Text>Box 1</Text>
        </Center>
      </Flex>
    </div>
  )
}

export default Homepage
