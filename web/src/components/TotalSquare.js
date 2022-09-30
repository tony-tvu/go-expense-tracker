import React, { useEffect, useState } from 'react'
import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  Container,
  Divider,
  HStack,
  Spacer,
  Text,
  VStack,
} from '@chakra-ui/react'
import { currency, timeSince } from '../commons'
import logger from '../logger'

export default function TotalSquare({ total, title }) {
  return (
    <Box bg={'pink'} w={'100%'} minH={['90px', '130px', '130px', '270px']} mb={5}>
      {title}
      {currency.format(total)}
    </Box>
  )
}
