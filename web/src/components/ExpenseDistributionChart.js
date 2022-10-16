import React from 'react'
import {
  BarChart,
  Bar,
  Cell,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import logger from '../logger'

export default function ExpenseDistributionChart({ transactionsData }) {
  if (!transactionsData) return null

  function calculateExpenseDistribution() {
    const expenseMap = {
      bills: {
        name: 'Bills',
        total: 0,
        color: '#004CA3',
      },
      entertainment: {
        name: 'Entertainment',
        total: 0,
        color: '#8A51A5',
      },
      groceries: {
        name: 'Groceries',
        total: 0,
        color: '#CB5E99',
      },
      restaurant: {
        name: 'Restaurant',
        total: 0,
        color: '#F47B89',
      },
      transportation: {
        name: 'Transportation',
        total: 0,
        color: '#FFA47E',
      },
      vacation: {
        name: 'Vacation',
        total: 0,
        color: '#FFD286',
      },
      uncategorized: {
        name: 'Uncategorized',
        total: 0,
        color: '#FFFFA6',
      },
    }

    transactionsData.forEach((transaction) => {
      switch (transaction.category) {
        case 'bills':
          expenseMap['bills'].total += transaction.amount
          break
        case 'entertainment':
          expenseMap['entertainment'].total += transaction.amount
          break
        case 'groceries':
          expenseMap['groceries'].total += transaction.amount
          break
        case 'income':
          break
        case 'ignore':
          break
        case 'restaurant':
          expenseMap['restaurant'].total += transaction.amount
          break
        case 'transportation':
          expenseMap['transportation'].total += transaction.amount
          break
        case 'vacation':
          expenseMap['vacation'].total += transaction.amount
          break
        case 'uncategorized':
          expenseMap['uncategorized'].total += transaction.amount
          break
        default:
          logger('unknown expense category: ', transaction.category)
      }
    })

    let graphData = []

    Object.keys(expenseMap).forEach((key) => {
      expenseMap[key].total = -1 * expenseMap[key].total
    })

    Object.keys(expenseMap)
      .filter((key) => expenseMap[key].total !== 0)
      .map((key) => {
        return graphData.push(expenseMap[key])
      })

    return graphData
  }

  // TODO: set up state/useEffect to avoid calling calculateExpenseDistribution 3 times
  return (
    <ResponsiveContainer width="100%" height="100%">
      <BarChart width={500} height={300} data={calculateExpenseDistribution()}>
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
