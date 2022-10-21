export const currency = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
})

export function loadScript(elementId, src) {
  if (!document.getElementById(elementId)) {
    const script = document.createElement('script')
    script.src = src
    script.id = elementId
    document.head.appendChild(script)
  }
}
