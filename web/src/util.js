export const currency = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
})

export function timeSince(date) {
  var seconds = Math.floor((new Date() - date) / 1000)

  var interval = seconds / 31536000

  if (interval > 1) {
    return (
      Math.floor(interval) +
      `${Math.floor(interval) === 1 ? ' year' : ' years'}`
    )
  }
  interval = seconds / 2592000
  if (interval > 1) {
    return (
      Math.floor(interval) +
      `${Math.floor(interval) === 1 ? ' month' : ' months'}`
    )
  }
  interval = seconds / 86400
  if (interval > 1) {
    return (
      Math.floor(interval) + `${Math.floor(interval) === 1 ? ' day' : ' days'}`
    )
  }
  interval = seconds / 3600
  if (interval > 1) {
    return (
      Math.floor(interval) +
      `${Math.floor(interval) === 1 ? ' hour' : ' hours'}`
    )
  }
  interval = seconds / 60
  if (interval > 1) {
    return (
      Math.floor(interval) +
      `${Math.floor(interval) === 1 ? ' minute' : ' minutes'}`
    )
  }
  return (
    Math.floor(seconds) +
    `${Math.floor(interval) === 1 ? ' second' : ' seconds'}`
  )
}

export function loadScript(elementId, src) {
  if (!document.getElementById(elementId)) {
    const script = document.createElement('script')
    script.src = src
    script.id = elementId
    document.head.appendChild(script)
  }
}

export const formatUTC = (dateInt, addOffset = false) => {
  let date = !dateInt || dateInt.length < 1 ? new Date() : new Date(dateInt)
  if (typeof dateInt === 'string') {
    return date
  } else {
    const offset = addOffset
      ? date.getTimezoneOffset()
      : -date.getTimezoneOffset()
    const offsetDate = new Date()
    offsetDate.setTime(date.getTime() + offset * 60000)
    return offsetDate
  }
}
