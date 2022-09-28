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
import AccountsSummary from '../components/AccountsSummary'
import Transactions from '../components/Transactions'

export default function OverviewPage() {
  return (
    <VStack>
      <Container maxW="container.md" centerContent>
        <AccountsSummary />
      </Container>
    </VStack>
  )
}
