import React, { createContext, useState } from 'react'
import { DateTime } from 'luxon'

export const AppStateContext = createContext()

const AppStateProvider = (props) => {
  const [appState, setAppState] = useState({
    selectedMonth: DateTime.now().month,
    selectedYear: DateTime.now().year,
  })

  return (
    <AppStateContext.Provider value={[appState, setAppState]}>
      {props.children}
    </AppStateContext.Provider>
  )
}

export default AppStateProvider
