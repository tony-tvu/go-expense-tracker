import React, { useEffect, useState, useContext } from 'react'
import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  chakra,
  Text,
  HStack,
  Select,
  useColorModeValue,
  Box,
  Spacer,
  FormControl,
  FormLabel,
  Input,
} from '@chakra-ui/react'
import { TriangleDownIcon, TriangleUpIcon } from '@chakra-ui/icons'
import {
  useReactTable,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  createColumnHelper,
} from '@tanstack/react-table'
import { useNavigate } from 'react-router-dom'
import { DateTime } from 'luxon'
import { FaCircle, FaSitemap, FaDollarSign } from 'react-icons/fa'
import logger from '../logger'
import EditTransactionBtn from './EditTransactionBtn'
import DeleteTransactionBtn from './DeleteTransactionBtn'
import { currency } from '../util'
import CreateRuleBtn from './CreateRuleBtn'
import AddTransactionBtn from './AddTransactionBtn'
import { AppStateContext } from '../hooks/AppStateProvider'
import DownloadCsvBtn from './DownloadCsvBtn'

export default function TransactionsTable({
  transactionsData,
  onSuccess,
  forceRefresh,
}) {
  const [appState] = useContext(AppStateContext)
  const [data, setData] = useState([])
  const [sorting, setSorting] = useState([])
  const [searchStr, setSearchStr] = useState('')
  const [filteredCategory, setFilteredCategory] = useState('all')
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

  useEffect(() => {
    setData([])
  }, [searchStr])

  useEffect(() => {
    setData([])
  }, [filteredCategory])

  useEffect(() => {
    setData([])
  }, [appState])

  useEffect(() => {
    if (transactionsData) {
      let transactions = []
      transactionsData.forEach((t) => {
        // apply searchStr filter
        if (t.name.toLowerCase().includes(searchStr.toLowerCase())) {
          if (filteredCategory === 'all') {
            transactions.push({
              date: t.date,
              name: t.name,
              category: t.category,
              amount: t.amount,
              transactionId: t.transaction_id,
              enrollmentId: t.enrollment_id,
              options: '',
            })
          } else {
            // apply category filter
            if (t.category === filteredCategory) {
              transactions.push({
                date: t.date,
                name: t.name,
                category: t.category,
                amount: t.amount,
                transactionId: t.transaction_id,
                enrollmentId: t.enrollment_id,
                options: '',
              })
            }
          }
        }
      })
      setData(transactions)
    }
  }, [transactionsData, searchStr, filteredCategory])

  const columnHelper = createColumnHelper()
  const columns = [
    columnHelper.accessor('date', {
      cell: (info) => info.getValue(),
      header: 'Date',
    }),
    columnHelper.accessor('name', {
      cell: (info) => info.getValue(),
      header: 'Name',
    }),
    columnHelper.accessor('category', {
      cell: (info) => info.getValue(),
      header: 'Category',
    }),
    columnHelper.accessor('amount', {
      cell: (info) => info.getValue(),
      header: 'Amount',
    }),
    columnHelper.accessor('options', {
      cell: (info) => info.getValue(),
      header: '',
    }),
  ]

  const table = useReactTable({
    columns,
    data,
    getCoreRowModel: getCoreRowModel(),
    onSortingChange: setSorting,
    getSortedRowModel: getSortedRowModel(),
    state: {
      sorting,
    },
  })

  if (data.length === 0 && !transactionsData) return null

  // TODO: put colors in theme file
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

  function getHeaderWidth(header) {
    switch (header) {
      case 'date':
        return '120px'
      case 'name':
        return '300px'
      case 'category':
        return '200px'
      case 'amount':
        return '200px'
      case 'options':
        return '100px'
      default: {
        return '100px'
      }
    }
  }

  return (
    <Box bg={bgColor} p={5} borderRadius={10}>
      <HStack mb={5} h={'65px'}>
        <FormControl w={'300px'}>
          <FormLabel>Search</FormLabel>
          <Input
            type="text"
            placeholder="Start typing..."
            onChange={(e) => setSearchStr(e.target.value)}
          />
        </FormControl>
        <FormControl w={'300px'}>
          <FormLabel>Category</FormLabel>
          <Select
            defaultValue={'all'}
            onChange={(event) => {
              setFilteredCategory(event.target.value)
            }}
          >
            <option value={'all'}>All</option>
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
        </FormControl>
        <Spacer />
        <Box pt={'32px'}>
          <DownloadCsvBtn transactionsData={transactionsData} />
        </Box>
        <Box pt={'32px'}>
          <CreateRuleBtn
            onSuccess={() => {
              setData([])
              onSuccess()
            }}
            forceRefresh={forceRefresh}
            icon={<FaSitemap />}
          />
        </Box>
        <Box pl={2} pt={'32px'}>
          <AddTransactionBtn
            onSuccess={() => {
              setData([])
              onSuccess()
            }}
            icon={<FaDollarSign />}
          />
        </Box>
      </HStack>

      <Table>
        <Thead p={10}>
          {table.getHeaderGroups().map((headerGroup) => (
            <Tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <Th
                    key={header.id}
                    onClick={header.column.getToggleSortingHandler()}
                    w={getHeaderWidth(header.id)}
                  >
                    {flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
                    {header.id !== 'options' ? (
                      <chakra.span pl="4">
                        {header.column.getIsSorted() ? (
                          header.column.getIsSorted() === 'desc' ? (
                            <TriangleDownIcon aria-label="sorted descending" />
                          ) : (
                            <TriangleUpIcon aria-label="sorted ascending" />
                          )
                        ) : (
                          <TriangleUpIcon color={'transparent'} />
                        )}
                      </chakra.span>
                    ) : (
                      <></>
                    )}
                  </Th>
                )
              })}
            </Tr>
          ))}
        </Thead>
        <Tbody>
          {table.getRowModel().rows.map((row) => {
            return (
              <Tr key={row.id}>
                {row.getVisibleCells().map((cell) => {
                  switch (cell.column.id) {
                    case 'date': {
                      return (
                        <Td key={cell.id}>
                          <Text alignItems={'center'}>
                            {DateTime.fromISO(cell.getValue(), {
                              zone: 'utc',
                            }).toFormat('LL/dd/yy')}
                          </Text>
                        </Td>
                      )
                    }
                    case 'name': {
                      return (
                        <Td key={cell.id}>
                          <Text>{cell.getValue()}</Text>
                        </Td>
                      )
                    }
                    case 'category': {
                      return (
                        <Td key={cell.id}>
                          <HStack>
                            <FaCircle
                              color={getCategoryColor(cell.getValue())}
                            />
                            <Select
                              w={'165px'}
                              value={cell.getValue()}
                              borderColor={bgColor}
                              onChange={async (event) => {
                                await updateCategory(
                                  cell.row.original.transactionId,
                                  event.target.value
                                )
                              }}
                            >
                              <option value={'bills'}>Bills</option>
                              <option value={'entertainment'}>
                                Entertainment
                              </option>
                              <option value={'groceries'}>Groceries</option>
                              <option value={'ignore'}>Ignore</option>
                              <option value={'income'}>Income</option>
                              <option value={'restaurant'}>Restaurant</option>
                              <option value={'transportation'}>
                                Transportation
                              </option>
                              <option value={'vacation'}>Vacation</option>
                              <option value={'uncategorized'}>
                                Uncategorized
                              </option>
                            </Select>
                          </HStack>
                        </Td>
                      )
                    }
                    case 'amount': {
                      return (
                        <Td key={cell.id}>
                          <Text>{currency.format(cell.getValue())}</Text>
                        </Td>
                      )
                    }
                    case 'options': {
                      return (
                        <Td key={cell.id}>
                          <HStack>
                            <EditTransactionBtn
                              onSuccess={() => {
                                setSearchStr('')
                                setData([])
                                onSuccess()
                              }}
                              forceRefresh={forceRefresh}
                              transaction={cell.row.original}
                              transactionsData={transactionsData}
                            />
                            {cell.row.original.enrollmentId ===
                            'user_created' ? (
                              <DeleteTransactionBtn
                                onSuccess={onSuccess}
                                forceRefresh={forceRefresh}
                                transaction={cell.row.original}
                              />
                            ) : (
                              <></>
                            )}
                          </HStack>
                        </Td>
                      )
                    }
                    default: {
                      logger(`unknown column: ${cell.column.id}`)
                      return <></>
                    }
                  }
                })}
              </Tr>
            )
          })}
        </Tbody>
      </Table>
    </Box>
  )
}
