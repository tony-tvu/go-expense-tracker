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
import Accounts from '../components/AccountSummary'
import Transactions from '../components/TransactionSummary'

export default function TransactionsPage() {
  return (
    <VStack>
      <Container maxW="container.md" centerContent>
        <Transactions />
      </Container>
    </VStack>
  )
}
