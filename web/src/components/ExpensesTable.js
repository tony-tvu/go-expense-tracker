import {
  Box,
  Center,
  Divider,
  Select,
  Spinner,
  Text,
  useColorModeValue,
} from '@chakra-ui/react'
import React from 'react'
import { currency } from '../util'
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { DateTime } from 'luxon'

export default function ExpensesTable({ transactionsData }) {
  const selectorBg = useColorModeValue('gray.100', '#1E1E1E')

  function renderRows() {
    return transactionsData.map((transaction) => {
      return (
        <Box mb={2} borderColor={'#464646'}>
          <Row key={transaction.id}>
            <Col
              xs={3}
              sm={3}
              md={1}
              className="d-flex align-items-center justify-content-center"
            >
              <Text alignItems={'center'}>
                {DateTime.fromISO(transaction.date).toLocaleString()}
              </Text>
            </Col>
            <Col xs={3} sm={3} md={5} className="d-flex align-items-center">
              <Text>{transaction.name}</Text>
            </Col>
            <Col xs={3} sm={3} md={3} className="d-flex align-items-center">
              <Text>{currency.format(transaction.amount)}</Text>
            </Col>
            <Col xs={3} sm={3} md={3} className="d-flex align-items-center">
              <Select borderColor={selectorBg}>
                <option value={1}>January</option>
                <option value={2}>February</option>
              </Select>
            </Col>
          </Row>
          <Divider mt={2} />
        </Box>
      )
    })
  }

  if (!transactionsData) {
    return (
      <Center w={'100%'} minH={'200px'}>
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
    <Container style={{ paddingLeft: 0, paddingRight: 0 }}>
      <Box mb={2}>
        <Row>
          <Col xs={3} sm={3} md={1} className="d-flex align-items-center">
            <Text fontWeight={'600'} fontSize={'lg'}>
              Date
            </Text>
          </Col>
          <Col xs={3} sm={3} md={5} className="d-flex align-items-center">
            <Text fontWeight={'600'} fontSize={'lg'}>
              Name
            </Text>
          </Col>
          <Col xs={3} sm={3} md={3} className="d-flex align-items-center">
            <Text fontWeight={'600'} fontSize={'lg'}>
              Amount
            </Text>
          </Col>
          <Col xs={3} sm={3} md={3} className="d-flex align-items-center">
            <Text fontWeight={'600'} fontSize={'lg'}>
              Category
            </Text>
          </Col>
        </Row>
        <Divider borderColor={'#464646'} mt={3} />
      </Box>
      {renderRows()}
    </Container>
  )
}
