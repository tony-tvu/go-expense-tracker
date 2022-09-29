import React, { useEffect, useState } from 'react'
import logger from '../logger'
import TotalSquare from '../components/TotalSquare'
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { Box } from '@chakra-ui/react'

export default function Dashboard() {
  const [accountsData, setAccountsData] = useState([])
  const [cashTotal, setCashTotal] = useState(null)
  const [creditTotal, setCreditTotal] = useState(null)
  const [netWorth, setNetWorth] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (loading) {
      fetch(`${process.env.REACT_APP_API_URL}/accounts`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then(async (res) => {
          if (!res) return
          const resData = await res.json().catch((err) => logger(err))
          if (res.status === 200 && resData.accounts) {
            setAccountsData(resData.accounts)
            let cashTotal = 0
            let creditTotal = 0
            resData.accounts.forEach((acc) => {
              if (acc.type === 'checking' || acc.type === 'savings') {
                cashTotal += acc.current_balance
              } else if (acc.type === 'credit card') {
                creditTotal += acc.current_balance
              }
            })
            setNetWorth(cashTotal - creditTotal)

            setCashTotal(cashTotal)
            creditTotal = -1 * creditTotal
            setCreditTotal(creditTotal)
            setLoading(false)
          }
        })
        .catch((err) => {
          logger('error getting transactions', err)
        })
    }
  }, [accountsData, loading])

  return (
    <Container style={{ paddingLeft: 0, paddingRight: 0 }}>
      <Row style={{ paddingLeft: 0, paddingRight: 0 }}>
        {/* Main Col 1 - Totals and Monthly overview */}
        <Col xs={12} sm={12} md={8}>
          {/* Totals */}
          <Row>
            <Col xs={4} sm={4} md={4}>
              <TotalSquare total={cashTotal} title={'Total balance'} />
            </Col>
            <Col xs={4} sm={4} md={4}>
              <TotalSquare total={cashTotal} title={'Total balance'} />
            </Col>
            <Col xs={4} sm={4} md={4}>
              <TotalSquare total={cashTotal} title={'Total balance'} />
            </Col>
          </Row>

          {/* Monthly overview */}
          <Row>
            <Col xs={12} sm={12} md={12}>
              <Box bg={'red'} h={'300px'}></Box>
            </Col>
          </Row>
        </Col>

        {/* Main Col 2 - Activity */}
        <Col xs={12} sm={12} md={4}>
          <Box bg={'green'} minH={['100px', '100px', '130px', '650px']}></Box>
        </Col>
      </Row>
    </Container>
  )
}

{
  /* <Container>
<Row>
  <Col sm={12} md={6}>
    <AccountSummary />
  </Col>
  <Col sm={12} md={6}>
    <ExpenseSummary />
  </Col>
</Row>
</Container> */
}
