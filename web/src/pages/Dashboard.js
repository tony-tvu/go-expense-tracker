import React, { useEffect, useState } from 'react'
import logger from '../logger'
import TotalSquare from '../components/TotalSquare'
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { Box } from '@chakra-ui/react'
import MonthYearPicker from '../components/MonthYearPicker'
import { DateTime } from 'luxon'
import ExpensesTable from '../components/ExpensesTable'

export default function Dashboard() {
  const [selectedMonth, setSelectedMonth] = useState(DateTime.now().month)
  const [selectedYear, setSelectedYear] = useState(DateTime.now().year)

  const [accountsData, setAccountsData] = useState([])
  const [transactionsData, setTransactionsData] = useState(null)
  const [cashTotal, setCashTotal] = useState(null)
  const [creditTotal, setCreditTotal] = useState(null)
  const [transactionsTotal, setTransactionsTotal] = useState(null)
  const [loading, setLoading] = useState(true)

  async function fetchAccounts() {
    await fetch(`${process.env.REACT_APP_API_URL}/accounts`, {
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
          resData.accounts.forEach((account) => {
            if (account.subtype === 'credit_card') {
              creditTotal += account.balance
            } else {
              cashTotal += account.balance
            }
          })
          setCashTotal(cashTotal)
          creditTotal = -1 * creditTotal
          setCreditTotal(creditTotal)
          setLoading(false)
        }
      })
      .catch((err) => {
        logger('error fetching accounts', err)
      })
  }

  async function fetchTransactions() {
    await fetch(`${process.env.REACT_APP_API_URL}/transactions`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    })
      .then(async (res) => {
        if (!res) return
        const resData = await res.json().catch((err) => logger(err))
        if (res.status === 200 && resData.transactions) {
          setTransactionsData(resData.transactions)
          let total = 0
          resData.transactions.forEach((transaction) => {
            total += transaction.amount
          })
          setTransactionsTotal(total)
        }
      })
      .catch((err) => {
        logger('error fetching accounts', err)
      })
  }

  useEffect(() => {
    document.title = 'Dashboard'
    if (loading) {
      fetchAccounts()
      fetchTransactions()
    }
  }, [accountsData, loading])

  console.log(selectedMonth)
  console.log(selectedYear)

  return (
    <Container style={{ paddingLeft: 0, paddingRight: 0 }}>
      <Row style={{ paddingLeft: 0, paddingRight: 0 }}>
        {/* Main Col 1 - Totals and Monthly overview */}
        <Col xs={12} sm={12} md={12}>
          {/* Totals */}
          <Row>
            <Col xs={4} sm={4} md={3}>
              <TotalSquare total={cashTotal} title={'Income'} />
            </Col>
            <Col xs={4} sm={4} md={3}>
              <TotalSquare total={transactionsTotal} title={'Expenses'} />
            </Col>
            <Col xs={4} sm={4} md={3}>
              <TotalSquare total={transactionsTotal} title={'Profit'} />
            </Col>
            <Col xs={12} sm={12} md={3}>
              <MonthYearPicker
                selectedMonth={selectedMonth}
                setSelectedMonth={setSelectedMonth}
                selectedYear={selectedYear}
                setSelectedYear={setSelectedYear}
                transactionsData={transactionsData ?? null}
              />
            </Col>
          </Row>
        </Col>
      </Row>
      <Row>
        <Col xs={12} sm={12} md={12}>
          <ExpensesTable transactionsData={transactionsData ?? null} />
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
