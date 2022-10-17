import React, { useCallback, useContext, useEffect, useState } from 'react'
import logger from '../logger'
import TotalSquare from '../components/TotalSquare'
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import MonthYearPicker from '../components/MonthYearPicker'
import TransactionsTable from '../components/TransactionsTable'
import { Flex } from '@chakra-ui/react'
import ExpenseDistributionChart from '../components/ExpenseDistributionChart'
import { useNavigate } from 'react-router-dom'
import { AppStateContext } from '../hooks/AppStateProvider'

export default function Transactions() {
  const [appState] = useContext(AppStateContext)
  const [transactionsData, setTransactionsData] = useState(null)
  const [expensesTotal, setExpensesTotal] = useState(null)
  const [incomeTotal, setIncomeTotal] = useState(null)
  const [profit, setProfit] = useState(null)
  const [availableYears, setAvailableYears] = useState(null)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  const fetchTransactions = useCallback(async () => {
    setLoading(true)
    await fetch(
      `${process.env.REACT_APP_API_URL}/transactions?month=${appState.selectedMonth}&year=${appState.selectedYear}`,
      {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      }
    )
      .then(async (res) => {
        if (!res) return
        const resData = await res.json().catch((err) => logger(err))
        if (res.status === 401) navigate('/login')
        if (res.status === 200 && resData.transactions) {
          setTransactionsData(resData.transactions)
          setAvailableYears(resData.years)

          let calculatedExpenses = 0
          let calculatedIncome = 0
          resData.transactions.forEach((transaction) => {
            if (
              transaction.category !== 'ignore' &&
              transaction.category !== 'income'
            ) {
              calculatedExpenses += transaction.amount
            } else if (
              transaction.category !== 'ignore' &&
              transaction.category === 'income'
            ) {
              calculatedIncome += transaction.amount
            }
          })
          setExpensesTotal(calculatedExpenses)
          setIncomeTotal(calculatedIncome)
          setProfit(calculatedIncome + calculatedExpenses)
        }
      })
      .catch((err) => {
        logger('error fetching accounts', err)
      })
  }, [navigate, appState])

  useEffect(() => {
    document.title = 'Transactions'
    fetchTransactions()
  }, [fetchTransactions, loading, appState])

  function forceRefresh() {
    setTransactionsData(null)
  }

  return (
    <Flex>
      <Container>
        <Row>
          <Col xs={4} sm={4} md={3}>
            <TotalSquare total={incomeTotal ?? null} title={'Income'} />
          </Col>
          <Col xs={4} sm={4} md={3}>
            <TotalSquare total={expensesTotal ?? null} title={'Expenses'} />
          </Col>
          <Col xs={4} sm={4} md={3}>
            <TotalSquare total={profit ?? null} title={'Profit'} />
          </Col>
          <Col xs={12} sm={12} md={3}>
            <MonthYearPicker availableYears={availableYears ?? null} />
          </Col>
        </Row>

        <Row style={{ height: '300px', marginBottom: '25px' }}>
          <Col xs={12} sm={12} md={12}>
            <ExpenseDistributionChart
              transactionsData={transactionsData ?? null}
            />
          </Col>
        </Row>

        <Row>
          <Col xs={12} sm={12} md={12}>
            <TransactionsTable
              transactionsData={transactionsData ?? null}
              onSuccess={fetchTransactions}
              forceRefresh={forceRefresh}
            />
          </Col>
        </Row>
      </Container>
    </Flex>
  )
}
