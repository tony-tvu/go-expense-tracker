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
import AccountSummary from '../components/AccountSummary'
import Transactions from '../components/TransactionSummary'

export default function OverviewPage() {
  return (
    <VStack>
      <Container maxW="container.md" centerContent>
        <AccountSummary />
      </Container>
    </VStack>
  )
}
