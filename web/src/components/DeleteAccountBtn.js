import React from 'react'
import {
  Button,
  useDisclosure,
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogBody,
  AlertDialogFooter,
} from '@chakra-ui/react'
import logger from '../logger'
import { useNavigate } from 'react-router-dom'

export default function DeleteAccountBtn({ enrollment, onSuccess }) {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const cancelRef = React.useRef()
  const navigate = useNavigate()

  async function deleteAccount() {
    await fetch(`${process.env.REACT_APP_API_URL}/enrollments/${enrollment.enrollment_id}`, {
      method: 'DELETE',
      credentials: 'include',
    })
      .then((res) => {
        if (res.status === 401) {
          navigate('/login')
        }
        if (res.status === 200) {
          onSuccess()
          onClose()
        }
      })
      .catch((e) => {
        logger('error deleting account', e)
      })
  }

  return (
    <>
      <Button
        variant="solid"
        justifyContent="space-between"
        fontWeight="normal"
        fontSize="md"
        as="b"
        onClick={onOpen}
        color={'white'}
        bg={'#E63E3F'}
        _hover={{
          bg: '#AA2E2F',
        }}
      >
        Delete
      </Button>

      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Account
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure you want to remove {enrollment.institution}?
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              <Button
                color={'white'}
                bg={'#E63E3F'}
                _hover={{
                  bg: '#AA2E2F',
                }}
                onClick={deleteAccount}
                ml={3}
              >
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  )
}
