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

export default function OverviewPage() {

  return (
    <VStack>
      <Container maxW="container.md" centerContent>
        <HStack
          alignItems="end"
          justifyContent={'center'}
          width="100%"
          mt={5}
          mb={5}
        >
          <Text fontSize="3xl" as="b" pl={1}>
            Overview
          </Text>
        </HStack>
        <CashAccounts />
      

      </Container>
    </VStack>
  )
}
