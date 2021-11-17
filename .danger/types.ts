export interface TestCoverageEntry {
  fileName: string
  method: string
  coverage: number
  isBelowThreshold: boolean
}
