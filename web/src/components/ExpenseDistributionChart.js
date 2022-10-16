import React from 'react'
import {
  BarChart,
  Bar,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import logger from '../logger'

export default function ExpenseDistributionChart({ transactionsData }) {
  if (!transactionsData) return null

  function calculateExpenseDistribution() {
    let bills = 0
    let entertainment = 0
    let groceries = 0
    let restaurant = 0
    let transportation = 0
    let vacation = 0
    let uncategorized = 0

    transactionsData.forEach((transaction) => {
      switch (transaction.category) {
        case 'bills':
          bills += transaction.amount
          break
        case 'entertainment':
          entertainment += transaction.amount
          break
        case 'groceries':
          groceries += transaction.amount
          break
        case 'restaurant':
          restaurant += transaction.amount
          break
        case 'transportation':
          transportation += transaction.amount
          break
        case 'vacation':
          vacation += transaction.amount
          break
        case 'uncategorized':
          uncategorized += transaction.amount
          break
        default:
          logger('unknown expense category')
      }
    })

    return [
      {
        name: 'Bills',
        total: -1 * bills,
        color: '#004CA3',
      },
      {
        name: 'Entertainment',
        total: -1 * entertainment,
        color: '#8A51A5',
      },
      {
        name: 'Groceries',
        total: -1 * groceries,
        color: '#CB5E99',
      },
      {
        name: 'Restaurant',
        total: -1 * restaurant,
        color: '#F47B89',
      },
      {
        name: 'Transportation',
        total: -1 * transportation,
        color: '#FFA47E',
      },
      {
        name: 'Vacation',
        total: -1 * vacation,
        color: '#FFD286',
      },
      {
        name: 'Uncategorized',
        total: -1 * uncategorized,
        color: '#FFFFA6',
      },
    ]
  }

  return (
    <ResponsiveContainer width="100%" height="100%">
      <BarChart
        width={500}
        height={300}
        data={calculateExpenseDistribution()}
        margin={{
          top: 5,
          right: 30,
          left: 20,
          bottom: 5,
        }}
      >
        <XAxis dataKey="name" />
        <YAxis />
        <Tooltip />
        <Bar dataKey="total">
          {calculateExpenseDistribution().map((entry, index) => (
            <Cell fill={calculateExpenseDistribution()[index].color} />
          ))}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  )
}
