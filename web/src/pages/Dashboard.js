import React, { useCallback, useEffect, useState } from 'react'
import logger from '../logger'
import AccountSummary from '../components/AccountSummary'
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import ExpenseSummary from '../components/ExpenseSummary'

export default function Dashboard() {
  return (
    <Container>
      <Row>
        <Col sm={12} md={6}>
          <AccountSummary />
        </Col>
        <Col sm={12} md={6}>
          <ExpenseSummary />
        </Col>
      </Row>
    </Container>
  )
}
