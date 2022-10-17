import {
  Box,
  Divider,
  HStack,
  Select,
  Spacer,
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
import { FaCircle, FaSitemap, FaDollarSign } from 'react-icons/fa'
import AddTransactionBtn from './AddTransactionBtn'
import CreateRuleBtn from './CreateRuleBtn'
import EditTransactionBtn from './EditTransactionBtn'
import DeleteTransactionBtn from './DeleteTransactionBtn'

export default function TransactionsTable({
  transactionsData,
  onSuccess,
  forceRefresh,
}) {
  const bgColor = useColorModeValue('white', '#252526')
  const navigate = useNavigate()

  async function updateCategory(transactionId, category) {
    await fetch(`${process.env.REACT_APP_API_URL}/transactions/category`, {
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
        return '#004CA3'
      case 'entertainment':
        return '#8A51A5'
      case 'groceries':
        return '#CB5E99'
      case 'ignore':
        return 'grey'
      case 'income':
        return 'green'
      case 'restaurant':
        return '#F47B89'
      case 'transportation':
        return '#FFA47E'
      case 'vacation':
        return '#FFD286'
      case 'uncategorized':
        return '#FFFFA6'
      default:
        return 'black'
    }
  }

  function renderRows() {
    return transactionsData.map((transaction) => {
      return (
        <Box key={transaction.id} mb={2} borderColor={'#464646'}>
          <Row key={transaction.id}>
            <Col xs={3} sm={3} md={1} className="d-flex align-items-center">
              <Text alignItems={'center'}>
                {DateTime.fromISO(transaction.date, { zone: 'UTC' }).toFormat(
                  'LL/dd/yy'
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
                  borderColor={bgColor}
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
            <Col xs={3} sm={3} md={2} className="d-flex align-items-center">
              <Text>{currency.format(transaction.amount)}</Text>
            </Col>
            <Col xs={1} sm={1} md={1} className="d-flex align-items-center">
              <HStack>
                <EditTransactionBtn
                  onSuccess={onSuccess}
                  forceRefresh={forceRefresh}
                  transaction={transaction}
                />
                {transaction.enrollment_id === 'user_created' ? (
                  <DeleteTransactionBtn
                    onSuccess={onSuccess}
                    forceRefresh={forceRefresh}
                    transaction={transaction}
                  />
                ) : (
                  <></>
                )}
              </HStack>
            </Col>
          </Row>
          <Divider mt={2} />
        </Box>
      )
    })
  }

  return (
    <Container
      style={{
        padding: '20px',
        backgroundColor: bgColor,
        borderRadius: '10px',
      }}
    >
      <Box mb={2}>
        <HStack mb={5}>
          <Spacer />
          <CreateRuleBtn
            onSuccess={onSuccess}
            forceRefresh={forceRefresh}
            icon={<FaSitemap />}
          />
          <AddTransactionBtn onSuccess={onSuccess} icon={<FaDollarSign />} />
        </HStack>
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
          <Col xs={3} sm={3} md={2} className="d-flex align-items-center">
            <Text fontWeight={'600'} fontSize={'lg'}>
              Amount
            </Text>
          </Col>
          <Col xs={3} sm={3} md={1} className="d-flex align-items-center"></Col>
        </Row>
        <Divider borderColor={'#464646'} mt={3} />
      </Box>
      {renderRows()}
    </Container>
  )
}
