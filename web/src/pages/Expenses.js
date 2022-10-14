import React, { useCallback, useEffect, useState } from 'react'
import logger from '../logger'
import TotalSquare from '../components/TotalSquare'
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import MonthYearPicker from '../components/MonthYearPicker'
import { DateTime } from 'luxon'
import ExpensesTable from '../components/ExpensesTable'

export default function Expenses() {
  const [selectedMonth, setSelectedMonth] = useState(DateTime.now().month)
  const [selectedYear, setSelectedYear] = useState(DateTime.now().year)

  const [transactionsData, setTransactionsData] = useState(null)
  const [transactionsTotal, setTransactionsTotal] = useState(null)
  const [availableYears, setAvailableYears] = useState(null)
  const [loading, setLoading] = useState(true)

  const fetchTransactions = useCallback(async () => {
    setTransactionsData(null)
    setLoading(true)
    await fetch(
      `${process.env.REACT_APP_API_URL}/transactions?month=${selectedMonth}&year=${selectedYear}`,
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
        if (res.status === 200 && resData.transactions) {
          setTransactionsData(resData.transactions)
          setAvailableYears(resData.years)
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
  }, [selectedMonth, selectedYear])

  useEffect(() => {
    document.title = 'Expenses'
    if (loading) {
      fetchTransactions()
    }
  }, [fetchTransactions, loading])

  useEffect(() => {
    fetchTransactions()
  }, [fetchTransactions, selectedMonth, selectedYear])

  return (
    <Container style={{ paddingLeft: 0, paddingRight: 0 }}>
      <Row style={{ paddingLeft: 0, paddingRight: 0 }}>
        <Col xs={12} sm={12} md={12}>
          <Row>
            <Col xs={4} sm={4} md={3}>
              <TotalSquare total={transactionsTotal} title={'Expenses'} />
            </Col>
            <Col xs={4} sm={4} md={3}>
              <TotalSquare total={1} title={'Income'} />
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
                availableYears={availableYears ?? null}
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
