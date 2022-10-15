import React, { useState } from 'react'
import {
  Button,
  useDisclosure,
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogBody,
  AlertDialogFooter,
  FormControl,
  FormLabel,
  Input,
  Select,
  Spinner,
  useToast,
} from '@chakra-ui/react'
import logger from '../logger'
import { BsPlus } from 'react-icons/bs'
import { useNavigate } from 'react-router-dom'
import { colors } from '../theme'

export default function CreateRuleBtn({ onSuccess }) {
  const [loading, setLoading] = useState(false)
  const [substring, setSubstring] = useState(null)
  const [category, setCategory] = useState('bills')
  const { isOpen, onOpen, onClose } = useDisclosure()
  const cancelRef = React.useRef()
  const navigate = useNavigate()
  const toast = useToast()

  async function createRule() {
    if (!substring) {
      toast({
        title: 'Error',
        description: 'Substring cannot be empty',
        status: 'error',
        position: 'top-right',
        duration: 5000,
        isClosable: true,
      })
      setLoading(false)
      return
    } else {
      setLoading(true)
      await fetch(`${process.env.REACT_APP_API_URL}/rules`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          substring: substring,
          category: category,
        }),
      })
        .then((res) => {
          if (res.status === 401) {
            navigate('/login')
            setLoading(false)
          }
          if (res.status === 200) {
            onClose()
            onSuccess()
            setLoading(false)
          }
        })
        .catch((e) => {
          logger('error creating rule', e)
        })
    }
  }

  return (
    <>
      <Button
        leftIcon={<BsPlus />}
        type="button"
        variant="solid"
        onClick={onOpen}
        bg={colors.primary}
        color={'white'}
        _hover={{
          bg: colors.primaryFaded,
        }}
      >
        Create Rule
      </Button>

      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Create Rule
            </AlertDialogHeader>

            <AlertDialogBody>
              <FormControl>
                <FormLabel>Substring (case-sensitive)</FormLabel>
                <Input
                  onChange={(event) => setSubstring(event.target.value)}
                  mb={3}
                />
                <FormLabel>Category</FormLabel>
                <Select
                  defaultValue={category}
                  onChange={async (event) => {
                    setCategory(event.target.value)
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
              </FormControl>
            </AlertDialogBody>
            <AlertDialogFooter>
              {loading && (
                <Spinner
                  thickness="4px"
                  speed="0.65s"
                  emptyColor="gray.200"
                  color="blue.500"
                  size="md"
                  mr={5}
                />
              )}
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              <Button
                bg={colors.primary}
                color={'white'}
                _hover={{
                  bg: colors.primaryFaded,
                }}
                onClick={() => {
                  setLoading(true)
                  createRule()
                }}
                ml={3}
                disabled={loading}
              >
                Create
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  )
}
