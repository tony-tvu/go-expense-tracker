import React, { useEffect, useState } from 'react'
import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  Center,
  Container,
  Divider,
  HStack,
  Spacer,
  Text,
  VStack,
  useColorModeValue,
} from '@chakra-ui/react'
import { currency, timeSince } from '../util'
import logger from '../logger'

export default function TotalSquare({ total, title }) {
  const bgColor = useColorModeValue('white', '#252526')
  const textColor = useColorModeValue('black', '#DCDCE2')

  return (
    <Box
      bg={bgColor}
      w={'100%'}
      minH={['50px', '50px', '132px', '132px']}
      mb={5}
    >
      <VStack alignItems={'start'} p={5}>
        <Text
        pl={'2px'}
          fontSize={{
            base: '14px',
            sm: '14px',
            md: '18px',
            lg: '18px',
          }}
          color={textColor}
        >
          {title}
        </Text>
        <Text
          fontSize={{
            base: '18px',
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
