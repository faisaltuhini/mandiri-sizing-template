export interface NumberFormatOptions {
  decimals?: number
  minDecimals?: number
  maxDecimals?: number
  prefix?: string
  suffix?: string
  fallback?: string
  locale?: string
}

export function formatNumber(
  value: number | string | null | undefined,
  options: NumberFormatOptions = {}
): string {
  if (value === null || value === undefined || value === '') {
    return options.fallback ?? ''
  }

  const num = typeof value === 'string' ? parseFormattedNumber(value) : value
  if (isNaN(num)) {
    return options.fallback ?? ''
  }

  const { decimals, minDecimals, maxDecimals, prefix = '', suffix = '', locale = 'id-ID' } = options

  const intlOptions: Intl.NumberFormatOptions = {
    useGrouping: true,
  }

  if (decimals !== undefined) {
    intlOptions.minimumFractionDigits = decimals
    intlOptions.maximumFractionDigits = decimals
  } else {
    if (minDecimals !== undefined) intlOptions.minimumFractionDigits = minDecimals
    if (maxDecimals !== undefined) intlOptions.maximumFractionDigits = maxDecimals
  }

  const formatted = new Intl.NumberFormat(locale, intlOptions).format(num)
  return `${prefix}${formatted}${suffix}`
}

export function parseFormattedNumber(value: string | number | null | undefined): number {
  if (value === null || value === undefined) return 0
  if (typeof value === 'number') return isNaN(value) ? 0 : value

  let str = value.toString().trim()
  if (!str) return 0

  // If string contains both dots and commas (e.g. "15.000.000,50" or "15,000,000.50"):
  if (str.includes('.') && str.includes(',')) {
    if (str.lastIndexOf(',') > str.lastIndexOf('.')) {
      // Indonesian: dots are thousands, comma is decimal -> "15.000.000,50"
      str = str.replace(/\./g, '').replace(',', '.')
    } else {
      // English: commas are thousands, dot is decimal -> "15,000,000.50"
      str = str.replace(/,/g, '')
    }
  } else if (str.includes(',')) {
    const commaCount = (str.match(/,/g) || []).length
    if (commaCount > 1) {
      // Multiple commas -> thousands separator
      str = str.replace(/,/g, '')
    } else {
      // Single comma -> Indonesian decimal separator
      str = str.replace(',', '.')
    }
  } else if (str.includes('.')) {
    const dotCount = (str.match(/\./g) || []).length
    if (dotCount > 1) {
      // Multiple dots -> Indonesian thousands separator
      str = str.replace(/\./g, '')
    } else {
      // Single dot: check if it represents thousands (e.g. "2.000" in ID) or decimal ("1.5")
      // In ID locale inputs, typing digits without comma creates dot separators.
    }
  }

  // Remove any remaining invalid characters except digits, minus, and dot
  const cleaned = str.replace(/[^0-9.-]/g, '')
  const parsed = parseFloat(cleaned)
  return isNaN(parsed) ? 0 : parsed
}

export function useFormat() {
  return {
    formatNumber,
    parseFormattedNumber,
  }
}
