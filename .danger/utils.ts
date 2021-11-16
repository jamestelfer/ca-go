import * as fs from 'fs'
import * as path from 'path'
import { TestCoverageEntry } from './types'

export function getTestCoverageEntries(threshold: number, artifactPath: string): Array<TestCoverageEntry> {
  const artifact = path.join(process.cwd(), artifactPath)
  const data = fs.readFileSync(artifact).toString()

  return data.split(/\r?\n/)
  .filter((line) => line.trim().length)
  .map((line) => {
    const [ fileName, method, percentage ] = line.split('\t').filter((entry) => entry.length)
    const coverage = parseInt(percentage.slice(0, -1))

    return {
      fileName: fileName.slice(0, -1),
      method,
      coverage,
      isBelowThreshold: coverage < threshold
    }
  })
}
