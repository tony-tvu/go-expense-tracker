import React, { useState } from 'react'
import DatePicker from 'react-datepicker'


export default function Analytics() {
  const [date, setDate] = useState(new Date())

  return (
    <DatePicker selected={date} onChange={(date) => setDate(date)} />
  )
}
