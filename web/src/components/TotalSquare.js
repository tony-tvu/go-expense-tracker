import React from 'react'
import {
  Box,
  Text,
  VStack,
  useColorModeValue,
  Center,
  Spinner,
} from '@chakra-ui/react'
import { currency } from '../util'

export default function TotalSquare({ total, title }) {
  const bgColor = useColorModeValue('white', '#252526')
  const textColor = useColorModeValue('black', '#DCDCE2')

  if (!total && total !== 0) {
    return (
      <Center
        w={'100%'}
        minH={['50px', '50px', '145px', '145px']}
        bg={bgColor}
        borderRadius="7"
      >
        <Spinner
          thickness="4px"
          speed="0.65s"
          emptyColor="gray.200"
          color="blue.500"
          size="xl"
        />
      </Center>
    )
  }

  return (
    <Box
      bg={bgColor}
      w={'100%'}
      minH={['50px', '50px', '145px', '145px']}
      mb={5}
      borderRadius="7"
    >
      <VStack alignItems={'start'} p={5}>
        <Text
          pl={'2px'}
          fontSize={{
            base: '8px',
            sm: '10px',
            md: '18px',
            lg: '18px',
          }}
          color={textColor}
        >
          {title}
        </Text>
        <Text
          fontSize={{
            base: '8px',
            sm: '24px',
            md: '30px',
            lg: '36px',
          }}
          fontWeight={500}
        >
          {currency.format(total)}
        </Text>
      </VStack>
    </Box>
  )
}
