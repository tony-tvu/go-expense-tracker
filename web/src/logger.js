var logger = function logger(...messages) {
  if (process.env.REACT_APP_ENV === "development") {
    console.log(messages)
  }
}

export default logger
