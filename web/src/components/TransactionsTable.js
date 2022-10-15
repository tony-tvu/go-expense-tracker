import {
  Box,
  Divider,
  HStack,
  Select,
  Text,
  useColorModeValue,
} from '@chakra-ui/react'
import React from 'react'
import { currency } from '../util'
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { DateTime } from 'luxon'
import { useNavigate } from 'react-router-dom'
import logger from '../logger'
import { FaCircle } from 'react-icons/fa'

export default function TransactionsTable({ transactionsData, onSuccess }) {
  const selectorBg = useColorModeValue('gray.100', '#1E1E1E')
  const navigate = useNavigate()

  async function updateCategory(transactionId, category) {
    await fetch(`${process.env.REACT_APP_API_URL}/transactions`, {
      method: 'PATCH',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        transaction_id: transactionId,
        category: category,
      }),
    })
      .then((res) => {
        if (res.status === 401) navigate('/login')
        if (res.status === 200) onSuccess()
      })
      .catch((e) => {
        logger('error updating transaction category', e)
      })
  }

  if (!transactionsData) {
    return null
  }

  function getCategoryColor(category) {
    switch (category) {
      case 'bills':
        return 'orange'
      case 'entertainment':
        return 'yellow'
      case 'groceries':
        return 'blue'
      case 'ignore':
        return 'grey'
      case 'income':
        return 'green'
      case 'restaurant':
        return 'purple'
      case 'transportation':
        return 'pink'
      case 'vacation':
        return 'brown'
      case 'uncategorized':
        return 'red'
      default:
        return 'black'
    }
  }

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
                {DateTime.fromISO(transaction.date, { zone: 'UTC' }).toFormat(
                  'LL/dd/yyyy'
                )}
              </Text>
            </Col>
            <Col xs={3} sm={3} md={5} className="d-flex align-items-center">
              <Text>{transaction.name}</Text>
            </Col>
            <Col xs={3} sm={3} md={3} className="d-flex align-items-center">
              <HStack>
                <FaCircle color={getCategoryColor(transaction.category)} />
                <Select
                  defaultValue={transaction.category}
                  borderColor={selectorBg}
                  onChange={async (event) => {
                    await updateCategory(
                      transaction.transaction_id,
                      event.target.value
                    )
                  }}
                >
                  <option value={'bills'}>Bills</option>
                  <option value={'entertainment'}>Entertainment</option>
                  <option value={'groceries'}>Groceries</option>
                  <option value={'ignore'}>Ignore</option>
                  <option value={'income'}>Income</option>
                  <option value={'restaurant'}>Restaurant</option>
                  <option value={'transportation'}>Transportation</option>
                  <option value={'vacation'}>Vacation</option>
                  <option value={'uncategorized'}>Uncategorized</option>
                </Select>
              </HStack>
            </Col>
            <Col xs={3} sm={3} md={3} className="d-flex align-items-center">
              <Text>{currency.format(transaction.amount)}</Text>
            </Col>
          </Row>
          <Divider mt={2} />
        </Box>
      )
    })
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
              Category
            </Text>
          </Col>
          <Col xs={3} sm={3} md={3} className="d-flex align-items-center">
            <Text fontWeight={'600'} fontSize={'lg'}>
              Amount
            </Text>
          </Col>
        </Row>
        <Divider borderColor={'#464646'} mt={3} />
      </Box>
      {renderRows()}
    </Container>
  )
}