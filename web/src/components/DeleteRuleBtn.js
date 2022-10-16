import React from 'react'
import { AlertDialog, AlertDialogBody, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogOverlay, Button, useDisclosure } from '@chakra-ui/react'
import logger from '../logger'
import { useNavigate } from 'react-router-dom'

export default function DeleteRuleBtn({ rule, onSuccess }) {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const cancelRef = React.useRef()
  const navigate = useNavigate()

  async function deleteAccount() {
    await fetch(`${process.env.REACT_APP_API_URL}/rules/${rule.id}`, {
      method: 'DELETE',
      credentials: 'include',
    })
      .then((res) => {
        if (res.status === 401) {
          navigate('/login')
        }
        if (res.status === 200) {
          onSuccess()
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
      type="button"
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
              Are you sure you want to remove "{rule.substring}"={rule.category}?
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
                onClick={async () => await deleteAccount()}
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
