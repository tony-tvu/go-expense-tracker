import {
  Center,
  Container,
  HStack,
  Spinner,
  Text,
  VStack,
} from '@chakra-ui/react'
import React, { useCallback, useEffect, useState } from 'react'
import logger from '../logger'
import CashAccounts from "../components/CashAccounts"
import Transactions from "../components/Transactions"

export default function OverviewPage() {

  return (
    <VStack>
      <Container maxW="container.md" centerContent>
        <CashAccounts />
      
        <Transactions />
      </Container>
    </VStack>
  )
}
